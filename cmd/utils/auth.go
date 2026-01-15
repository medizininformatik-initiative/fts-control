package utils

import (
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/spf13/viper"
)

// AuthMethod represents the type of authentication configured.
type AuthMethod int

const (
	AuthMethodNone AuthMethod = iota
	AuthMethodBasic
	AuthMethodOAuth2
	AuthMethodCertificate
)

// String returns the string representation of the auth method.
func (m AuthMethod) String() string {
	switch m {
	case AuthMethodBasic:
		return "basic"
	case AuthMethodOAuth2:
		return "oauth2"
	case AuthMethodCertificate:
		return "certificate"
	default:
		return "none"
	}
}

// envVarPattern matches ${VAR} or $VAR syntax for environment variables.
var envVarPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// expandEnvVars expands environment variables in the given string.
// Supports both ${VAR} and $VAR syntax.
// Returns the original string with all variables expanded.
func expandEnvVars(s string) string {
	return envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := extractVarName(match)
		return os.Getenv(varName)
	})
}

// extractVarName extracts the variable name from a match.
// For "${VAR}" format, returns the content between braces.
// For "$VAR" format, returns everything after the dollar sign.
func extractVarName(match string) string {
	if len(match) > 2 && match[1] == '{' {
		return match[2 : len(match)-1]
	}
	return match[1:]
}

// ValidateAuthConfig validates the authentication configuration.
// It ensures only one auth method is configured and all required fields are present.
//
// Note: This function mutates the config by expanding environment variables
// in secret fields (password, client_secret). The expanded values replace
// the original ${VAR} placeholders in the config struct.
func ValidateAuthConfig(auth *AuthConfig) error {
	if auth == nil {
		return nil
	}

	count := countConfiguredAuthMethods(auth)
	if count == 0 {
		return fmt.Errorf("auth section present but no authentication method configured")
	}
	if count > 1 {
		return fmt.Errorf("only one authentication method can be configured at a time")
	}

	if auth.Basic != nil {
		return validateBasicAuth(auth.Basic)
	}
	if auth.OAuth2 != nil {
		return validateOAuth2Auth(auth.OAuth2)
	}
	return validateCertificateAuth(auth.Certificate)
}

func countConfiguredAuthMethods(auth *AuthConfig) int {
	count := 0
	if auth.Basic != nil {
		count++
	}
	if auth.OAuth2 != nil {
		count++
	}
	if auth.Certificate != nil {
		count++
	}
	return count
}

func validateBasicAuth(basic *BasicAuthConfig) error {
	if basic.Username == "" {
		return fmt.Errorf("auth.basic.user is required")
	}
	if basic.Password == "" {
		return fmt.Errorf("auth.basic.password is required")
	}
	basic.Password = expandEnvVars(basic.Password)
	return nil
}

func validateOAuth2Auth(oauth2 *OAuth2Config) error {
	if oauth2.TokenURL == "" {
		return fmt.Errorf("auth.oauth2.token_url is required")
	}
	if oauth2.ClientID == "" {
		return fmt.Errorf("auth.oauth2.client_id is required")
	}
	if oauth2.ClientSecret == "" {
		return fmt.Errorf("auth.oauth2.client_secret is required")
	}
	if _, err := url.Parse(oauth2.TokenURL); err != nil {
		return fmt.Errorf("auth.oauth2.token_url is not a valid URL: %w", err)
	}
	oauth2.ClientSecret = expandEnvVars(oauth2.ClientSecret)
	return nil
}

func validateFileExists(path, description string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s not found: %s: %w", description, path, err)
	}
	return nil
}

// GetAuthMethod returns which authentication method is configured.
func GetAuthMethod(auth *AuthConfig) AuthMethod {
	if auth == nil {
		return AuthMethodNone
	}
	switch {
	case auth.Basic != nil:
		return AuthMethodBasic
	case auth.OAuth2 != nil:
		return AuthMethodOAuth2
	case auth.Certificate != nil:
		return AuthMethodCertificate
	default:
		return AuthMethodNone
	}
}

// GetAuthenticatedClient returns an HTTP client with authentication applied.
// If no auth is configured, returns a basic client.
func GetAuthenticatedClient() (*Client, error) {
	client := NewClient()

	if !viper.IsSet("auth") {
		return client, nil
	}

	var auth AuthConfig
	if err := viper.UnmarshalKey("auth", &auth); err != nil {
		return nil, fmt.Errorf("failed to parse auth configuration: %w", err)
	}

	method := GetAuthMethod(&auth)
	if method == AuthMethodNone {
		return client, nil
	}

	if err := ValidateAuthConfig(&auth); err != nil {
		return nil, fmt.Errorf("invalid auth configuration: %w", err)
	}

	return applyAuthentication(client, method, &auth)
}

func applyAuthentication(client *Client, method AuthMethod, auth *AuthConfig) (*Client, error) {
	switch method {
	case AuthMethodBasic:
		return client.WithBasicAuth(auth.Basic.Username, auth.Basic.Password), nil
	case AuthMethodOAuth2:
		tokenProvider := NewOAuth2TokenProvider(auth.OAuth2)
		token, err := tokenProvider.GetToken()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch OAuth2 token: %w", err)
		}
		return client.WithBearerToken(token), nil
	case AuthMethodCertificate:
		// Validate HTTPS is required for certificate authentication
		baseURL := viper.GetString("api.base_url")
		if err := validateCertificateAuthRequiresHTTPS(baseURL); err != nil {
			return nil, err
		}
		return applyCertificateAuth(client, auth.Certificate)
	}
	return client, nil
}
