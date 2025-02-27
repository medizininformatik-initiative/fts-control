package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// struct for a single Process
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

type TimeArray [7]int

// listProcessesCmd represents the listProcesses command
var listProcessesCmd = &cobra.Command{
	Use:   "listProcesses",
	Short: "List of all transfer process statuses",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listProcesses called")
	},
}

func init() {
	rootCmd.AddCommand(listProcessesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listProcessesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listProcessesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
