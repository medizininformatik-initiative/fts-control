package utils

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestOAuth2Token_IsExpired(t *testing.T) {
	tests := []struct {
		name        string
		token       *OAuth2Token
		wantExpired bool
	}{
		{
			name:        "nil token is expired",
			token:       nil,
			wantExpired: true,
		},
		{
			name: "expired token",
			token: &OAuth2Token{
				AccessToken: "token",
				ExpiresAt:   time.Now().Add(-1 * time.Minute),
			},
			wantExpired: true,
		},
		{
			name: "token expiring within buffer (30s)",
			token: &OAuth2Token{
				AccessToken: "token",
				ExpiresAt:   time.Now().Add(20 * time.Second),
			},
			wantExpired: true,
		},
		{
			name: "valid token",
			token: &OAuth2Token{
				AccessToken: "token",
				ExpiresAt:   time.Now().Add(5 * time.Minute),
			},
			wantExpired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.IsExpired(); got != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExpired)
			}
		})
	}
}

func TestOAuth2TokenProvider_GetToken_Success(t *testing.T) {
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
		if r.FormValue("client_id") != "test_client" {
			t.Errorf("expected client_id=test_client, got %s", r.FormValue("client_id"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token": "test_token", "token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "test_client",
		ClientSecret: "test_secret",
	})

	token, err := provider.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if token != "test_token" {
		t.Errorf("GetToken() = %q, want %q", token, "test_token")
	}
}

func TestOAuth2TokenProvider_GetToken_WithScope(t *testing.T) {
	var receivedScope string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		receivedScope = r.FormValue("scope")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token": "scoped_token", "expires_in": 3600}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
		Scope:        "read write",
	})

	_, err := provider.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if receivedScope != "read write" {
		t.Errorf("scope = %q, want %q", receivedScope, "read write")
	}
}

func TestOAuth2TokenProvider_GetToken_NoScope(t *testing.T) {
	var scopePresent bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		_, scopePresent = r.Form["scope"]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token": "token", "expires_in": 3600}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
		// No scope
	})

	_, err := provider.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if scopePresent {
		t.Error("scope should not be present when not configured")
	}
}

func TestOAuth2TokenProvider_GetToken_Caching(t *testing.T) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token": "cached_token", "expires_in": 3600}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
	})

	// First call should fetch token
	token1, err := provider.GetToken()
	if err != nil {
		t.Fatalf("First GetToken() error = %v", err)
	}

	// Second call should return cached token
	token2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	// Third call should also return cached token
	token3, err := provider.GetToken()
	if err != nil {
		t.Fatalf("Third GetToken() error = %v", err)
	}

	if token1 != token2 || token2 != token3 {
		t.Errorf("tokens should be the same: %q, %q, %q", token1, token2, token3)
	}

	if count := atomic.LoadInt32(&requestCount); count != 1 {
		t.Errorf("expected 1 request (caching), got %d", count)
	}
}

func TestOAuth2TokenProvider_GetToken_RefreshOnExpiry(t *testing.T) {
	var requestCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return token that expires in 1 second
		_, _ = w.Write([]byte(`{"access_token": "token_` + string(rune('0'+count)) + `", "expires_in": 1}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
	})

	// First call
	_, err := provider.GetToken()
	if err != nil {
		t.Fatalf("First GetToken() error = %v", err)
	}

	// Wait for token to expire (expires_in=1s, plus 30s buffer means it's already expired)
	// The 30s buffer in IsExpired() means a 1s token is immediately considered expired
	// So second call should fetch a new token
	_, err = provider.GetToken()
	if err != nil {
		t.Fatalf("Second GetToken() error = %v", err)
	}

	if count := atomic.LoadInt32(&requestCount); count != 2 {
		t.Errorf("expected 2 requests (token refresh), got %d", count)
	}
}

func TestOAuth2TokenProvider_GetToken_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "invalid_client"}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "wrong_client",
		ClientSecret: "wrong_secret",
	})

	_, err := provider.GetToken()
	if err == nil {
		t.Error("expected error for 401 response")
	}
}

func TestOAuth2TokenProvider_GetToken_MissingAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token_type": "Bearer", "expires_in": 3600}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
	})

	_, err := provider.GetToken()
	if err == nil {
		t.Error("expected error for missing access_token")
	}
}

func TestOAuth2TokenProvider_GetToken_DefaultExpiresIn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// No expires_in field
		_, _ = w.Write([]byte(`{"access_token": "token"}`))
	}))
	defer server.Close()

	provider := NewOAuth2TokenProvider(&OAuth2Config{
		TokenURL:     server.URL,
		ClientID:     "client",
		ClientSecret: "secret",
	})

	token, err := provider.GetToken()
	if err != nil {
		t.Fatalf("GetToken() error = %v", err)
	}
	if token != "token" {
		t.Errorf("GetToken() = %q, want %q", token, "token")
	}

	// Token should not be expired (default 5 min expiry)
	if provider.token.IsExpired() {
		t.Error("token should not be expired with default expiry")
	}
}
