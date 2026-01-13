package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestExpandEnvVars(t *testing.T) {
	// Set up test environment variables using t.Setenv (automatically cleaned up)
	t.Setenv("TEST_VAR", "secret_value")
	t.Setenv("ANOTHER_VAR", "another_secret")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no variables",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "single ${VAR} syntax",
			input:    "${TEST_VAR}",
			expected: "secret_value",
		},
		{
			name:     "single $VAR syntax",
			input:    "$TEST_VAR",
			expected: "secret_value",
		},
		{
			name:     "variable with surrounding text",
			input:    "prefix_${TEST_VAR}_suffix",
			expected: "prefix_secret_value_suffix",
		},
		{
			name:     "multiple variables",
			input:    "${TEST_VAR}:${ANOTHER_VAR}",
			expected: "secret_value:another_secret",
		},
		{
			name:     "undefined variable expands to empty",
			input:    "${UNDEFINED_VAR}",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("expandEnvVars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateAuthConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		config    *AuthConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "nil config is valid",
			config:  nil,
			wantErr: false,
		},
		{
			name:      "empty config returns error",
			config:    &AuthConfig{},
			wantErr:   true,
			errSubstr: "no authentication method configured",
		},
		{
			name: "multiple methods returns error",
			config: &AuthConfig{
				Basic: &BasicAuthConfig{
					Username: "user",
					Password: "pass",
				},
				OAuth2: &OAuth2Config{
					TokenURL:     "https://auth.example.com/token",
					ClientID:     "client",
					ClientSecret: "secret",
				},
			},
			wantErr:   true,
			errSubstr: "only one authentication method can be configured at a time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthConfig(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errSubstr)
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errSubstr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAuthConfig_BasicAuth(t *testing.T) {
	tests := []struct {
		name      string
		config    *AuthConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid basic auth",
			config: &AuthConfig{
				Basic: &BasicAuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
			wantErr: false,
		},
		{
			name: "missing username",
			config: &AuthConfig{
				Basic: &BasicAuthConfig{
					Password: "pass",
				},
			},
			wantErr:   true,
			errSubstr: "user is required",
		},
		{
			name: "missing password",
			config: &AuthConfig{
				Basic: &BasicAuthConfig{
					Username: "user",
				},
			},
			wantErr:   true,
			errSubstr: "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthConfig(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errSubstr)
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errSubstr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAuthConfig_BasicAuth_EnvVarExpansion(t *testing.T) {
	t.Setenv("TEST_PASSWORD", "expanded_secret")

	auth := &AuthConfig{
		Basic: &BasicAuthConfig{
			Username: "user",
			Password: "${TEST_PASSWORD}",
		},
	}

	err := ValidateAuthConfig(auth)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if auth.Basic.Password != "expanded_secret" {
		t.Errorf("password not expanded: got %q, want %q", auth.Basic.Password, "expanded_secret")
	}
}

func TestValidateAuthConfig_OAuth2(t *testing.T) {
	tests := []struct {
		name      string
		config    *AuthConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid oauth2",
			config: &AuthConfig{
				OAuth2: &OAuth2Config{
					TokenURL:     "https://auth.example.com/token",
					ClientID:     "client",
					ClientSecret: "secret",
				},
			},
			wantErr: false,
		},
		{
			name: "missing token_url",
			config: &AuthConfig{
				OAuth2: &OAuth2Config{
					ClientID:     "client",
					ClientSecret: "secret",
				},
			},
			wantErr:   true,
			errSubstr: "token_url is required",
		},
		{
			name: "missing client_id",
			config: &AuthConfig{
				OAuth2: &OAuth2Config{
					TokenURL:     "https://auth.example.com/token",
					ClientSecret: "secret",
				},
			},
			wantErr:   true,
			errSubstr: "client_id is required",
		},
		{
			name: "missing client_secret",
			config: &AuthConfig{
				OAuth2: &OAuth2Config{
					TokenURL: "https://auth.example.com/token",
					ClientID: "client",
				},
			},
			wantErr:   true,
			errSubstr: "client_secret is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthConfig(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errSubstr)
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errSubstr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateAuthConfig_Certificate(t *testing.T) {
	// Create temporary certificate files for testing
	tempDir := t.TempDir()
	certFile := filepath.Join(tempDir, "client.crt")
	keyFile := filepath.Join(tempDir, "client.key")

	// Generate test certificate
	if err := generateTestCertificate(certFile, keyFile); err != nil {
		t.Fatalf("failed to generate test certificate: %v", err)
	}

	tests := []struct {
		name      string
		config    *AuthConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid certificate",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{
					CertFile: certFile,
					KeyFile:  keyFile,
				},
			},
			wantErr: false,
		},
		{
			name: "missing cert_file",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{
					KeyFile: keyFile,
				},
			},
			wantErr:   true,
			errSubstr: "cert_file is required",
		},
		{
			name: "missing key_file",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{
					CertFile: certFile,
				},
			},
			wantErr:   true,
			errSubstr: "key_file is required",
		},
		{
			name: "cert file not found",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{
					CertFile: "/nonexistent/cert.crt",
					KeyFile:  keyFile,
				},
			},
			wantErr:   true,
			errSubstr: "certificate file not found",
		},
		{
			name: "key file not found",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{
					CertFile: certFile,
					KeyFile:  "/nonexistent/key.key",
				},
			},
			wantErr:   true,
			errSubstr: "key file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthConfig(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errSubstr)
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("expected error containing %q, got %q", tt.errSubstr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAuthMethod(t *testing.T) {
	tests := []struct {
		name     string
		config   *AuthConfig
		expected AuthMethod
	}{
		{
			name:     "nil config",
			config:   nil,
			expected: AuthMethodNone,
		},
		{
			name:     "empty config",
			config:   &AuthConfig{},
			expected: AuthMethodNone,
		},
		{
			name: "basic auth",
			config: &AuthConfig{
				Basic: &BasicAuthConfig{Username: "user", Password: "pass"},
			},
			expected: AuthMethodBasic,
		},
		{
			name: "oauth2",
			config: &AuthConfig{
				OAuth2: &OAuth2Config{TokenURL: "url", ClientID: "id", ClientSecret: "secret"},
			},
			expected: AuthMethodOAuth2,
		},
		{
			name: "certificate",
			config: &AuthConfig{
				Certificate: &CertificateAuthConfig{CertFile: "cert", KeyFile: "key"},
			},
			expected: AuthMethodCertificate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAuthMethod(tt.config)
			if result != tt.expected {
				t.Errorf("GetAuthMethod() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAuthMethodString(t *testing.T) {
	tests := []struct {
		method   AuthMethod
		expected string
	}{
		{AuthMethodNone, "none"},
		{AuthMethodBasic, "basic"},
		{AuthMethodOAuth2, "oauth2"},
		{AuthMethodCertificate, "certificate"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.method.String() != tt.expected {
				t.Errorf("AuthMethod.String() = %q, want %q", tt.method.String(), tt.expected)
			}
		})
	}
}

func TestFetchOAuth2Token(t *testing.T) {
	// Set up mock OAuth2 server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		if r.FormValue("grant_type") != "client_credentials" {
			t.Errorf("expected grant_type=client_credentials, got %s", r.FormValue("grant_type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token": "test_token", "token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer server.Close()

	config := &OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "test_client",
		ClientSecret: "test_secret",
	}

	token, err := fetchOAuth2Token(config)
	if err != nil {
		t.Fatalf("fetchOAuth2Token() error = %v", err)
	}

	if token != "test_token" {
		t.Errorf("fetchOAuth2Token() = %q, want %q", token, "test_token")
	}
}

func TestFetchOAuth2Token_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "invalid_client"}`))
	}))
	defer server.Close()

	config := &OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "wrong_client",
		ClientSecret: "wrong_secret",
	}

	_, err := fetchOAuth2Token(config)
	if err == nil {
		t.Error("fetchOAuth2Token() expected error for 401 response")
	}
}

func TestGetAuthenticatedClient_NoAuth(t *testing.T) {
	// Reset viper
	viper.Reset()
	viper.Set("api.base_url", "http://localhost:8080")

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}
	if client == nil {
		t.Error("GetAuthenticatedClient() returned nil client")
	}
}

func TestGetAuthenticatedClient_BasicAuth(t *testing.T) {
	viper.Reset()
	viper.Set("api.base_url", "http://localhost:8080")
	viper.Set("auth.basic.user", "testuser")
	viper.Set("auth.basic.password", "testpass")

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}
	if client == nil {
		t.Error("GetAuthenticatedClient() returned nil client")
	}
}

func TestGetAuthenticatedClient_BasicAuth_SendsAuthorizationHeader(t *testing.T) {
	// Set up a test server that verifies the Authorization header
	var capturedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuthHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`["Project1"]`))
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL)
	viper.Set("auth.basic.user", "testuser")
	viper.Set("auth.basic.password", "testpass")

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}

	// Make a request to verify the header is sent
	var result []string
	err = client.GetJSON(context.Background(), "/api/v2/projects", &result)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	// Verify Authorization header
	// "testuser:testpass" in base64 is "dGVzdHVzZXI6dGVzdHBhc3M="
	expectedHeader := "Basic dGVzdHVzZXI6dGVzdHBhc3M="
	if capturedAuthHeader != expectedHeader {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, expectedHeader)
	}
}

func TestGetAuthenticatedClient_BasicAuth_WithEnvVarExpansion(t *testing.T) {
	t.Setenv("TEST_AUTH_PASSWORD", "secretpass")

	var capturedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuthHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL)
	viper.Set("auth.basic.user", "admin")
	viper.Set("auth.basic.password", "${TEST_AUTH_PASSWORD}")

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}

	var result []string
	err = client.GetJSON(context.Background(), "/api/v2/projects", &result)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	// "admin:secretpass" in base64 is "YWRtaW46c2VjcmV0cGFzcw=="
	expectedHeader := "Basic YWRtaW46c2VjcmV0cGFzcw=="
	if capturedAuthHeader != expectedHeader {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, expectedHeader)
	}
}

func TestGetAuthenticatedClient_BasicAuth_POST_SendsAuthorizationHeader(t *testing.T) {
	var capturedAuthHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuthHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL)
	viper.Set("auth.basic.user", "user")
	viper.Set("auth.basic.password", "pass")

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}

	// Make a POST request
	err = client.PostJSON(context.Background(), "/api/v2/process/start", nil, nil)
	if err != nil {
		t.Fatalf("PostJSON() error = %v", err)
	}

	// "user:pass" in base64 is "dXNlcjpwYXNz"
	expectedHeader := "Basic dXNlcjpwYXNz"
	if capturedAuthHeader != expectedHeader {
		t.Errorf("Authorization header = %q, want %q", capturedAuthHeader, expectedHeader)
	}
}

func TestGetAuthenticatedClient_NoAuth_NoAuthorizationHeader(t *testing.T) {
	var capturedAuthHeader string
	var headerPresent bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuthHeader = r.Header.Get("Authorization")
		_, headerPresent = r.Header["Authorization"]
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("api.base_url", server.URL)
	// No auth configured

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("GetAuthenticatedClient() error = %v", err)
	}

	var result []string
	err = client.GetJSON(context.Background(), "/api/v2/projects", &result)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}

	if headerPresent || capturedAuthHeader != "" {
		t.Errorf("Authorization header should not be present when no auth configured, got: %q", capturedAuthHeader)
	}
}

// generateTestCertificate generates a self-signed certificate for testing
func generateTestCertificate(certFile, keyFile string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer func() { _ = certOut.Close() }()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return err
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer func() { _ = keyOut.Close() }()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return err
	}

	return nil
}
