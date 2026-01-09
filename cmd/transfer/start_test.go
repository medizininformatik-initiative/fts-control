package transfer

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
	"ftsctl/cmd/utils"
)

func TestExecuteStartTransfer_Success(t *testing.T) {
	tests := []struct {
		name             string
		projectName      string
		ids              []string
		expectedInOutput []string
		notInOutput      []string
		checkRequest     func(*testing.T, *http.Request)
	}{
		{
			name:        "all consented patients",
			projectName: "TestProject",
			ids:         nil,
			expectedInOutput: []string{
				"Transfer of project 'TestProject' has started",
				"all consented patients",
			},
			notInOutput: []string{
				"patient IDs:",
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				// Should be POST with no body
				if req.Method != "POST" {
					t.Errorf("Expected POST method, got: %s", req.Method)
				}
				if req.Body != nil {
					body, _ := io.ReadAll(req.Body)
					if len(body) > 0 {
						t.Errorf("Expected empty body for all patients, got: %s", string(body))
					}
				}
			},
		},
		{
			name:        "specific patient IDs",
			projectName: "Project2",
			ids:         []string{"patient1", "patient2", "patient3"},
			expectedInOutput: []string{
				"Transfer of project 'Project2' has started",
				"patient IDs:",
				"- patient1",
				"- patient2",
				"- patient3",
			},
			notInOutput: []string{
				"all consented patients",
			},
			checkRequest: func(t *testing.T, req *http.Request) {
				if req.Method != "POST" {
					t.Errorf("Expected POST method, got: %s", req.Method)
				}
				if req.Body != nil {
					body, _ := io.ReadAll(req.Body)
					bodyStr := string(body)
					if !strings.Contains(bodyStr, "patient1") {
						t.Errorf("Expected body to contain patient1, got: %s", bodyStr)
					}
				}
			},
		},
		{
			name:        "single patient ID",
			projectName: "Project3",
			ids:         []string{"single-patient"},
			expectedInOutput: []string{
				"Transfer of project 'Project3' has started",
				"- single-patient",
			},
			checkRequest: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				if tt.checkRequest != nil {
					tt.checkRequest(t, req)
				}
				// Return 202 Accepted for successful transfer start
				return h.MockResponse(202, ""), nil
			}

			err := ExecuteStartTransfer(h.Client, h.Stdout, tt.projectName, tt.ids)
			h.AssertNoError(err)

			for _, expected := range tt.expectedInOutput {
				h.AssertOutputContains(expected)
			}
			for _, unexpected := range tt.notInOutput {
				h.AssertOutputNotContains(unexpected)
			}
		})
	}
}

func TestExecuteStartTransfer_HTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		response       string
		expectedStatus int
		errorContains  string
	}{
		{
			name:           "404 Project Not Found",
			statusCode:     404,
			response:       testutil.Error404,
			expectedStatus: 404,
			errorContains:  "failed to start transfer",
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			response:       testutil.Error500,
			expectedStatus: 500,
			errorContains:  "failed to start transfer",
		},
		{
			name:           "400 Bad Request",
			statusCode:     400,
			response:       testutil.Error400,
			expectedStatus: 400,
			errorContains:  "failed to start transfer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return h.MockResponse(tt.statusCode, tt.response), nil
			}

			err := ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", nil)
			h.AssertError(err, tt.errorContains)
			_ = h.AssertHTTPError(err, tt.expectedStatus)
		})
	}
}

func TestExecuteStartTransfer_NetworkError(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	err := ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", nil)
	h.AssertError(err, "failed to start transfer")
}

func TestExecuteStartTransfer_ProjectNameInURL(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
	}{
		{name: "simple name", projectName: "MyProject"},
		{name: "with numbers", projectName: "Project123"},
		{name: "with dash", projectName: "my-project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			var capturedPath string
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				capturedPath = req.URL.Path
				return h.MockResponse(202, ""), nil
			}

			_ = ExecuteStartTransfer(h.Client, h.Stdout, tt.projectName, nil)

			expectedPath := utils.EndpointTransferStart(tt.projectName)
			if !strings.HasSuffix(capturedPath, expectedPath) {
				t.Errorf("Expected URL path to end with %q, got: %s", expectedPath, capturedPath)
			}
		})
	}
}

func TestExecuteStartTransfer_EmptyIDsSlice(t *testing.T) {
	// Empty slice should behave same as nil (transfer all patients)
	h := testutil.NewTestHarness(t)
	var capturedBody []byte
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.Body != nil {
			capturedBody, _ = io.ReadAll(req.Body)
		}
		return h.MockResponse(202, ""), nil
	}

	err := ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", []string{})
	h.AssertNoError(err)

	// Empty slice should not send a body
	if len(capturedBody) > 0 {
		t.Errorf("Expected no body for empty IDs slice, got: %s", string(capturedBody))
	}

	h.AssertOutputContains("all consented patients")
}

func TestExecuteStartTransfer_HTTP200Success(t *testing.T) {
	// Some APIs return 200 instead of 202, should still work
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, ""), nil
	}

	err := ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", nil)
	h.AssertNoError(err)
	h.AssertOutputContains("Transfer of project 'TestProject' has started")
}

func TestExecuteStartTransfer_OutputFormatting(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(202, ""), nil
	}

	err := ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", []string{"id1", "id2"})
	h.AssertNoError(err)

	output := h.GetOutput()

	// Check dividers are present
	if !strings.Contains(output, "==================================================================") {
		t.Error("Expected large divider in output")
	}
	if !strings.Contains(output, "------------------------------------------------------------------") {
		t.Error("Expected small divider in output")
	}

	// Check proper indentation for IDs
	if !strings.Contains(output, "   - id1") {
		t.Error("Expected proper indentation for patient IDs")
	}
}

func TestExecuteStartTransfer_RequestMethod(t *testing.T) {
	h := testutil.NewTestHarness(t)
	var capturedMethod string
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		capturedMethod = req.Method
		return h.MockResponse(202, ""), nil
	}

	_ = ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", nil)

	if capturedMethod != "POST" {
		t.Errorf("Expected POST method, got: %s", capturedMethod)
	}
}

func TestExecuteStartTransfer_ContentTypeHeader(t *testing.T) {
	h := testutil.NewTestHarness(t)

	t.Run("with IDs sets content-type", func(t *testing.T) {
		var contentType string
		h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			contentType = req.Header.Get("Content-Type")
			return h.MockResponse(202, ""), nil
		}

		_ = ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", []string{"id1"})

		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json when sending IDs, got: %s", contentType)
		}
	})

	t.Run("without IDs no content-type", func(t *testing.T) {
		var contentType string
		h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
			contentType = req.Header.Get("Content-Type")
			return h.MockResponse(202, ""), nil
		}

		_ = ExecuteStartTransfer(h.Client, h.Stdout, "TestProject", nil)

		if contentType != "" {
			t.Errorf("Expected no Content-Type when not sending IDs, got: %s", contentType)
		}
	})
}
