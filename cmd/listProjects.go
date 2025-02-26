package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/url"
)

// listProjectsCmd represents the listProjects command
var listProjectsCmd = &cobra.Command{
	Use:   "listProjects", // /api/v2/projects
	Short: "List available projects",
	Long: `Cheese triangles cheesy feet emmental.
Emmental swiss cut the cheese gouda parmesan monterey jack lancashire cheeseburger. 
Macaroni cheese chalk and cheese jarlsberg ricotta cow fondue ricotta mascarpone. 
Port-salut hard cheese caerphilly babybel lancashire melted cheese.`,

	Run: func(cmd *cobra.Command, args []string) {

		// Import baseUrl from config.yaml
		baseUrlString := viper.GetString("api.base_url")
		if baseUrlString == "" {
			log.Fatal("Base API URL is not set in the configuration file")
		}

		// Parse baseUrl into url.URL
		baseUrl, err := url.Parse(baseUrlString)
		if err != nil {
			log.Fatalf("Invalid URL format: %v", err)
		}

		// Add the API endpoint to baseUrl
		apiUrl := baseUrl.ResolveReference(&url.URL{Path: "/api/v2/projects"}).String()

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
		fmt.Printf("List of all available projects:\n---------------------------------\n")
		for i, sliceProjects := range sliceProjects {
			fmt.Printf("%d. %s\n", i+1, sliceProjects)
		}
		fmt.Println("---------------------------------")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listProjectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listProjectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
