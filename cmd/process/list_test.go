package process

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
	"ftsctl/cmd/utils"
)

func TestExecuteListProcesses_Success(t *testing.T) {
	tests := []struct {
		name             string
		response         string
		expectedInOutput []string
	}{
		{
			name:     "multiple processes",
			response: testutil.SampleProcessList,
			expectedInOutput: []string{
				"List of all transfer process statuses",
				"ProcessID: proc-123",
				"Phase: RUNNING",
				"ProcessID: proc-456",
				"Phase: COMPLETED",
				"TotalPatients:",
				"TotalBundles:",
			},
		},
		{
			name:     "empty process list",
			response: testutil.EmptyProcessList,
			expectedInOutput: []string{
				"List of all transfer process statuses",
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
				if !strings.HasSuffix(req.URL.Path, utils.EndpointProcessStatuses) {
					t.Errorf("Expected path to end with %s, got: %s", utils.EndpointProcessStatuses, req.URL.Path)
				}
				return h.MockResponse(200, tt.response), nil
			}

			err := ExecuteListProcesses(h.Client, h.Stdout)
			h.AssertNoError(err)

			for _, expected := range tt.expectedInOutput {
				h.AssertOutputContains(expected)
			}
		})
	}
}

func TestExecuteListProcesses_HTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		response       string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "404 Not Found",
			statusCode:     404,
			response:       testutil.Error404,
			expectedStatus: 404,
			errorContains:  "failed to fetch process statuses",
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			response:       testutil.Error500,
			expectedStatus: 500,
			errorContains:  "failed to fetch process statuses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return h.MockResponse(tt.statusCode, tt.response), nil
			}

			err := ExecuteListProcesses(h.Client, h.Stdout)
			h.AssertError(err, tt.errorContains)
			_ = h.AssertHTTPError(err, tt.expectedStatus)
		})
	}
}

func TestExecuteListProcesses_MalformedJSON(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.MalformedJSON), nil
	}

	err := ExecuteListProcesses(h.Client, h.Stdout)
	h.AssertError(err, "failed to fetch process statuses")
}

func TestExecuteListProcesses_NetworkError(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network timeout")
	}

	err := ExecuteListProcesses(h.Client, h.Stdout)
	h.AssertError(err, "failed to fetch process statuses")
}

func TestExecuteListProcesses_RunningProcess(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, `[`+testutil.SampleProcessRunning+`]`), nil
	}

	err := ExecuteListProcesses(h.Client, h.Stdout)
	h.AssertNoError(err)

	// Running process should show "still running" message
	h.AssertOutputContains("The process is still running")
	h.AssertOutputContains("proc-123")
	h.AssertOutputContains("RUNNING")
}

func TestExecuteListProcesses_CompletedProcess(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, `[`+testutil.SampleProcessCompleted+`]`), nil
	}

	err := ExecuteListProcesses(h.Client, h.Stdout)
	h.AssertNoError(err)

	// Completed process should show FinishedAt time (not "still running")
	h.AssertOutputNotContains("The process is still running")
	h.AssertOutputContains("proc-456")
	h.AssertOutputContains("COMPLETED")
	h.AssertOutputContains("FinishedAt:")
}

func TestExecuteListProcesses_ProcessFields(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, `[`+testutil.SampleProcessRunning+`]`), nil
	}

	err := ExecuteListProcesses(h.Client, h.Stdout)
	h.AssertNoError(err)

	// Verify all process fields are present
	expectedFields := []string{
		"ProcessID:",
		"Phase:",
		"CreatedAt:",
		"TotalPatients: 100",
		"TotalBundles: 500",
		"DeidentifiedBundles: 300",
		"SentBundles: 200",
		"SkippedBundles: 10",
	}

	for _, field := range expectedFields {
		h.AssertOutputContains(field)
	}
}
