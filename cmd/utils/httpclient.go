package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"time"
)

// HTTPClient defines the interface for HTTP operations.
// This enables easy mocking in tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client provides high-level HTTP operations with retry logic.
type Client struct {
	httpClient  HTTPClient
	maxRetries  int
	baseBackoff time.Duration
}

// testClient is used for testing to inject a mock client.
// In production, this is always nil.
var testClient *Client

// NewClient creates a new HTTP client with default configuration:
// - 30 second timeout
// - 3 retry attempts with exponential backoff (1s, 2s, 4s)
// In tests, returns the injected test client if set.
func NewClient() *Client {
	if testClient != nil {
		return testClient
	}
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries:  3,
		baseBackoff: 1 * time.Second,
	}
}

// SetTestClient injects a test client for use in tests.
// This should only be used in test code.
func SetTestClient(client *Client) {
	testClient = client
}

// ClearTestClient removes the injected test client.
// This should be called in test cleanup.
func ClearTestClient() {
	testClient = nil
}

// NewClientWithHTTPClient creates a client with a custom HTTPClient.
// Useful for testing with mocks.
func NewClientWithHTTPClient(httpClient HTTPClient) *Client {
	return &Client{
		httpClient:  httpClient,
		maxRetries:  3,
		baseBackoff: 1 * time.Second,
	}
}

// SetBaseBackoff sets the base backoff duration for retries.
// This is primarily useful in tests to speed up retry scenarios.
func (c *Client) SetBaseBackoff(d time.Duration) {
	c.baseBackoff = d
}

// GetJSON performs a GET request and unmarshals the JSON response into target.
// The context controls cancellation and timeouts for the request.
func (c *Client) GetJSON(ctx context.Context, endpoint string, target interface{}) error {
	url, err := BuildAPIURL(endpoint)
	if err != nil {
		return fmt.Errorf("failed to build API URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return c.doWithRetry(ctx, req, target)
}

// PostJSON performs a POST request with optional JSON body and unmarshals the response.
// If body is nil, sends an empty POST request.
// If target is nil, response body is not unmarshaled.
// The context controls cancellation and timeouts for the request.
func (c *Client) PostJSON(ctx context.Context, endpoint string, body interface{}, target interface{}) error {
	url, err := BuildAPIURL(endpoint)
	if err != nil {
		return fmt.Errorf("failed to build API URL: %w", err)
	}

	var bodyReader io.Reader
	var jsonData []byte
	if body != nil {
		var err error
		jsonData, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		slog.Debug("JSON payload created", "payload", string(jsonData))
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		// Set GetBody to allow body recreation for retries
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(jsonData)), nil
		}
	}

	return c.doWithRetry(ctx, req, target)
}

// doWithRetry executes the request with exponential backoff retry logic.
// It respects context cancellation and will abort retries if the context is done.
func (c *Client) doWithRetry(ctx context.Context, req *http.Request, target interface{}) error {
	var lastErr error

	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		// Check context before each attempt
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("request cancelled: %w", err)
		}

		if attempt > 1 {
			backoff := c.calculateBackoff(attempt - 1)
			slog.Debug("Retrying request", "attempt", attempt, "backoff", backoff, "url", req.URL.String())

			// Use select to allow context cancellation during backoff
			select {
			case <-ctx.Done():
				return fmt.Errorf("request cancelled during retry backoff: %w", ctx.Err())
			case <-time.After(backoff):
			}

			// Recreate body reader for retry if needed
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return fmt.Errorf("failed to recreate request body for retry: %w", err)
				}
				req.Body = body
			}
		}

		slog.Debug("Sending request",
			"url", req.URL.String(),
			"method", req.Method,
			"attempt", attempt,
			"maxAttempts", c.maxRetries)

		resp, err := c.httpClient.Do(req) //nolint:bodyclose // body is closed in handleResponse
		if err != nil {
			lastErr = fmt.Errorf("request failed (attempt %d/%d): %w", attempt, c.maxRetries, err)
			slog.Debug("Request attempt failed", "error", err, "attempt", attempt)
			continue
		}

		// Handle response (closes body via defer)
		result, err := c.handleResponse(resp, target)
		if err != nil {
			lastErr = err
			// Retry on server errors (5xx)
			if httpErr, ok := err.(*HTTPError); ok && httpErr.StatusCode >= 500 {
				slog.Debug("Server error, will retry", "status", httpErr.StatusCode, "attempt", attempt)
				continue
			}
			// Don't retry client errors (4xx) or other errors
			return err
		}

		if result {
			return nil
		}
	}

	return fmt.Errorf("request failed after %d attempts: %w", c.maxRetries, lastErr)
}

// handleResponse processes the HTTP response, returning true on success.
func (c *Client) handleResponse(resp *http.Response, target interface{}) (bool, error) {
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			slog.Warn("Error closing response body", "error", cerr)
		}
	}()

	// Read body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for success status (2xx)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(bodyBytes),
			URL:        resp.Request.URL.String(),
		}
	}

	// Unmarshal if target provided
	if target != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return false, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	slog.Debug("Request completed successfully", "status", resp.Status)
	return true, nil
}

// calculateBackoff returns exponential backoff delay: 1s, 2s, 4s
func (c *Client) calculateBackoff(attempt int) time.Duration {
	multiplier := math.Pow(2, float64(attempt-1))
	return time.Duration(multiplier) * c.baseBackoff
}

// HTTPError represents an HTTP error response.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s (URL: %s)", e.StatusCode, e.Status, e.URL)
}

// IsHTTPError checks if error is an HTTPError and returns it.
func IsHTTPError(err error) (*HTTPError, bool) {
	httpErr, ok := err.(*HTTPError)
	return httpErr, ok
}
