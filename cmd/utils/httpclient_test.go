package utils

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func init() {
	// Set up a default base URL for tests
	viper.Set("api.base_url", "http://localhost:8080")
}

// MockHTTPClient is a mock implementation of HTTPClient for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// newMockResponse creates a mock HTTP response.
func newMockResponse(statusCode int, body string) *http.Response {
	// Create a minimal request to avoid nil pointer issues
	req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
		Request:    req,
	}
}

func TestClient_GetJSON_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != "GET" {
				t.Errorf("Expected GET method, got: %s", req.Method)
			}
			return newMockResponse(200, `["Project1", "Project2"]`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	var projects []string

	err := client.GetJSON(context.Background(), "/api/v2/projects", &projects)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got: %d", len(projects))
	}

	if projects[0] != "Project1" {
		t.Errorf("Expected first project 'Project1', got: %s", projects[0])
	}
}

func TestClient_GetJSON_EmptyResponse(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(200, ""), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	var result interface{}

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected no error for empty response, got: %v", err)
	}
}

func TestClient_GetJSON_HTTPError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(404, `{"error": "not found"}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	var result interface{}

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error for 404 response")
	}

	httpErr, ok := IsHTTPError(err)
	if !ok {
		t.Fatalf("Expected HTTPError, got: %T", err)
	}

	if httpErr.StatusCode != 404 {
		t.Errorf("Expected status code 404, got: %d", httpErr.StatusCode)
	}

	if httpErr.Body != `{"error": "not found"}` {
		t.Errorf("Expected body to be preserved, got: %s", httpErr.Body)
	}
}

func TestClient_PostJSON_Success(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != "POST" {
				t.Errorf("Expected POST method, got: %s", req.Method)
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type application/json, got: %s", req.Header.Get("Content-Type"))
			}
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	payload := map[string]string{"key": "value"}
	var result map[string]bool

	err := client.PostJSON(context.Background(), "/api/v2/test", payload, &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !result["success"] {
		t.Error("Expected success to be true")
	}
}

func TestClient_PostJSON_NilBody(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Method != "POST" {
				t.Errorf("Expected POST method, got: %s", req.Method)
			}
			// Content-Type should not be set for nil body
			if req.Header.Get("Content-Type") != "" {
				t.Errorf("Expected no Content-Type for nil body, got: %s", req.Header.Get("Content-Type"))
			}
			return newMockResponse(202, ""), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)

	err := client.PostJSON(context.Background(), "/api/v2/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestClient_PostJSON_AcceptedStatus(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(202, ""), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)

	err := client.PostJSON(context.Background(), "/api/v2/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error for 202 Accepted, got: %v", err)
	}
}

func TestClient_Retry_OnNetworkError(t *testing.T) {
	var attempts int32

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&attempts, 1)
			current := atomic.LoadInt32(&attempts)
			if current < 3 {
				return nil, errors.New("network error")
			}
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	client.baseBackoff = 10 * time.Millisecond // Speed up test
	var result map[string]bool

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected success after retries, got: %v", err)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

func TestClient_Retry_OnServerError(t *testing.T) {
	var attempts int32

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&attempts, 1)
			current := atomic.LoadInt32(&attempts)
			if current < 3 {
				return newMockResponse(500, "Server Error"), nil
			}
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	client.baseBackoff = 10 * time.Millisecond // Speed up test
	var result map[string]bool

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected success after retries, got: %v", err)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

func TestClient_NoRetry_OnClientError(t *testing.T) {
	var attempts int32

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&attempts, 1)
			return newMockResponse(400, "Bad Request"), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	client.baseBackoff = 10 * time.Millisecond
	var result interface{}

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error for 400 response")
	}

	// Should not retry on 4xx errors
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected 1 attempt (no retry for 4xx), got: %d", attempts)
	}
}

func TestClient_MaxRetries_Exceeded(t *testing.T) {
	var attempts int32

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&attempts, 1)
			return nil, errors.New("network error")
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	client.baseBackoff = 10 * time.Millisecond
	var result interface{}

	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error after max retries")
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

func TestClient_CalculateBackoff(t *testing.T) {
	client := NewClient()

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{1, 1 * time.Second},  // 2^0 = 1
		{2, 2 * time.Second},  // 2^1 = 2
		{3, 4 * time.Second},  // 2^2 = 4
		{4, 8 * time.Second},  // 2^3 = 8
	}

	for _, tt := range tests {
		got := client.calculateBackoff(tt.attempt)
		if got != tt.expected {
			t.Errorf("calculateBackoff(%d) = %v, expected %v", tt.attempt, got, tt.expected)
		}
	}
}

func TestHTTPError_Formatting(t *testing.T) {
	err := &HTTPError{
		StatusCode: 404,
		Status:     "404 Not Found",
		URL:        "http://localhost/api/test",
		Body:       "Resource not found",
	}

	expected := "HTTP 404: 404 Not Found (URL: http://localhost/api/test)"
	if err.Error() != expected {
		t.Errorf("Expected error message: %s, got: %s", expected, err.Error())
	}
}

func TestIsHTTPError(t *testing.T) {
	httpErr := &HTTPError{StatusCode: 500}
	otherErr := errors.New("other error")

	if got, ok := IsHTTPError(httpErr); !ok || got.StatusCode != 500 {
		t.Error("IsHTTPError should return true for HTTPError")
	}

	if _, ok := IsHTTPError(otherErr); ok {
		t.Error("IsHTTPError should return false for non-HTTPError")
	}
}

// Integration test using httptest server
func TestClient_Integration_WithHTTPTestServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2/projects":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`["Project1", "Project2", "Project3"]`))
		case "/api/v2/process/start":
			if r.Method != "POST" {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client that points to test server
	client := &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		maxRetries:  3,
		baseBackoff: 10 * time.Millisecond,
	}

	t.Run("GET projects", func(t *testing.T) {
		var projects []string
		// Note: We need to build URL manually for test server
		req, _ := http.NewRequestWithContext(context.Background(), "GET", server.URL+"/api/v2/projects", nil)
		err := client.doWithRetry(context.Background(), req, &projects)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(projects) != 3 {
			t.Errorf("Expected 3 projects, got: %d", len(projects))
		}
	})

	t.Run("POST start process", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(context.Background(), "POST", server.URL+"/api/v2/process/start", nil)
		err := client.doWithRetry(context.Background(), req, nil)
		if err != nil {
			t.Fatalf("Expected no error for 202 response, got: %v", err)
		}
	})
}

// Context cancellation tests

func TestClient_GetJSON_ContextCancelled(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Simulate a slow request that respects context
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(100 * time.Millisecond):
				return newMockResponse(200, `{"success": true}`), nil
			}
		},
	}

	client := NewClientWithHTTPClient(mockClient)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result interface{}
	err := client.GetJSON(ctx, "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestClient_GetJSON_ContextTimeout(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Simulate a slow request
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(200 * time.Millisecond):
				return newMockResponse(200, `{"success": true}`), nil
			}
		},
	}

	client := NewClientWithHTTPClient(mockClient)

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var result interface{}
	err := client.GetJSON(ctx, "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error for context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded error, got: %v", err)
	}
}

func TestClient_Retry_ContextCancelledDuringBackoff(t *testing.T) {
	var attempts int32

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			atomic.AddInt32(&attempts, 1)
			// Always fail to trigger retry
			return nil, errors.New("network error")
		},
	}

	client := NewClientWithHTTPClient(mockClient)
	client.baseBackoff = 100 * time.Millisecond // Long enough backoff to cancel during

	// Create a context that we'll cancel after the first attempt
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after a short delay (after first attempt but during backoff)
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	var result interface{}
	err := client.GetJSON(ctx, "/api/v2/test", &result)
	if err == nil {
		t.Fatal("Expected error when context cancelled during backoff")
	}

	// Should have only made 1 attempt before context was cancelled during backoff
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected 1 attempt before cancel, got: %d", attempts)
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestClient_PostJSON_ContextCancelled(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(100 * time.Millisecond):
				return newMockResponse(200, `{"success": true}`), nil
			}
		},
	}

	client := NewClientWithHTTPClient(mockClient)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	payload := map[string]string{"key": "value"}
	var result map[string]bool
	err := client.PostJSON(ctx, "/api/v2/test", payload, &result)
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

// Builder method tests

func TestClient_WithHeader(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithHeader("X-Custom-Header", "custom-value")

	var result map[string]bool
	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if capturedHeaders.Get("X-Custom-Header") != "custom-value" {
		t.Errorf("Expected X-Custom-Header to be 'custom-value', got: %s",
			capturedHeaders.Get("X-Custom-Header"))
	}
}

func TestClient_WithHeaders(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	headers := map[string]string{
		"X-Header-1": "value1",
		"X-Header-2": "value2",
	}
	client := NewClientWithHTTPClient(mockClient).WithHeaders(headers)

	var result map[string]bool
	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if capturedHeaders.Get("X-Header-1") != "value1" {
		t.Errorf("Expected X-Header-1 to be 'value1', got: %s",
			capturedHeaders.Get("X-Header-1"))
	}
	if capturedHeaders.Get("X-Header-2") != "value2" {
		t.Errorf("Expected X-Header-2 to be 'value2', got: %s",
			capturedHeaders.Get("X-Header-2"))
	}
}

func TestClient_WithBearerToken(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithBearerToken("secret-token-123")

	var result map[string]bool
	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expected := "Bearer secret-token-123"
	if capturedHeaders.Get("Authorization") != expected {
		t.Errorf("Expected Authorization to be '%s', got: %s",
			expected, capturedHeaders.Get("Authorization"))
	}
}

func TestClient_WithAPIKey(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return newMockResponse(200, `{"success": true}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithAPIKey("api-key-456")

	var result map[string]bool
	err := client.GetJSON(context.Background(), "/api/v2/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if capturedHeaders.Get("X-API-Key") != "api-key-456" {
		t.Errorf("Expected X-API-Key to be 'api-key-456', got: %s",
			capturedHeaders.Get("X-API-Key"))
	}
}

func TestClient_ImmutableHeaders(t *testing.T) {
	baseClient := NewClientWithHTTPClient(&MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return newMockResponse(200, `{}`), nil
		},
	})

	derivedClient := baseClient.WithHeader("X-Test", "value")

	if len(baseClient.headers) != 0 {
		t.Error("Base client was modified - immutability violated")
	}

	if derivedClient.headers.Get("X-Test") != "value" {
		t.Error("Derived client should have the header")
	}
}

func TestClient_ChainedHeaders(t *testing.T) {
	var capturedHeaders http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			capturedHeaders = req.Header
			return newMockResponse(200, `{}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithHeader("X-Header-1", "value1").
		WithHeader("X-Header-2", "value2").
		WithBearerToken("token")

	err := client.GetJSON(context.Background(), "/api/v2/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if capturedHeaders.Get("X-Header-1") != "value1" {
		t.Error("First header not applied")
	}
	if capturedHeaders.Get("X-Header-2") != "value2" {
		t.Error("Second header not applied")
	}
	if capturedHeaders.Get("Authorization") != "Bearer token" {
		t.Error("Bearer token not applied")
	}
}

func TestClient_HeadersAppliedToRetries(t *testing.T) {
	var attemptHeaders []http.Header
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			headerCopy := make(http.Header)
			for key, values := range req.Header {
				headerCopy[key] = append([]string(nil), values...)
			}
			attemptHeaders = append(attemptHeaders, headerCopy)

			if len(attemptHeaders) < 2 {
				return nil, errors.New("network error")
			}
			return newMockResponse(200, `{}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithHeader("X-Custom", "value")
	client.SetBaseBackoff(10 * time.Millisecond)

	err := client.GetJSON(context.Background(), "/api/v2/test", nil)
	if err != nil {
		t.Fatalf("Expected success after retry, got: %v", err)
	}

	if len(attemptHeaders) != 2 {
		t.Fatalf("Expected 2 attempts, got: %d", len(attemptHeaders))
	}

	for i, headers := range attemptHeaders {
		if headers.Get("X-Custom") != "value" {
			t.Errorf("Header missing in attempt %d", i+1)
		}
	}
}

func TestClient_HeadersNotOverwriteRequestHeaders(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Request Content-Type should not be overwritten, got: %s",
					req.Header.Get("Content-Type"))
			}
			if req.Header.Get("X-Custom") != "value" {
				t.Error("Custom header should be applied")
			}
			return newMockResponse(200, `{}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithHeader("Content-Type", "text/plain").
		WithHeader("X-Custom", "value")

	payload := map[string]string{"key": "value"}
	var result map[string]interface{}
	err := client.PostJSON(context.Background(), "/api/v2/test", payload, &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestClient_CopyHeaders_DeepCopy(t *testing.T) {
	original := make(http.Header)
	original.Set("X-Original", "value1")
	original.Add("X-Multi", "value1")
	original.Add("X-Multi", "value2")

	copied := copyHeaders(original)

	copied.Set("X-Original", "modified")
	copied.Set("X-New", "newvalue")
	copied.Add("X-Multi", "value3")

	if original.Get("X-Original") != "value1" {
		t.Error("Original header was modified")
	}
	if original.Get("X-New") != "" {
		t.Error("Original should not have new header")
	}
	if len(original["X-Multi"]) != 2 {
		t.Error("Original multi-value header was modified")
	}
}

func TestClient_WithBaseBackoff_Immutable(t *testing.T) {
	client := NewClient()
	originalBackoff := client.baseBackoff

	clientWithNewBackoff := client.WithBaseBackoff(5 * time.Second)

	if client.baseBackoff != originalBackoff {
		t.Error("Original client baseBackoff was mutated")
	}

	if clientWithNewBackoff.baseBackoff != 5*time.Second {
		t.Error("New client should have updated baseBackoff")
	}
}

func TestClient_HeadersInPostJSON(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Trace-ID") != "trace-123" {
				t.Error("Custom header not applied to POST request")
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Error("Content-Type should be application/json for POST with body")
			}
			return newMockResponse(201, `{"id": 42}`), nil
		},
	}

	client := NewClientWithHTTPClient(mockClient).
		WithHeader("X-Trace-ID", "trace-123")

	payload := map[string]string{"name": "test"}
	var result map[string]int
	err := client.PostJSON(context.Background(), "/api/v2/test", payload, &result)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result["id"] != 42 {
		t.Errorf("Expected id=42, got: %d", result["id"])
	}
}
