// Package mockserver provides a mock CDA server for e2e testing.
package mockserver

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// AuthMode represents the authentication mode for the mock server.
type AuthMode int

const (
	// AuthNone means no authentication is required.
	AuthNone AuthMode = iota
	// AuthBasic requires HTTP Basic authentication.
	AuthBasic
	// AuthOAuth2 requires a Bearer token in the Authorization header.
	AuthOAuth2
	// AuthCertificate requires mTLS client certificate authentication.
	AuthCertificate
)

// MockCDAServer is a mock CDA server for testing ftsctl commands.
type MockCDAServer struct {
	server    *httptest.Server
	authMode  AuthMode
	basicUser string
	basicPass string
	certs     *TestCertificates
	mux       *http.ServeMux
	t         testing.TB
}

// New creates a new MockCDAServer with no authentication.
func New() *MockCDAServer {
	return &MockCDAServer{
		authMode: AuthNone,
		mux:      http.NewServeMux(),
	}
}

// WithAuthNone configures the server to accept requests without authentication.
func (s *MockCDAServer) WithAuthNone() *MockCDAServer {
	s.authMode = AuthNone
	return s
}

// WithBasicAuth configures the server to require HTTP Basic authentication.
func (s *MockCDAServer) WithBasicAuth(user, pass string) *MockCDAServer {
	s.authMode = AuthBasic
	s.basicUser = user
	s.basicPass = pass
	return s
}

// WithOAuth2 configures the server to require a Bearer token.
// The server only validates that a non-empty Bearer token is present,
// it does NOT validate the token content or issuer.
func (s *MockCDAServer) WithOAuth2() *MockCDAServer {
	s.authMode = AuthOAuth2
	return s
}

// WithCertificate configures the server to require mTLS client certificate authentication.
func (s *MockCDAServer) WithCertificate(certs *TestCertificates) *MockCDAServer {
	s.authMode = AuthCertificate
	s.certs = certs
	return s
}

// Start starts the mock server and registers all API endpoints.
// Call Close() when done to clean up resources.
func (s *MockCDAServer) Start(t testing.TB) *MockCDAServer {
	t.Helper()
	s.t = t

	// Register API endpoints with auth middleware
	s.registerHandlers()

	handler := s.wrapWithAuth(s.mux)

	if s.authMode == AuthCertificate {
		s.startTLSServer(handler)
	} else {
		s.server = httptest.NewServer(handler)
	}

	return s
}

// startTLSServer starts an HTTPS server with client certificate verification.
func (s *MockCDAServer) startTLSServer(handler http.Handler) {
	if s.certs == nil {
		s.t.Fatal("certificate auth requires test certificates")
	}

	// Load CA for verifying client certificates
	caCert, err := os.ReadFile(s.certs.CACertFile)
	if err != nil {
		s.t.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Load server certificate
	serverCert, err := tls.LoadX509KeyPair(s.certs.ServerCertFile, s.certs.ServerKeyFile)
	if err != nil {
		s.t.Fatalf("failed to load server certificate: %v", err)
	}

	server := httptest.NewUnstartedServer(handler)
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	server.StartTLS()
	s.server = server
}

// Close stops the mock server and releases resources.
func (s *MockCDAServer) Close() {
	if s.server != nil {
		s.server.Close()
	}
}

// URL returns the base URL of the mock server.
func (s *MockCDAServer) URL() string {
	if s.server == nil {
		return ""
	}
	return s.server.URL
}
