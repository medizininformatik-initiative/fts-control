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
	Api ApiConfig `yaml:"api"`
}

// ApiConfig represents the api section of config
type ApiConfig struct {
	BaseUrl string `yaml:"base_url"`
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
