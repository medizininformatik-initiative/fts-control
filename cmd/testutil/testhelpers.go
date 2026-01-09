// Package testutil provides testing utilities for ftsctl commands.
package testutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestHarness provides utilities for testing commands.
type TestHarness struct {
	T          *testing.T
	MockClient *MockHTTPClient
	Client     *utils.Client
	Stdout     *bytes.Buffer
	Stderr     *bytes.Buffer
}

// NewTestHarness creates a new test harness with a mock HTTP client.
func NewTestHarness(t *testing.T) *TestHarness {
	t.Helper()

	// Configure viper for tests
	viper.Set("api.base_url", "http://localhost:8080")

	mockClient := &MockHTTPClient{}
	client := utils.NewClientWithHTTPClient(mockClient)
	// Speed up retries for tests
	client.SetBaseBackoff(10 * time.Millisecond)

	h := &TestHarness{
		T:          t,
		MockClient: mockClient,
		Client:     client,
		Stdout:     &bytes.Buffer{},
		Stderr:     &bytes.Buffer{},
	}

	return h
}

// SetupTestClient injects the test client globally and returns a cleanup function.
func (h *TestHarness) SetupTestClient() func() {
	utils.SetTestClient(h.Client)
	return func() {
		utils.ClearTestClient()
	}
}

// MockResponse creates a mock HTTP response.
func (h *TestHarness) MockResponse(statusCode int, body string) *http.Response {
	req, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
		Request:    req,
	}
}

// ExecuteCommand executes a Cobra command with the mock client injected.
func (h *TestHarness) ExecuteCommand(cmd *cobra.Command, args ...string) error {
	cleanup := h.SetupTestClient()
	defer cleanup()

	// Reset buffers
	h.Stdout.Reset()
	h.Stderr.Reset()

	// Configure command output
	cmd.SetOut(h.Stdout)
	cmd.SetErr(h.Stderr)
	cmd.SetArgs(args)

	// Reset flags to default values before execution
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})

	return cmd.Execute()
}

// Assertion helpers

// AssertNoError fails the test if err is not nil.
func (h *TestHarness) AssertNoError(err error) {
	h.T.Helper()
	if err != nil {
		h.T.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError fails the test if err is nil or doesn't contain msgContains.
func (h *TestHarness) AssertError(err error, msgContains string) {
	h.T.Helper()
	if err == nil {
		h.T.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), msgContains) {
		h.T.Errorf("Expected error to contain %q, got: %v", msgContains, err)
	}
}

// AssertHTTPError fails the test if err is not an HTTPError with the expected status.
func (h *TestHarness) AssertHTTPError(err error, expectedStatus int) *utils.HTTPError {
	h.T.Helper()
	if err == nil {
		h.T.Fatal("Expected HTTPError, got nil")
	}

	// Unwrap error chain to find HTTPError
	var httpErr *utils.HTTPError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != expectedStatus {
			h.T.Errorf("Expected status %d, got: %d", expectedStatus, httpErr.StatusCode)
		}
		return httpErr
	}

	h.T.Fatalf("Expected HTTPError, got: %T (%v)", err, err)
	return nil
}

// AssertOutputContains fails if stdout doesn't contain the expected string.
func (h *TestHarness) AssertOutputContains(expected string) {
	h.T.Helper()
	output := h.Stdout.String()
	if !strings.Contains(output, expected) {
		h.T.Errorf("Expected output to contain %q\nGot:\n%s", expected, output)
	}
}

// AssertOutputNotContains fails if stdout contains the unexpected string.
func (h *TestHarness) AssertOutputNotContains(unexpected string) {
	h.T.Helper()
	output := h.Stdout.String()
	if strings.Contains(output, unexpected) {
		h.T.Errorf("Expected output NOT to contain %q\nGot:\n%s", unexpected, output)
	}
}

// GetOutput returns the current stdout content.
func (h *TestHarness) GetOutput() string {
	return h.Stdout.String()
}

// GetStderr returns the current stderr content.
func (h *TestHarness) GetStderr() string {
	return h.Stderr.String()
}
