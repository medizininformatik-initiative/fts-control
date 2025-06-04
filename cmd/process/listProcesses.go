package process

import (
	"encoding/json"
	"fmt"
	"ftsctl/cmd/utils"
	"io"
	"log/slog"
	"net/http"
	"os"
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
			slog.Error("Error creating the request", "error", err)
			os.Exit(1)
		}

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Error sending the request", "error", err)
			os.Exit(1)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Warn("Warning: failed to close response body", "error", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			slog.Error("Unexpected status code", "status", resp.StatusCode)
			os.Exit(1)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading the response", "error", err)
			os.Exit(1)
		}

		var processes []utils.Process
		if err := json.Unmarshal(body, &processes); err != nil {
			slog.Error("Error parsing JSON", "error", err)
			os.Exit(1)
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
