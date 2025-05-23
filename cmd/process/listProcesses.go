package process

import (
	"encoding/json"
	"fmt"
	"ftsctl/cmd/utils"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var listProcessesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transfer process statuses",
	Long:  `Lists all available transfer process statuses.`,
	Run: func(cmdCobra *cobra.Command, args []string) {
		apiUrl := utils.BuildApiUrl("/api/v2/process/statuses")

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
				log.Printf("Warning: failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading the response: %v", err)
		}

		var processes []utils.Process
		if err := json.Unmarshal(body, &processes); err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
		}

		utils.DivdlnL()
		fmt.Println("List of all transfer process statuses:")
		utils.DivdlnL()

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
			utils.DivdlnS()
		}
	},
}

func init() {
	Cmd.AddCommand(listProcessesCmd)
}
