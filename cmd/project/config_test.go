package project

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
)

func TestExecuteProjectConfig_Success(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		if req.Method != "GET" {
			t.Errorf("Expected GET method, got: %s", req.Method)
		}
		if !strings.Contains(req.URL.Path, "/api/v2/projects/TestProject") {
			t.Errorf("Expected path to contain /api/v2/projects/TestProject, got: %s", req.URL.Path)
		}
		return h.MockResponse(200, testutil.SampleProjectConfig), nil
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertNoError(err)

	expectedOutputs := []string{
		"Project configuration with the Project Name: TestProject",
		"CohortSelector - TrustCenterAgent:",
		"https://tc.example.com",
		"test-domain",
		"http://example.com/pid",
		"http://example.com/policy",
		"policy1",
		"policy2",
		"DataSelector - Everything:",
		"https://fhir.example.com",
		"Deidentificator - Deidentifhir:",
		"P14D",
		"config.json",
		"scraper.json",
		"pseudonym-domain",
		"salt-domain",
		"dateshift-domain",
		"BundleSender - ResearchDomainAgent:",
		"https://research.example.com",
	}

	for _, expected := range expectedOutputs {
		h.AssertOutputContains(expected)
	}
}

func TestExecuteProjectConfig_HTTPErrors(t *testing.T) {
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
			errorContains:  "failed to fetch configuration",
		},
		{
			name:           "500 Server Error",
			statusCode:     500,
			response:       testutil.Error500,
			expectedStatus: 500,
			errorContains:  "failed to fetch configuration",
		},
		{
			name:           "400 Bad Request",
			statusCode:     400,
			response:       testutil.Error400,
			expectedStatus: 400,
			errorContains:  "failed to fetch configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return h.MockResponse(tt.statusCode, tt.response), nil
			}

			err := ExecuteProjectConfig(h.Client, h.Stdout, "NonExistent")
			h.AssertError(err, tt.errorContains)
			h.AssertHTTPError(err, tt.expectedStatus)
		})
	}
}

func TestExecuteProjectConfig_MalformedJSON(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.MalformedJSON), nil
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertError(err, "failed to fetch configuration")
}

func TestExecuteProjectConfig_EmptyResponse(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.EmptyBody), nil
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertNoError(err)
	h.AssertOutputContains("Project configuration with the Project Name: TestProject")
}

func TestExecuteProjectConfig_NetworkError(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertError(err, "failed to fetch configuration")
}

func TestExecuteProjectConfig_PartialJSON(t *testing.T) {
	// Test with minimal valid JSON structure
	partialConfig := `{
		"cohortSelector": {
			"trustCenterAgent": {
				"server": {"baseUrl": ""},
				"domain": "",
				"patientIdentifierSystem": "",
				"policySystem": "",
				"policies": []
			}
		},
		"dataSelector": {
			"everything": {
				"fhirServer": {"baseUrl": ""},
				"resolve": {"patientIdentifierSystem": ""}
			}
		},
		"deidentificator": {
			"deidentifhir": {
				"trustCenterAgent": {
					"server": {"baseUrl": ""},
					"domains": {
						"pseudonym": "",
						"salt": "",
						"dateShift": ""
					}
				},
				"maxDateShift": "",
				"deidentifhirConfig": "",
				"scraperConfig": ""
			}
		},
		"bundleSender": {
			"researchDomainAgent": {
				"server": {"baseUrl": ""},
				"project": ""
			}
		}
	}`

	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, partialConfig), nil
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertNoError(err)
	h.AssertOutputContains("Project configuration with the Project Name: TestProject")
}

func TestExecuteProjectConfig_OutputFormatting(t *testing.T) {
	h := testutil.NewTestHarness(t)
	h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return h.MockResponse(200, testutil.SampleProjectConfig), nil
	}

	err := ExecuteProjectConfig(h.Client, h.Stdout, "TestProject")
	h.AssertNoError(err)

	output := h.GetOutput()

	// Check section headers
	sections := []string{
		"CohortSelector - TrustCenterAgent:",
		"DataSelector - Everything:",
		"Deidentificator - Deidentifhir:",
		"BundleSender - ResearchDomainAgent:",
	}

	for _, section := range sections {
		if !strings.Contains(output, section) {
			t.Errorf("Expected section header %q in output", section)
		}
	}

	// Check indentation for nested items
	if !strings.Contains(output, "  Server BaseURL:") {
		t.Error("Expected proper indentation for nested items")
	}
}

func TestExecuteProjectConfig_ProjectNameInURL(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
	}{
		{name: "simple name", projectName: "MyProject"},
		{name: "with numbers", projectName: "Project123"},
		{name: "with dash", projectName: "my-project"},
		{name: "with underscore", projectName: "my_project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testutil.NewTestHarness(t)
			var capturedPath string
			h.MockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				capturedPath = req.URL.Path
				return h.MockResponse(200, testutil.SampleProjectConfig), nil
			}

			_ = ExecuteProjectConfig(h.Client, h.Stdout, tt.projectName)

			if !strings.Contains(capturedPath, tt.projectName) {
				t.Errorf("Expected URL path to contain %q, got: %s", tt.projectName, capturedPath)
			}
		})
	}
}
