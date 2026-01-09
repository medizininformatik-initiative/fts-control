package utils

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

// GetBaseURL returns the base URL from config.yaml or an error if not configured.
func GetBaseURL() (*url.URL, error) {
	baseUrlString := viper.GetString("api.base_url")
	if baseUrlString == "" {
		return nil, fmt.Errorf("base API URL is not set (use 'ftsctl config set-base-url' to configure)")
	}

	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	return baseUrl, nil
}

// BuildAPIURL combines the base URL with a specific API endpoint.
func BuildAPIURL(endpoint string) (string, error) {
	baseUrl, err := GetBaseURL()
	if err != nil {
		return "", err
	}
	return baseUrl.ResolveReference(&url.URL{Path: endpoint}).String(), nil
}
