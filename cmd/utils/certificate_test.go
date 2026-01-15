package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// TestCertificateAuthConfig holds all certificate files for a test scenario.
type TestCertificateAuthConfig struct {
	CACertFile     string
	CAKeyFile      string
	ServerCertFile string
	ServerKeyFile  string
	ClientCertFile string
	ClientKeyFile  string
}

// generateTestCA generates a CA certificate and key for testing.
func generateTestCA(certFile, keyFile string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "Test CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	if err := writeCertAndKey(certFile, keyFile, derBytes, priv); err != nil {
		return err
	}

	return nil
}

// generateTestCertSignedByCA generates a certificate signed by the given CA.
func generateTestCertSignedByCA(certFile, keyFile, caCertFile, caKeyFile string, isServer bool) error {
	// Load CA
	caCertPEM, err := os.ReadFile(caCertFile)
	if err != nil {
		return err
	}
	caKeyPEM, err := os.ReadFile(caKeyFile)
	if err != nil {
		return err
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return err
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return err
	}

	// Generate new key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	if isServer {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priv.PublicKey, caKey)
	if err != nil {
		return err
	}

	return writeCertAndKey(certFile, keyFile, derBytes, priv)
}

func writeCertAndKey(certFile, keyFile string, certDER []byte, key *rsa.PrivateKey) error {
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer func() { _ = certOut.Close() }()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer func() { _ = keyOut.Close() }()
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return err
	}

	return nil
}

// setupTestCertificates creates a full PKI for testing mTLS.
func setupTestCertificates(t *testing.T) *TestCertificateAuthConfig {
	t.Helper()
	tempDir := t.TempDir()

	cfg := &TestCertificateAuthConfig{
		CACertFile:     filepath.Join(tempDir, "ca.crt"),
		CAKeyFile:      filepath.Join(tempDir, "ca.key"),
		ServerCertFile: filepath.Join(tempDir, "server.crt"),
		ServerKeyFile:  filepath.Join(tempDir, "server.key"),
		ClientCertFile: filepath.Join(tempDir, "client.crt"),
		ClientKeyFile:  filepath.Join(tempDir, "client.key"),
	}

	// Generate CA
	if err := generateTestCA(cfg.CACertFile, cfg.CAKeyFile); err != nil {
		t.Fatalf("failed to generate CA: %v", err)
	}

	// Generate server certificate signed by CA
	if err := generateTestCertSignedByCA(cfg.ServerCertFile, cfg.ServerKeyFile, cfg.CACertFile, cfg.CAKeyFile, true); err != nil {
		t.Fatalf("failed to generate server certificate: %v", err)
	}

	// Generate client certificate signed by CA
	if err := generateTestCertSignedByCA(cfg.ClientCertFile, cfg.ClientKeyFile, cfg.CACertFile, cfg.CAKeyFile, false); err != nil {
		t.Fatalf("failed to generate client certificate: %v", err)
	}

	return cfg
}

func TestValidateCertificateAuthRequiresHTTPS(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid https URL",
			baseURL: "https://example.com",
			wantErr: false,
		},
		{
			name:      "http URL rejected",
			baseURL:   "http://example.com",
			wantErr:   true,
			errSubstr: "requires HTTPS",
		},
		{
			name:      "empty URL rejected",
			baseURL:   "",
			wantErr:   true,
			errSubstr: "api.base_url is required",
		},
		{
			name:      "invalid URL rejected",
			baseURL:   "://invalid",
			wantErr:   true,
			errSubstr: "invalid api.base_url",
		},
		{
			name:    "https with port",
			baseURL: "https://example.com:8443",
			wantErr: false,
		},
		{
			name:    "https with path",
			baseURL: "https://example.com/api/v1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCertificateAuthRequiresHTTPS(tt.baseURL)
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

func TestValidateCertificateAuth_CAFile(t *testing.T) {
	certs := setupTestCertificates(t)

	tests := []struct {
		name      string
		config    *CertificateAuthConfig
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid with CA file",
			config: &CertificateAuthConfig{
				CertFile: certs.ClientCertFile,
				KeyFile:  certs.ClientKeyFile,
				CAFile:   certs.CACertFile,
			},
			wantErr: false,
		},
		{
			name: "CA file not found",
			config: &CertificateAuthConfig{
				CertFile: certs.ClientCertFile,
				KeyFile:  certs.ClientKeyFile,
				CAFile:   "/nonexistent/ca.crt",
			},
			wantErr:   true,
			errSubstr: "CA certificate file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCertificateAuth(tt.config)
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

func TestValidateCertificateAuth_InvalidCAFile(t *testing.T) {
	certs := setupTestCertificates(t)
	tempDir := t.TempDir()

	// Create an invalid CA file (not a valid PEM)
	invalidCAFile := filepath.Join(tempDir, "invalid_ca.crt")
	if err := os.WriteFile(invalidCAFile, []byte("not a valid certificate"), 0600); err != nil {
		t.Fatalf("failed to write invalid CA file: %v", err)
	}

	config := &CertificateAuthConfig{
		CertFile: certs.ClientCertFile,
		KeyFile:  certs.ClientKeyFile,
		CAFile:   invalidCAFile,
	}

	err := validateCertificateAuth(config)
	if err == nil {
		t.Error("expected error for invalid CA file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse CA certificate") {
		t.Errorf("expected error about parsing CA certificate, got: %v", err)
	}
}

func TestValidateCertificateAuth_InvalidCertKeyPair(t *testing.T) {
	// Use setupTestCertificates to create two separate CA-signed certificates
	// so we can test with mismatched cert/key pairs
	certs1 := setupTestCertificates(t)
	certs2 := setupTestCertificates(t)

	// Use cert from one pair with key from another
	config := &CertificateAuthConfig{
		CertFile: certs1.ClientCertFile,
		KeyFile:  certs2.ClientKeyFile,
	}

	err := validateCertificateAuth(config)
	if err == nil {
		t.Error("expected error for mismatched cert/key pair, got nil")
	}
	if !strings.Contains(err.Error(), "failed to load certificate") {
		t.Errorf("expected error about loading certificate, got: %v", err)
	}
}

func TestApplyCertificateAuth(t *testing.T) {
	certs := setupTestCertificates(t)

	t.Run("preserves client headers", func(t *testing.T) {
		// Create a client with custom headers
		baseClient := NewClient().WithHeader("X-Custom-Header", "custom-value")

		config := &CertificateAuthConfig{
			CertFile: certs.ClientCertFile,
			KeyFile:  certs.ClientKeyFile,
		}

		client, err := applyCertificateAuth(baseClient, config)
		if err != nil {
			t.Fatalf("failed to apply certificate auth: %v", err)
		}
		if client == nil {
			t.Fatal("expected non-nil client")
		}

		// Verify headers are preserved
		if client.headers.Get("X-Custom-Header") != "custom-value" {
			t.Error("expected custom header to be preserved")
		}
	})

	t.Run("does not modify original client", func(t *testing.T) {
		baseClient := NewClient()
		originalHTTPClient := baseClient.httpClient

		config := &CertificateAuthConfig{
			CertFile: certs.ClientCertFile,
			KeyFile:  certs.ClientKeyFile,
		}

		newClient, err := applyCertificateAuth(baseClient, config)
		if err != nil {
			t.Fatalf("failed to apply certificate auth: %v", err)
		}

		// Original client should be unchanged
		if baseClient.httpClient != originalHTTPClient {
			t.Error("original client's httpClient should not be modified")
		}

		// New client should have a different httpClient
		if newClient.httpClient == originalHTTPClient {
			t.Error("new client should have a different httpClient with TLS config")
		}
	})
}

func TestCreateCertificateClient(t *testing.T) {
	certs := setupTestCertificates(t)

	t.Run("creates client with certificate", func(t *testing.T) {
		config := &CertificateAuthConfig{
			CertFile: certs.ClientCertFile,
			KeyFile:  certs.ClientKeyFile,
		}

		client, err := createCertificateClient(config)
		if err != nil {
			t.Fatalf("failed to create certificate client: %v", err)
		}
		if client == nil {
			t.Error("expected non-nil client")
		}
	})

	t.Run("creates client with certificate and CA", func(t *testing.T) {
		config := &CertificateAuthConfig{
			CertFile: certs.ClientCertFile,
			KeyFile:  certs.ClientKeyFile,
			CAFile:   certs.CACertFile,
		}

		client, err := createCertificateClient(config)
		if err != nil {
			t.Fatalf("failed to create certificate client: %v", err)
		}
		if client == nil {
			t.Error("expected non-nil client")
		}
	})
}

func TestGetAuthenticatedClient_Certificate_RequiresHTTPS(t *testing.T) {
	certs := setupTestCertificates(t)

	viper.Reset()
	viper.Set("api.base_url", "http://example.com") // HTTP, not HTTPS
	viper.Set("auth.certificate.cert_file", certs.ClientCertFile)
	viper.Set("auth.certificate.key_file", certs.ClientKeyFile)

	_, err := GetAuthenticatedClient()
	if err == nil {
		t.Error("expected error for HTTP with certificate auth, got nil")
	}
	if !strings.Contains(err.Error(), "requires HTTPS") {
		t.Errorf("expected error about requiring HTTPS, got: %v", err)
	}
}

func TestGetAuthenticatedClient_Certificate_AcceptsHTTPS(t *testing.T) {
	certs := setupTestCertificates(t)

	viper.Reset()
	viper.Set("api.base_url", "https://example.com")
	viper.Set("auth.certificate.cert_file", certs.ClientCertFile)
	viper.Set("auth.certificate.key_file", certs.ClientKeyFile)

	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// TestCertificateAuth_EndToEnd tests actual mTLS communication.
// This test creates a TLS server that requires client certificates and verifies
// that the certificate client can successfully communicate with it.
func TestCertificateAuth_EndToEnd(t *testing.T) {
	certs := setupTestCertificates(t)

	// Load CA for server to verify client certificates
	caCert, err := os.ReadFile(certs.CACertFile)
	if err != nil {
		t.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load server certificate
	serverCert, err := tls.LoadX509KeyPair(certs.ServerCertFile, certs.ServerKeyFile)
	if err != nil {
		t.Fatalf("failed to load server certificate: %v", err)
	}

	// Create TLS server that requires client certificates
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify client certificate was presented
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			t.Error("expected client certificate to be presented")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	server.StartTLS()
	defer server.Close()

	// Configure viper with the test server URL
	viper.Reset()
	viper.Set("api.base_url", server.URL)
	viper.Set("auth.certificate.cert_file", certs.ClientCertFile)
	viper.Set("auth.certificate.key_file", certs.ClientKeyFile)
	viper.Set("auth.certificate.ca_file", certs.CACertFile)

	// Create authenticated client
	client, err := GetAuthenticatedClient()
	if err != nil {
		t.Fatalf("failed to create authenticated client: %v", err)
	}

	// Make request to the TLS server
	var result map[string]string
	err = client.GetJSON(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", result)
	}
}

// TestCertificateAuth_EndToEnd_WithoutClientCert verifies that the server
// rejects requests without client certificates.
func TestCertificateAuth_EndToEnd_WithoutClientCert(t *testing.T) {
	certs := setupTestCertificates(t)

	// Load CA for server
	caCert, err := os.ReadFile(certs.CACertFile)
	if err != nil {
		t.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load server certificate
	serverCert, err := tls.LoadX509KeyPair(certs.ServerCertFile, certs.ServerKeyFile)
	if err != nil {
		t.Fatalf("failed to load server certificate: %v", err)
	}

	// Create TLS server that requires client certificates
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	server.StartTLS()
	defer server.Close()

	// Create a client WITHOUT certificate auth
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	// Request should fail due to missing client certificate
	_, err = httpClient.Get(server.URL + "/test")
	if err == nil {
		t.Error("expected error when connecting without client certificate")
	}
}

func TestLoadCACertPool_ValidCA(t *testing.T) {
	certs := setupTestCertificates(t)

	pool, err := loadCACertPool(certs.CACertFile)
	if err != nil {
		t.Fatalf("failed to load CA cert pool: %v", err)
	}
	if pool == nil {
		t.Error("expected non-nil cert pool")
	}
}

func TestLoadCACertPool_FileNotFound(t *testing.T) {
	_, err := loadCACertPool("/nonexistent/ca.crt")
	if err == nil {
		t.Error("expected error for nonexistent CA file")
	}
	if !strings.Contains(err.Error(), "failed to read CA certificate") {
		t.Errorf("expected error about reading CA certificate, got: %v", err)
	}
}

func TestLoadCACertPool_InvalidPEM(t *testing.T) {
	tempDir := t.TempDir()
	invalidCAFile := filepath.Join(tempDir, "invalid_ca.crt")
	if err := os.WriteFile(invalidCAFile, []byte("not a valid certificate"), 0600); err != nil {
		t.Fatalf("failed to write invalid CA file: %v", err)
	}

	_, err := loadCACertPool(invalidCAFile)
	if err == nil {
		t.Error("expected error for invalid PEM")
	}
	if !strings.Contains(err.Error(), "failed to parse CA certificate") {
		t.Errorf("expected error about parsing CA certificate, got: %v", err)
	}
}
