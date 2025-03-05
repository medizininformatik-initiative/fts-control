package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var processId string

// processStatusCmd represents the processStatus command
var processStatusCmd = &cobra.Command{
	Use:   "processStatus",
	Short: "Show process status",
	Long:  `Shows the process status of the selected process ID`,
	Run: func(cmd *cobra.Command, args []string) {

		apiUrl := BuildApiUrl("/api/v2/process/status/" + processId)

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

		//fmt.Println(string(body))

		prcss := Process{}
		err = json.Unmarshal(body, &prcss)
		if err != nil {
			log.Fatalf("Error unmarshalling response: %v", err)
		}

		fmt.Println("------------------------------------------------------------------")
		fmt.Printf("Representation of the Process with the ProcessID: %s \n", processId)
		fmt.Println("------------------------------------------------------------------")
		fmt.Printf("ProcessID: %s\n", prcss.ProcessID)
		fmt.Printf("Phase: %s\n", prcss.Phase)
		fmt.Printf("CreatedAt: %s\n", prcss.CreatedAt.ToTime().Local().Format(time.RFC1123))
		if prcss.FinishedAt != nil {
			fmt.Printf("FinishedAt: %s\n", prcss.FinishedAt.ToTime().Local().Format(time.RFC1123))
		} else {
			fmt.Println("FinishedAt: null")
		}
		fmt.Printf("TotalPatients: %d\n", prcss.TotalPatients)
		fmt.Printf("TotalBundles: %d\n", prcss.TotalBundles)
		fmt.Printf("DeidentifiedBundles: %d\n", prcss.DeidentifiedBundles)
		fmt.Printf("SentBundles: %d\n", prcss.SentBundles)
		fmt.Printf("SkippedBundles: %d\n", prcss.SkippedBundles)
		fmt.Println("------------------------------------------------------------------")
	},
}

func init() {
	processStatusCmd.Flags().StringVarP(&processId, "processId", "i", "", "Please Enter the processId")

	// Set Flag as Required
	if err := processStatusCmd.MarkFlagRequired("processId"); err != nil {
		log.Fatalf("Error marking processId as required: %v", err)
	}

	rootCmd.AddCommand(processStatusCmd)

}
