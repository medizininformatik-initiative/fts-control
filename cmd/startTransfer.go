package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

// requestPayload defines the structure for the optional ID list
type requestPayload struct {
	IDs []string `json:"ids,omitempty"`
}

// startTransferCmd represents the startTransfer command
var startTransferCmd = &cobra.Command{
	Use:   "startTransfer", // /api/v2/process/{project}/start
	Short: "Start a transfer process",
	Long:  `Start a transfer of patients with IDs given in the request body or if empty start a transfer of all consented patients of the selected project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve project name
		projectName, _ := cmd.Flags().GetString("projectName")

		// Check if project name is provided
		if projectName == "" {
			fmt.Println("ERROR: The --projectName (-n) flag is required!")
			fmt.Println("Usage: ftsctl startTransfer --projectName exampleProject [--ids id1,id2,id3]")
			return
		}

		// Retrieve IDs (if provided)
		idsStr, _ := cmd.Flags().GetString("ids")
		var ids []string
		if idsStr != "" {
			ids = strings.Split(idsStr, ",")
		}

		// Build API URL
		apiUrl := BuildApiUrl(fmt.Sprintf("/api/v2/process/%s/start", projectName))

		// Generate payload
		payload := requestPayload{IDs: ids}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Error creating JSON payload: %v", err)
		}

		// Send HTTP request
		resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		// Check if response status is 200 OK
		if resp.StatusCode == http.StatusOK {

			if len(ids) > 0 {
				DivdlnL()
				fmt.Printf("Transfer of the project %s \nwith the given patient IDs has started.\n", projectName)
				DivdlnL()
				fmt.Println("Patient-IDs:")
				for _, id := range ids {
					fmt.Printf("   - %s\n", id)
				}
			} else {
				DivdlnL()
				fmt.Printf("Transfer of the project %s has started.\n", projectName)
				DivdlnL()
				fmt.Printf("No patient IDs provided. \nThe transfer starts with all patient IDs from %s.\n", projectName)
			}

			DivdlnS()
			fmt.Printf("\nApiUrl: %s\n", apiUrl)
			fmt.Printf("\nResponse Status: %s\n", resp.Status)
			DivdlnS()
		} else {
			DivdlnL()
			fmt.Printf("Request failed! Response Status: %s\n", resp.Status)
			DivdlnL()
		}
	},
}

func init() {
	startTransferCmd.Flags().StringP("projectName", "n", "", "Project name (required)")
	startTransferCmd.Flags().StringP("ids", "i", "", "Comma-separated list of patient IDs (optional)")

	rootCmd.AddCommand(startTransferCmd)
}
