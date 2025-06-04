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

var processId string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show process status",
	Long:  `Shows the process status of the selected process ID`,
	Run: func(cmd *cobra.Command, args []string) {

		if !cmd.Flags().Changed("processId") {
			utils.DivdlnL()
			fmt.Printf("Warning: No --processId (-i) specified.\n\nPlease set the flag to retrieve the status of a process.\n")
			utils.DivdlnL()
			return
		}

		apiUrl := utils.BuildApiUrl("/api/v2/process/status/" + processId)
		fmt.Printf("API URL: %s\n", apiUrl)
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
				slog.Warn("Error closing response body", "error", err)
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

		prcss := utils.Process{}
		err = json.Unmarshal(body, &prcss)
		if err != nil {
			slog.Error("Error unmarshalling response", "error", err)
			os.Exit(1)
		}

		utils.DivdlnL()
		fmt.Printf("Representation of the Process with the ProcessID: %s \n", processId)
		utils.DivdlnL()
		fmt.Printf("ProcessID: %s\n", prcss.ProcessID)
		fmt.Printf("Phase: %s\n", prcss.Phase)
		fmt.Printf("CreatedAt: %s\n", prcss.CreatedAt.ToTime().Local().Format(time.RFC1123))
		if prcss.FinishedAt != nil {
			fmt.Printf("FinishedAt: %s\n", prcss.FinishedAt.ToTime().Local().Format(time.RFC1123))
		} else {
			fmt.Println("FinishedAt: The process is still running")
		}
		fmt.Printf("TotalPatients: %d\n", prcss.TotalPatients)
		fmt.Printf("TotalBundles: %d\n", prcss.TotalBundles)
		fmt.Printf("DeidentifiedBundles: %d\n", prcss.DeidentifiedBundles)
		fmt.Printf("SentBundles: %d\n", prcss.SentBundles)
		fmt.Printf("SkippedBundles: %d\n", prcss.SkippedBundles)
		utils.DivdlnS()
	},
}

func init() {
	StatusCmd.Flags().StringVarP(&processId, "processId", "i", "", "Please Enter the processId")
	Cmd.AddCommand(StatusCmd)
}
