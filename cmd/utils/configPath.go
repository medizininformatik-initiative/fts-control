package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	appName    = "ftsctl"
	configFile = "config.yaml"
)

// ConfigFile represents the structure of config.yaml
type ConfigFile struct {
	Api  ApiConfig   `yaml:"api"`
	Auth *AuthConfig `yaml:"auth,omitempty"`
}

// ApiConfig represents the api section of config
type ApiConfig struct {
	BaseUrl string `yaml:"base_url"`
}

// AuthConfig represents authentication configuration.
// Only one authentication method should be configured at a time.
type AuthConfig struct {
	Basic       *BasicAuthConfig       `yaml:"basic,omitempty"`
	OAuth2      *OAuth2Config          `yaml:"oauth2,omitempty"`
	Certificate *CertificateAuthConfig `yaml:"certificate,omitempty"`
}

// BasicAuthConfig represents HTTP Basic authentication credentials.
type BasicAuthConfig struct {
	Username string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"` // Supports ${VAR} syntax for env vars
}

// OAuth2Config represents OAuth2 client credentials configuration.
type OAuth2Config struct {
	TokenURL     string `yaml:"token_url" mapstructure:"token_url"`
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"` // Supports ${VAR} syntax for env vars
	Scope        string `yaml:"scope,omitempty" mapstructure:"scope"`       // Optional scope parameter
}

// CertificateAuthConfig represents mutual TLS authentication configuration.
type CertificateAuthConfig struct {
	CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	CAFile   string `yaml:"ca_file,omitempty" mapstructure:"ca_file"` // Optional CA certificate
}

// GetConfigDir returns the OS-specific config directory for ftsctl
func GetConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(userConfigDir, appName), nil
}

// GetConfigPath returns the full path to config.yaml
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFile), nil
}

// ConfigExists checks if config.yaml exists
func ConfigExists() bool {
	configPath, err := GetConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(configPath)
	return err == nil
}

// ValidateURL validates URL format and returns normalized URL
func ValidateURL(urlString string) (string, error) {
	if urlString == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedUrl.Scheme == "" {
		return "", fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return "", fmt.Errorf("URL scheme must be http or https, got: %s", parsedUrl.Scheme)
	}

	if parsedUrl.Host == "" {
		return "", fmt.Errorf("URL must include a host")
	}

	return parsedUrl.String(), nil
}

// WriteConfig writes config.yaml with the given base URL
func WriteConfig(baseUrl string) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, configFile)

	config := ConfigFile{
		Api: ApiConfig{
			BaseUrl: baseUrl,
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
