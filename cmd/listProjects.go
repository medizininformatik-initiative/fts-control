package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

/*
// Project represents a single project.
type Project struct {
}

// ProjectsResponse response structure of the endpoint /api/v2/projects
type ProjectsResponse struct {
	Projects []Project `json:"projects"` // list of projects
}
*/
// listProjectsCmd represents the listProjects command
var listProjectsCmd = &cobra.Command{
	Use:   "listProjects", // /api/v2/projects
	Short: "List available projects",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// API endpoint
		url := "http://0.0.0.0:32772/api/v2/projects"

		// Create HTTP client
		client := &http.Client{}

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Error creating the request: %v", err)
		}

		// Optional: Add headers (e.g., Authorization, Content-Type, etc.)
		// req.Header.Add("Accept", "application/json")

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending the request: %v", err)
		}

		defer resp.Body.Close()

		// Check the status code
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading the response: %v", err)
		}

		// Print the response
		fmt.Printf("Response: %s\n", string(body))
	},
}

func init() {
	rootCmd.AddCommand(listProjectsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listProjectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listProjectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
