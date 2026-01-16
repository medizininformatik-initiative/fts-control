package mockserver

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// wrapWithAuth wraps the handler with authentication middleware based on authMode.
func (s *MockCDAServer) wrapWithAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch s.authMode {
		case AuthNone:
			// No authentication required
			handler.ServeHTTP(w, r)

		case AuthBasic:
			if !s.validateBasicAuth(r) {
				w.Header().Set("WWW-Authenticate", `Basic realm="ftsctl"`)
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
				return
			}
			handler.ServeHTTP(w, r)

		case AuthOAuth2:
			if !s.validateBearerToken(r) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error": "unauthorized"}`))
				return
			}
			handler.ServeHTTP(w, r)

		case AuthCertificate:
			// Certificate validation is handled at TLS layer.
			// If we reach here, the TLS handshake succeeded.
			handler.ServeHTTP(w, r)
		}
	})
}

// validateBasicAuth validates HTTP Basic authentication credentials.
func (s *MockCDAServer) validateBasicAuth(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return false
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return false
	}

	return creds[0] == s.basicUser && creds[1] == s.basicPass
}

// validateBearerToken validates that a Bearer token is present (non-empty).
// This does NOT validate the token content, issuer, or expiry.
func (s *MockCDAServer) validateBearerToken(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}

	token := auth[len(prefix):]
	return token != ""
}
