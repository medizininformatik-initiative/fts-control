package mockserver

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Default test data
var (
	defaultProjects = []string{"Project1", "Project2", "Project3"}

	defaultProjectConfig = map[string]interface{}{
		"cohortSelector": map[string]interface{}{
			"trustCenterAgent": map[string]interface{}{
				"server":                  map[string]string{"baseUrl": "https://tc.example.com"},
				"domain":                  "test-domain",
				"patientIdentifierSystem": "http://example.com/pid",
				"policySystem":            "http://example.com/policy",
				"policies":                []string{"policy1", "policy2"},
			},
		},
		"dataSelector": map[string]interface{}{
			"everything": map[string]interface{}{
				"fhirServer": map[string]string{"baseUrl": "https://fhir.example.com"},
				"resolve":    map[string]string{"patientIdentifierSystem": "http://example.com/pid"},
			},
		},
	}

	defaultProcessList = []map[string]interface{}{
		{
			"processId":            "proc-123",
			"phase":                "RUNNING",
			"createdAt":            []int{2024, 1, 15, 10, 30, 0, 0},
			"finishedAt":           nil,
			"totalPatients":        100,
			"totalBundles":         500,
			"deidentifiedBundles":  300,
			"sentBundles":          200,
			"skippedBundles":       10,
		},
		{
			"processId":            "proc-456",
			"phase":                "COMPLETED",
			"createdAt":            []int{2024, 1, 15, 10, 30, 0, 0},
			"finishedAt":           []int{2024, 1, 15, 12, 0, 0, 0},
			"totalPatients":        50,
			"totalBundles":         250,
			"deidentifiedBundles":  250,
			"sentBundles":          250,
			"skippedBundles":       0,
		},
	}
)

// registerHandlers registers all API endpoint handlers.
func (s *MockCDAServer) registerHandlers() {
	s.mux.HandleFunc("/api/v2/projects", s.handleProjectList)
	s.mux.HandleFunc("/api/v2/projects/", s.handleProjectConfig)
	s.mux.HandleFunc("/api/v2/process/statuses", s.handleProcessList)
	s.mux.HandleFunc("/api/v2/process/status/", s.handleProcessStatus)
	s.mux.HandleFunc("/api/v2/process/", s.handleTransferStart)
}

// handleProjectList handles GET /api/v2/projects
func (s *MockCDAServer) handleProjectList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(defaultProjects); err != nil {
		s.t.Errorf("failed to encode projects: %v", err)
	}
}

// handleProjectConfig handles GET /api/v2/projects/{name}
func (s *MockCDAServer) handleProjectConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract project name from path
	name := strings.TrimPrefix(r.URL.Path, "/api/v2/projects/")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if project exists
	found := false
	for _, p := range defaultProjects {
		if p == name {
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "project not found"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(defaultProjectConfig); err != nil {
		s.t.Errorf("failed to encode project config: %v", err)
	}
}

// handleProcessList handles GET /api/v2/process/statuses
func (s *MockCDAServer) handleProcessList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(defaultProcessList); err != nil {
		s.t.Errorf("failed to encode process list: %v", err)
	}
}

// handleProcessStatus handles GET /api/v2/process/status/{id}
func (s *MockCDAServer) handleProcessStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract process ID from path
	id := strings.TrimPrefix(r.URL.Path, "/api/v2/process/status/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Find process
	for _, p := range defaultProcessList {
		if p["processId"] == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				s.t.Errorf("failed to encode process: %v", err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"error": "process not found"}`))
}

// handleTransferStart handles POST /api/v2/process/{project}/start
func (s *MockCDAServer) handleTransferStart(w http.ResponseWriter, r *http.Request) {
	// Only handle POST requests to /api/v2/process/{project}/start
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract path parts
	path := strings.TrimPrefix(r.URL.Path, "/api/v2/process/")
	if !strings.HasSuffix(path, "/start") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Extract project name
	project := strings.TrimSuffix(path, "/start")
	if project == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if project exists
	found := false
	for _, p := range defaultProjects {
		if p == project {
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "project not found"}`))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
