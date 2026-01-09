package project

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
)

func TestExecuteListProjects_Success(t *testing.T) {
	tests := []struct {
		name             string
		response         string
		expectedInOutput []string
	}{
		{
			name:     "multiple projects",
			response: testutil.SampleProjects,
			expectedInOutput: []string{
				"List of all available projects",
				"1. Project1",
				"2. Project2",
				"3. Project3",
			},
		},
		{
			name:     "single project",
			response: testutil.SingleProject,
			expectedInOutput: []string{
				"List of all available projects",
				"1. OnlyOne",
			},
		},
		{
			name:     "empty project list",
			response: testutil.EmptyProjects,
			expectedInOutput: []string{
				"List of all available projects",
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
				if !strings.HasSuffix(req.URL.Path, "/api/v2/projects") {
					t.Errorf("Expected path to end with /api/v2/projects, got: %s", req.URL.Path)
				}
				return h.MockResponse(200, tt.response), nil
			}

			err := ExecuteListProjects(h.Client, h.Stdout)
			h.AssertNoError(err)

			for _, expected := range tt.expectedInOutput {
				h.AssertOutputContains(expected)
			}
		})
	}
}

func TestExecuteListProjects_HTTPErrors(t *testing.T) {
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
			errorContains:  "failed to fetch projects",
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			response:       testutil.Error500,
			expectedStatus: 500,
			errorContains:  "failed to fetch projects",
		},
		{
			name:           "400 Bad Request",
			statusCode:     400,
			response:       testutil.Error400,
			expectedStatus: 400,
			errorContains:  "failed to fetch projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return h.MockResponse(tt.statusCode, tt.response), nil
			}

			err := ExecuteListProjects(h.Client, h.Stdout)
			h.AssertError(err, tt.errorContains)
			h.AssertHTTPError(err, tt.expectedStatus)
		})
	}
}

func TestExecuteListProjects_MalformedJSON(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.MalformedJSON), nil
	}

	err := ExecuteListProjects(h.Client, h.Stdout)
	h.AssertError(err, "failed to fetch projects")
}

func TestExecuteListProjects_EmptyResponse(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.EmptyBody), nil
	}

	err := ExecuteListProjects(h.Client, h.Stdout)
	h.AssertNoError(err)
	h.AssertOutputContains("List of all available projects")
}

func TestExecuteListProjects_NetworkError(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network timeout")
	}

	err := ExecuteListProjects(h.Client, h.Stdout)
	h.AssertError(err, "failed to fetch projects")
}

func TestExecuteListProjects_OutputFormatting(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.SampleProjects), nil
	}

	err := ExecuteListProjects(h.Client, h.Stdout)
	h.AssertNoError(err)

	output := h.GetOutput()

	// Check dividers are present
	if !strings.Contains(output, "==================================================================") {
		t.Error("Expected large divider in output")
	}
	if !strings.Contains(output, "------------------------------------------------------------------") {
		t.Error("Expected small divider in output")
	}

	// Check proper numbering
	lines := strings.Split(output, "\n")
	foundNumber := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "1. ") {
			foundNumber = true
			break
		}
	}
	if !foundNumber {
		t.Error("Expected numbered list starting with '1. '")
	}
}
