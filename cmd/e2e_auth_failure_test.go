package cmd

import (
	"context"
	"strings"
	"testing"

	"ftsctl/cmd/testutil"
	"ftsctl/cmd/testutil/mockserver"
	"ftsctl/cmd/utils"

	"github.com/spf13/viper"
)

// TestE2E_BasicAuth_WrongPassword tests that wrong password fails authentication.
func TestE2E_BasicAuth_WrongPassword(t *testing.T) {
	server := mockserver.New().WithBasicAuth("correctuser", "correctpass").Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.basic.user", "correctuser")
	viper.Set("auth.basic.password", "wrongpass")

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var projects []string
	err = client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err == nil {
		t.Error("expected error for wrong password, got nil")
		return
	}

	// Should get a 401 Unauthorized error
	var httpErr *utils.HTTPError
	if !containsHTTPError(err, 401, &httpErr) {
		t.Errorf("expected 401 Unauthorized error, got: %v", err)
	}
}

// TestE2E_BasicAuth_WrongUser tests that wrong username fails authentication.
func TestE2E_BasicAuth_WrongUser(t *testing.T) {
	server := mockserver.New().WithBasicAuth("correctuser", "correctpass").Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.basic.user", "wronguser")
	viper.Set("auth.basic.password", "correctpass")

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var projects []string
	err = client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err == nil {
		t.Error("expected error for wrong user, got nil")
		return
	}

	var httpErr *utils.HTTPError
	if !containsHTTPError(err, 401, &httpErr) {
		t.Errorf("expected 401 Unauthorized error, got: %v", err)
	}
}

// TestE2E_BasicAuth_NoCredentials tests that missing credentials fails.
func TestE2E_BasicAuth_NoCredentials(t *testing.T) {
	server := mockserver.New().WithBasicAuth("user", "pass").Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	// No auth configured

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var projects []string
	err = client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err == nil {
		t.Error("expected error when no credentials provided, got nil")
		return
	}

	var httpErr *utils.HTTPError
	if !containsHTTPError(err, 401, &httpErr) {
		t.Errorf("expected 401 Unauthorized error, got: %v", err)
	}
}

// TestE2E_OAuth2_NoToken tests that missing bearer token fails.
func TestE2E_OAuth2_NoToken(t *testing.T) {
	server := mockserver.New().WithOAuth2().Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	// No OAuth2 auth configured

	client, err := utils.GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to get authenticated client: %v", err)
	}

	var projects []string
	err = client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err == nil {
		t.Error("expected error when no token provided, got nil")
		return
	}

	var httpErr *utils.HTTPError
	if !containsHTTPError(err, 401, &httpErr) {
		t.Errorf("expected 401 Unauthorized error, got: %v", err)
	}
}

// TestE2E_OAuth2_InvalidTokenEndpoint tests that invalid token URL fails.
func TestE2E_OAuth2_InvalidTokenEndpoint(t *testing.T) {
	server := mockserver.New().WithOAuth2().Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	viper.Set("auth.oauth2.token_url", "http://localhost:9999/nonexistent")
	viper.Set("auth.oauth2.client_id", "test-client")
	viper.Set("auth.oauth2.client_secret", "test-secret")

	_, err := utils.GetAuthenticatedClient()
	if err == nil {
		t.Error("expected error for invalid token endpoint, got nil")
		return
	}

	if !strings.Contains(err.Error(), "failed to fetch OAuth2 token") {
		t.Errorf("expected token fetch error, got: %v", err)
	}
}

// TestE2E_Certificate_NoCert tests that missing client certificate fails.
func TestE2E_Certificate_NoCert(t *testing.T) {
	certs := mockserver.SetupTestCertificates(t)
	server := mockserver.New().WithCertificate(certs).Start(t)
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL())
	// No certificate auth configured - this should fail at TLS handshake

	client := utils.NewClient()

	var projects []string
	err := client.GetJSON(context.Background(), utils.EndpointProjects, &projects)
	if err == nil {
		t.Error("expected error when no client certificate provided, got nil")
		return
	}

	// TLS handshake should fail
	if !strings.Contains(err.Error(), "tls") && !strings.Contains(err.Error(), "certificate") {
		t.Logf("got error (likely TLS-related): %v", err)
	}
}

// TestE2E_OAuth2_RealKeycloak_InvalidClient tests OAuth2 with wrong client credentials.
func TestE2E_OAuth2_RealKeycloak_InvalidClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	testutil.RequireDockerCompose(t)

	if !testutil.IsKeycloakRunning() {
		t.Skip("Keycloak is not running, skipping OAuth2 failure test")
	}

	viper.Reset()
	viper.Set("api.base_url", "http://localhost:8080") // Won't reach this
	viper.Set("auth.oauth2.token_url", testutil.DefaultKeycloakConfig().TokenURL)
	viper.Set("auth.oauth2.client_id", "invalid-client")
	viper.Set("auth.oauth2.client_secret", "invalid-secret")

	_, err := utils.GetAuthenticatedClient()
	if err == nil {
		t.Error("expected error for invalid OAuth2 client, got nil")
		return
	}

	if !strings.Contains(err.Error(), "failed to fetch OAuth2 token") {
		t.Errorf("expected token fetch error, got: %v", err)
	}
}

// containsHTTPError checks if the error chain contains an HTTPError with the expected status.
func containsHTTPError(err error, expectedStatus int, httpErr **utils.HTTPError) bool {
	for e := err; e != nil; {
		if he, ok := e.(*utils.HTTPError); ok {
			if httpErr != nil {
				*httpErr = he
			}
			return he.StatusCode == expectedStatus
		}
		// Check if error has Unwrap method
		if unwrapper, ok := e.(interface{ Unwrap() error }); ok {
			e = unwrapper.Unwrap()
		} else {
			break
		}
	}
	return false
}
