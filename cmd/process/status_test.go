package process

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
	"ftsctl/cmd/utils"
)

func TestExecuteProcessStatus_Success(t *testing.T) {
	tests := []struct {
		name             string
		processID        string
		response         string
		expectedInOutput []string
	}{
		{
			name:      "running process",
			processID: "proc-123",
			response:  testutil.SampleProcessRunning,
			expectedInOutput: []string{
				"Representation of the Process with the ProcessID: proc-123",
				"ProcessID: proc-123",
				"Phase: RUNNING",
				"The process is still running",
				"TotalPatients: 100",
				"TotalBundles: 500",
			},
		},
		{
			name:      "completed process",
			processID: "proc-456",
			response:  testutil.SampleProcessCompleted,
			expectedInOutput: []string{
				"Representation of the Process with the ProcessID: proc-456",
				"ProcessID: proc-456",
				"Phase: COMPLETED",
				"FinishedAt:",
				"TotalPatients: 50",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				if req.Method != "GET" {
					t.Errorf("Expected GET method, got: %s", req.Method)
				}
				expectedPath := utils.EndpointProcessStatus(tt.processID)
				if !strings.Contains(req.URL.Path, expectedPath) {
					t.Errorf("Expected path to contain %s, got: %s", expectedPath, req.URL.Path)
				}
				if !strings.HasSuffix(req.URL.Path, tt.processID) {
					t.Errorf("Expected path to end with %s, got: %s", tt.processID, req.URL.Path)
				}
				return h.MockResponse(200, tt.response), nil
			}

			err := ExecuteProcessStatus(h.Client, h.Stdout, tt.processID)
			h.AssertNoError(err)

			for _, expected := range tt.expectedInOutput {
				h.AssertOutputContains(expected)
			}
		})
	}
}

func TestExecuteProcessStatus_HTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		response       string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "404 Process Not Found",
			statusCode:     404,
			response:       testutil.Error404,
			expectedStatus: 404,
			errorContains:  "failed to fetch status for process",
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			response:       testutil.Error500,
			expectedStatus: 500,
			errorContains:  "failed to fetch status for process",
		},
		{
			name:           "400 Bad Request",
			statusCode:     400,
			response:       testutil.Error400,
			expectedStatus: 400,
			errorContains:  "failed to fetch status for process",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return h.MockResponse(tt.statusCode, tt.response), nil
			}

			err := ExecuteProcessStatus(h.Client, h.Stdout, "non-existent")
			h.AssertError(err, tt.errorContains)
			_ = h.AssertHTTPError(err, tt.expectedStatus)
		})
	}
}

func TestExecuteProcessStatus_MalformedJSON(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.MalformedJSON), nil
	}

	err := ExecuteProcessStatus(h.Client, h.Stdout, "proc-123")
	h.AssertError(err, "failed to fetch status for process")
}

func TestExecuteProcessStatus_EmptyResponse(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.EmptyBody), nil
	}

	// Empty response with 200 OK should still work (empty process struct)
	err := ExecuteProcessStatus(h.Client, h.Stdout, "proc-123")
	h.AssertNoError(err)
	h.AssertOutputContains("Representation of the Process with the ProcessID: proc-123")
}

func TestExecuteProcessStatus_NetworkError(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	err := ExecuteProcessStatus(h.Client, h.Stdout, "proc-123")
	h.AssertError(err, "failed to fetch status for process")
}

func TestExecuteProcessStatus_ProcessIDInURL(t *testing.T) {
	tests := []struct {
		name      string
		processID string
	}{
		{name: "simple id", processID: "proc-123"},
		{name: "uuid-like", processID: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"},
		{name: "numeric", processID: "12345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			var capturedPath string
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				capturedPath = req.URL.Path
				return h.MockResponse(200, testutil.SampleProcessRunning), nil
			}

			_ = ExecuteProcessStatus(h.Client, h.Stdout, tt.processID)

			if !strings.HasSuffix(capturedPath, tt.processID) {
				t.Errorf("Expected URL path to end with %q, got: %s", tt.processID, capturedPath)
			}
		})
	}
}

func TestExecuteProcessStatus_ProcessIDInOutput(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.SampleProcessRunning), nil
	}

	processID := "my-custom-process-id"
	err := ExecuteProcessStatus(h.Client, h.Stdout, processID)
	h.AssertNoError(err)

	// The header should contain the provided process ID
	h.AssertOutputContains("Representation of the Process with the ProcessID: " + processID)
}

func TestExecuteProcessStatus_AllFieldsPresent(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.SampleProcessRunning), nil
	}

	err := ExecuteProcessStatus(h.Client, h.Stdout, "proc-123")
	h.AssertNoError(err)

	// Verify all process fields are present in output
	expectedFields := []string{
		"ProcessID:",
		"Phase:",
		"CreatedAt:",
		"FinishedAt:",
		"TotalPatients:",
		"TotalBundles:",
		"DeidentifiedBundles:",
		"SentBundles:",
		"SkippedBundles:",
	}

	for _, field := range expectedFields {
		h.AssertOutputContains(field)
	}
}

func TestExecuteProcessStatus_OutputFormatting(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.SampleProcessRunning), nil
	}

	err := ExecuteProcessStatus(h.Client, h.Stdout, "proc-123")
	h.AssertNoError(err)

	output := h.GetOutput()

	// Check dividers are present
	if !strings.Contains(output, "==================================================================") {
		t.Error("Expected large divider in output")
	}
	if !strings.Contains(output, "------------------------------------------------------------------") {
		t.Error("Expected small divider in output")
	}
}
