package cmd

import (
	"log"
	"net/url"

	"github.com/spf13/viper"
)

// Process struct for a single process
type Process struct {
	ProcessID           string     `json:"processId"`
	Phase               string     `json:"phase"`
	CreatedAt           TimeArray  `json:"createdAt"`
	FinishedAt          *TimeArray `json:"finishedAt"`
	TotalPatients       int        `json:"totalPatients"`
	TotalBundles        int        `json:"totalBundles"`
	DeidentifiedBundles int        `json:"deidentifiedBundles"`
	SentBundles         int        `json:"sentBundles"`
	SkippedBundles      int        `json:"skippedBundles"`
}

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
