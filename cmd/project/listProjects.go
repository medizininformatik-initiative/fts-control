package project

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"ftsctl/cmd/utils"
	"github.com/spf13/cobra"
)

var ListProjectsCmd = &cobra.Command{
	Use:   "list", // /api/v2/projects
	Short: "List of all available projects",
	Long:  `List of all available projects from the API.`,

	Run: func(cmd *cobra.Command, args []string) {

		// Generate API URL
		apiUrl := utils.BuildApiUrl("/api/v2/projects")

		// Create HTTP client
		client := &http.Client{}

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			slog.Error("Error creating the request", "error", err)
			return
		}
		slog.Debug("Sending request", "url", req.URL.String())

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Error sending the request", "error", err)
			return
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				slog.Warn("Error closing response body", "error", err)
			}
		}()

		// Check the status code (StatusOk = 200)
		if resp.StatusCode != http.StatusOK {
			slog.Error("Unexpected status code", "status", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading the response", "error", err)
			return
		}

		var sliceProjects []string
		err = json.Unmarshal(body, &sliceProjects)
		if err != nil {
			slog.Error("Unmarshal error", "error", err)
			return
		}

		utils.DivdlnL()
		fmt.Println("List of all available projects:")
		utils.DivdlnL()
		for i, sliceProject := range sliceProjects {
			fmt.Printf("%d. %s\n", i+1, sliceProject)
		}
		utils.DivdlnS()
	},
}

func init() {

	Cmd.AddCommand(ListProjectsCmd)
}
