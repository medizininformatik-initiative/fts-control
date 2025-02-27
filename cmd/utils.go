package cmd

import (
	"log"
	"net/url"

	"github.com/spf13/viper"
)

// GetBaseURL Import baseUrl from config.yaml
func GetBaseURL() *url.URL {
	baseUrlString := viper.GetString("api.base_url")
	if baseUrlString == "" {
		log.Fatal("Base API URL is not set in the configuration file")
	}

	// Parse baseUrl into url.URL
	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		log.Fatalf("Invalid URL format: %v", err)
	}

	return baseUrl
}

// BuildApiUrl combines the base URL with a specific API endpoint
func BuildApiUrl(endpoint string) string {
	return GetBaseURL().ResolveReference(&url.URL{Path: endpoint}).String()
}
