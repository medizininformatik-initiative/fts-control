package utils

import (
	"bytes"
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

	err := client.GetJSON("/api/v2/projects", &projects)
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

	err := client.GetJSON("/api/v2/test", &result)
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

	err := client.GetJSON("/api/v2/test", &result)
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

	err := client.PostJSON("/api/v2/test", payload, &result)
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

	err := client.PostJSON("/api/v2/test", nil, nil)
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

	err := client.PostJSON("/api/v2/test", nil, nil)
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

	err := client.GetJSON("/api/v2/test", &result)
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

	err := client.GetJSON("/api/v2/test", &result)
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

	err := client.GetJSON("/api/v2/test", &result)
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

	err := client.GetJSON("/api/v2/test", &result)
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
		req, _ := http.NewRequest("GET", server.URL+"/api/v2/projects", nil)
		err := client.doWithRetry(req, &projects)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(projects) != 3 {
			t.Errorf("Expected 3 projects, got: %d", len(projects))
		}
	})

	t.Run("POST start process", func(t *testing.T) {
		req, _ := http.NewRequest("POST", server.URL+"/api/v2/process/start", nil)
		err := client.doWithRetry(req, nil)
		if err != nil {
			t.Fatalf("Expected no error for 202 response, got: %v", err)
		}
	})
}
