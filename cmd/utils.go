package cmd

import (
	"log"
	"net/url"
	"time"

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

func (t TimeArray) ToTime() time.Time {
	return time.Date(t[0], time.Month(t[1]), t[2], t[3], t[4], t[5], t[6], time.UTC)
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
