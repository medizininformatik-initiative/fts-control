package utils

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
