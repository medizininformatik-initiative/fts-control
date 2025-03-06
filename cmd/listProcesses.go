package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"time"
)

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
		// Generate API URL
		apiUrl := BuildApiUrl("/api/v2/process/statuses")

		client := &http.Client{}
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Fatalf("Error creating the request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending the request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading the response: %v", err)
		}

		var processes []Process
		if err := json.Unmarshal(body, &processes); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		DivdlnL()
		fmt.Println("List of all transfer process statuses:")
		DivdlnL()
		for _, p := range processes {
			fmt.Printf("ProcessID: %s\n", p.ProcessID)
			fmt.Printf("Phase: %s\n", p.Phase)
			fmt.Printf("CreatedAt: %s\n", p.CreatedAt.ToTime().Local().Format(time.RFC1123))
			if p.FinishedAt != nil {
				fmt.Printf("FinishedAt: %s\n", p.FinishedAt.ToTime().Local().Format(time.RFC1123))
			} else {
				fmt.Println("FinishedAt: The process is still running.")
			}
			fmt.Printf("TotalPatients: %d\n", p.TotalPatients)
			fmt.Printf("TotalBundles: %d\n", p.TotalBundles)
			fmt.Printf("DeidentifiedBundles: %d\n", p.DeidentifiedBundles)
			fmt.Printf("SentBundles: %d\n", p.SentBundles)
			fmt.Printf("SkippedBundles: %d\n", p.SkippedBundles)
			DivdlnS()
		}
	},
}

func init() {
	rootCmd.AddCommand(listProcessesCmd)

}
