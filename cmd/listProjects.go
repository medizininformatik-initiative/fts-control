package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

// listProjectsCmd represents the listProjects command
var listProjectsCmd = &cobra.Command{
	Use:   "listProjects", // /api/v2/projects
	Short: "List of all available projects",
	Long:  `List of all available projects from the API.`,

	Run: func(cmd *cobra.Command, args []string) {

		// Generate API URL
		apiUrl := BuildApiUrl("/api/v2/projects")

		// Create HTTP client
		client := &http.Client{}

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Fatalf("Error creating the request: %v", err)
		}

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending the request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		// Check the status code (StatusOk = 200)
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected stats code: %d", resp.StatusCode)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading the response: %v", err)
		}

		//
		// json formating
		//

		// Slice for projects
		var sliceProjects []string

		// JSON unmarshal
		err = json.Unmarshal(body, &sliceProjects)
		if err != nil {
			fmt.Println("Unmarshal error:", err)
			return
		}

		// Output
		DivdlnL()
		fmt.Println("List of all available projects:")
		DivdlnL()
		for i, sliceProjects := range sliceProjects {
			fmt.Printf("%d. %s\n", i+1, sliceProjects)
		}
		DivdlnS()
	},
}

func init() {

	// load config.yaml
	viper.SetConfigName("config") // Filename
	viper.SetConfigType("yaml")   // Filetype
	viper.AddConfigPath("..")     // directory where the file is located

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	rootCmd.AddCommand(listProjectsCmd)

}
