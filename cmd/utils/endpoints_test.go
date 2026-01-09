package utils

import (
	"testing"
)

func TestEndpointConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"projects endpoint", EndpointProjects, "/api/v2/projects"},
		{"process statuses endpoint", EndpointProcessStatuses, "/api/v2/process/statuses"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, tt.constant)
			}
		})
	}
}

func TestEndpointProjectConfig(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantPath    string
	}{
		{"valid project name", "MyProject", "/api/v2/projects/MyProject"},
		{"project with dash", "my-project", "/api/v2/projects/my-project"},
		{"project with numbers", "project123", "/api/v2/projects/project123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndpointProjectConfig(tt.projectName)
			if got != tt.wantPath {
				t.Errorf("Expected path %q, got %q", tt.wantPath, got)
			}
		})
	}
}

func TestEndpointProcessStatus(t *testing.T) {
	tests := []struct {
		name      string
		processID string
		wantPath  string
	}{
		{"valid process ID", "proc-123", "/api/v2/process/status/proc-123"},
		{"UUID process ID", "550e8400-e29b-41d4-a716-446655440000", "/api/v2/process/status/550e8400-e29b-41d4-a716-446655440000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndpointProcessStatus(tt.processID)
			if got != tt.wantPath {
				t.Errorf("Expected path %q, got %q", tt.wantPath, got)
			}
		})
	}
}

func TestEndpointTransferStart(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		wantPath    string
	}{
		{"valid project name", "TestProject", "/api/v2/process/TestProject/start"},
		{"project with dash", "my-project", "/api/v2/process/my-project/start"},
		{"project with numbers", "project123", "/api/v2/process/project123/start"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndpointTransferStart(tt.projectName)
			if got != tt.wantPath {
				t.Errorf("Expected path %q, got %q", tt.wantPath, got)
			}
		})
	}
}
