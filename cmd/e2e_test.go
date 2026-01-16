package cmd

import (
	"bytes"
	"context"
	"testing"

	"ftsctl/cmd/testutil"
	"ftsctl/cmd/testutil/mockserver"
	"ftsctl/cmd/utils"

	"github.com/spf13/viper"
)

// TestE2E_NoAuth tests commands against a server with no authentication.
func TestE2E_NoAuth(t *testing.T) {
	server := mockserver.New().WithAuthNone().Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"project list", testProjectList},
		{"project config", testProjectConfig},
		{"process list", testProcessList},
		{"process status", testProcessStatus},
		{"transfer start", testTransferStart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestE2E_BasicAuth tests commands against a server with Basic authentication.
func TestE2E_BasicAuth(t *testing.T) {
	const testUser = "testuser"
	const testPass = "testpass"

	server := mockserver.New().WithBasicAuth(testUser, testPass).Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.basic.user", testUser)
	viper.Set("auth.basic.password", testPass)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"project list", testProjectList},
		{"project config", testProjectConfig},
		{"process list", testProcessList},
		{"process status", testProcessStatus},
		{"transfer start", testTransferStart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestE2E_Certificate tests commands against a server with mTLS certificate authentication.
func TestE2E_Certificate(t *testing.T) {
	certs := mockserver.SetupTestCertificates(t)
	server := mockserver.New().WithCertificate(certs).Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.certificate.cert_file", certs.ClientCertFile)
	viper.Set("auth.certificate.key_file", certs.ClientKeyFile)
	viper.Set("auth.certificate.ca_file", certs.CACertFile)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"project list", testProjectList},
		{"project config", testProjectConfig},
		{"process list", testProcessList},
		{"process status", testProcessStatus},
		{"transfer start", testTransferStart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestE2E_OAuth2_RealKeycloak tests commands with real Keycloak OAuth2 authentication.
// This test requires a running Keycloak instance. In CI, Keycloak is started by the workflow.
// For local testing, run: docker compose -f test/oauth2/compose.yaml up -d --build
func TestE2E_OAuth2_RealKeycloak(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	if !testutil.IsKeycloakRunning() {
		t.Skip("Keycloak is not running. Start it with: docker compose -f test/oauth2/compose.yaml up -d --build")
	}

	keycloakConfig := testutil.DefaultKeycloakConfig()

	// Start mock CDA server with OAuth2 (just validates bearer token exists)
	server := mockserver.New().WithOAuth2().Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.oauth2.token_url", keycloakConfig.TokenURL)
	viper.Set("auth.oauth2.client_id", keycloakConfig.ClientID)
	viper.Set("auth.oauth2.client_secret", keycloakConfig.ClientSecret)

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{"project list", testProjectList},
		{"project config", testProjectConfig},
		{"process list", testProcessList},
		{"process status", testProcessStatus},
		{"transfer start", testTransferStart},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// Test helper functions

func testProjectList(t *testing.T) {
	t.Helper()

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var projects []string
	err = client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err != nil {
		t.Fatalf("failed to fetch projects: %v", err)
	}

	if len(projects) != 3 {
		t.Errorf("expected 3 projects, got %d", len(projects))
	}
	if projects[0] != "Project1" {
		t.Errorf("expected first project to be 'Project1', got %s", projects[0])
	}
}

func testProjectConfig(t *testing.T) {
	t.Helper()

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var config map[string]interface{}
	err = client.GetJSON(context.Background(), utils.EndpointProjects+"/Project1", &config)
	if err != nil {
		t.Fatalf("failed to fetch project config: %v", err)
	}

	if config["cohortSelector"] == nil {
		t.Error("expected cohortSelector in config")
	}
}

func testProcessList(t *testing.T) {
	t.Helper()

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var processes []map[string]interface{}
	err = client.GetJSON(context.Background(), utils.EndpointProcessStatuses, &processes)
	if err != nil {
		t.Fatalf("failed to fetch process list: %v", err)
	}

	if len(processes) != 2 {
		t.Errorf("expected 2 processes, got %d", len(processes))
	}
}

func testProcessStatus(t *testing.T) {
	t.Helper()

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var process map[string]interface{}
	err = client.GetJSON(context.Background(), utils.EndpointProcessStatus("proc-123"), &process)
	if err != nil {
		t.Fatalf("failed to fetch process status: %v", err)
	}

	if process["processId"] != "proc-123" {
		t.Errorf("expected processId 'proc-123', got %v", process["processId"])
	}
	if process["phase"] != "RUNNING" {
		t.Errorf("expected phase 'RUNNING', got %v", process["phase"])
	}
}

func testTransferStart(t *testing.T) {
	t.Helper()

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	// Use PostJSON for transfer start which returns 202 Accepted with no body
	buf := &bytes.Buffer{}
	err = client.PostJSON(context.Background(), utils.EndpointTransferStart("Project1"), nil, buf)
	if err != nil {
		t.Fatalf("failed to start transfer: %v", err)
	}
}
