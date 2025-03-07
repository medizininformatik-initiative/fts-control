package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// Structs for projectConfig JSON

type Config struct {
	CohortSelector  CohortSelector  `json:"cohortSelector"`
	DataSelector    DataSelector    `json:"dataSelector"`
	Deidentificator Deidentificator `json:"deidentificator"`
	BundleSender    BundleSender    `json:"bundleSender"`
}

type CohortSelector struct {
	TrustCenterAgent TrustCenterAgent `json:"trustCenterAgent"`
}

type TrustCenterAgent struct {
	Server                  Server   `json:"server"`
	Domain                  string   `json:"domain"`
	PatientIdentifierSystem string   `json:"patientIdentifierSystem"`
	PolicySystem            string   `json:"policySystem"`
	Policies                []string `json:"policies"`
}

type Server struct {
	BaseUrl string `json:"baseUrl"`
}

type DataSelector struct {
	Everything Everything `json:"everything"`
}

type Everything struct {
	FhirServer Server  `json:"fhirServer"`
	Resolve    Resolve `json:"resolve"`
}

type Resolve struct {
	PatientIdentifierSystem string `json:"patientIdentifierSystem"`
}

type Deidentificator struct {
	Deidentifhir Deidentifhir `json:"deidentifhir"`
}

type Deidentifhir struct {
	TrustCenterAgent   TrustCenterAgentDeident `json:"trustCenterAgent"`
	MaxDateShift       string                  `json:"maxDateShift"`
	DeidentifhirConfig string                  `json:"deidentifhirConfig"`
	ScraperConfig      string                  `json:"scraperConfig"`
}

type TrustCenterAgentDeident struct {
	Server  Server  `json:"server"`
	Domains Domains `json:"domains"`
}

type Domains struct {
	Pseudonym string `json:"pseudonym"`
	Salt      string `json:"salt"`
	DateShift string `json:"dateShift"`
}

type BundleSender struct {
	ResearchDomainAgent ResearchDomainAgent `json:"researchDomainAgent"`
}

type ResearchDomainAgent struct {
	Server  Server `json:"server"`
	Project string `json:"project"`
}

// Name of the Project

var prjName string

// projectConfigCmd represents the projectConfig command
var projectConfigCmd = &cobra.Command{
	Use:   "projectConfig",
	Short: "Project Configuration",
	Long:  `Shows the project configuration of the selected project.`,
	Run: func(cmd *cobra.Command, args []string) {

		if !cmd.Flags().Changed("projectName") {
			DivdlnL()
			fmt.Printf("Warning: No --projectName (-n) specified.\n\nPlease set the Flag to retrieve the Project Configuration of a Project.\n")
			DivdlnL()
			return
		}

		apiUrl := BuildApiUrl("/api/v2/projects/" + prjName)
		//apiUrl := BuildApiUrl("/api/v2/projects/example")

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

		pConfig := Config{}
		err = json.Unmarshal(body, &pConfig)
		if err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}
		/////////////////////////
		// Formated Output
		DivdlnL()
		fmt.Printf("Project configuration with the Project Name: %s \n", prjName)
		DivdlnL()
		fmt.Printf("CohortSelector - TrustCenterAgent:\n")
		fmt.Printf("  Server BaseURL: %s\n", pConfig.CohortSelector.TrustCenterAgent.Server.BaseUrl)
		fmt.Printf("  Domain: %s\n", pConfig.CohortSelector.TrustCenterAgent.Domain)
		fmt.Printf("  Patient Identifier System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PatientIdentifierSystem)
		fmt.Printf("  Policy System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PolicySystem)
		fmt.Println("  Policies:")
		for _, policy := range pConfig.CohortSelector.TrustCenterAgent.Policies {
			fmt.Printf("    - %s\n", policy)
		}

		DivdlnS()
		fmt.Printf("DataSelector - Everything:\n")
		fmt.Printf("  FHIR Server BaseURL: %s\n", pConfig.DataSelector.Everything.FhirServer.BaseUrl)
		fmt.Printf("  Resolve Patient Identifier System: %s\n", pConfig.DataSelector.Everything.Resolve.PatientIdentifierSystem)

		DivdlnS()
		fmt.Printf("Deidentificator - Deidentifhir:\n")
		fmt.Printf("  TrustCenterAgent Server BaseURL: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Server.BaseUrl)
		fmt.Printf("  Max Date Shift: %s\n", pConfig.Deidentificator.Deidentifhir.MaxDateShift)
		fmt.Printf("  Deidentifhir Config: %s\n", pConfig.Deidentificator.Deidentifhir.DeidentifhirConfig)
		fmt.Printf("  Scraper Config: %s\n", pConfig.Deidentificator.Deidentifhir.ScraperConfig)
		fmt.Println("  Domains:")
		fmt.Printf("    Pseudonym: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Pseudonym)
		fmt.Printf("    Salt: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Salt)
		fmt.Printf("    DateShift: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.DateShift)

		DivdlnS()
		fmt.Printf("BundleSender - ResearchDomainAgent:\n")
		fmt.Printf("  Server BaseURL: %s\n", pConfig.BundleSender.ResearchDomainAgent.Server.BaseUrl)
		fmt.Printf("  Project: %s\n", pConfig.BundleSender.ResearchDomainAgent.Project)
		DivdlnS()

	},
}

func init() {
	projectConfigCmd.Flags().StringVarP(&prjName, "projectName", "n", "", "Please Enter the Name of the Project")

	// Set Flag as Required
	//if err := projectConfigCmd.MarkFlagRequired("projectName"); err != nil {
	//	log.Fatalf("Error marking projectName as required: %v", err)
	//}

	rootCmd.AddCommand(projectConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
