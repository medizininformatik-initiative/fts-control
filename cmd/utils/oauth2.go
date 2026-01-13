package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// OAuth2Token represents a cached OAuth2 access token with expiry information.
type OAuth2Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

// IsExpired returns true if the token is expired or will expire within the buffer period.
// Uses a 30-second buffer to avoid using tokens that are about to expire.
func (t *OAuth2Token) IsExpired() bool {
	if t == nil {
		return true
	}
	return time.Now().Add(30 * time.Second).After(t.ExpiresAt)
}

// OAuth2TokenProvider manages OAuth2 token acquisition and caching.
// It fetches tokens using the client credentials grant and caches them
// until they expire.
type OAuth2TokenProvider struct {
	config     *OAuth2Config
	token      *OAuth2Token
	httpClient *http.Client
}

// NewOAuth2TokenProvider creates a new token provider for the given configuration.
func NewOAuth2TokenProvider(config *OAuth2Config) *OAuth2TokenProvider {
	return &OAuth2TokenProvider{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetToken returns a valid access token, fetching a new one if the cached token
// is expired or doesn't exist.
func (p *OAuth2TokenProvider) GetToken() (string, error) {
	if p.token != nil && !p.token.IsExpired() {
		return p.token.AccessToken, nil
	}

	token, err := p.fetchToken()
	if err != nil {
		return "", err
	}

	p.token = token
	return token.AccessToken, nil
}

// fetchToken requests a new access token from the OAuth2 token endpoint.
func (p *OAuth2TokenProvider) fetchToken() (*OAuth2Token, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", p.config.ClientID)
	data.Set("client_secret", p.config.ClientSecret)
	if p.config.Scope != "" {
		data.Set("scope", p.config.Scope)
	}

	resp, err := p.httpClient.PostForm(p.config.TokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("access_token not found in response")
	}

	// Default to 5 minutes if expires_in is not provided
	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 300
	}

	return &OAuth2Token{
		AccessToken: tokenResp.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}
