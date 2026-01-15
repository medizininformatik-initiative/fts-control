package testutil

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// KeycloakConfig contains the Keycloak configuration for tests.
type KeycloakConfig struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
}

// DefaultKeycloakConfig returns the default Keycloak configuration
// matching the fts-realm.json setup.
func DefaultKeycloakConfig() *KeycloakConfig {
	return &KeycloakConfig{
		TokenURL:     "http://localhost:8080/realms/fts/protocol/openid-connect/token",
		ClientID:     "cd-client",
		ClientSecret: "tIQfOvBuhyR1dw9OQ3E4tCeTvcHtiW84",
	}
}

// IsKeycloakRunning checks if Keycloak is already running and healthy.
// Keycloak exposes health on port 9000 in dev mode.
func IsKeycloakRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:9000/health")
	if err != nil {
		return false
	}

	var health struct {
		Status string `json:"status"`
	}
	decodeErr := json.NewDecoder(resp.Body).Decode(&health)
	_ = resp.Body.Close()
	if decodeErr != nil {
		return false
	}

	return health.Status == "UP"
}

// RequireDockerCompose skips the test if docker compose is not available.
func RequireDockerCompose(t testing.TB) {
	t.Helper()

	cmd := exec.Command("docker", "compose", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("docker compose not available: %v", err)
	}

	if !strings.Contains(string(output), "Docker Compose") {
		t.Skipf("docker compose not available: unexpected output: %s", output)
	}
}
