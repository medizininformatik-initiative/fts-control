package project

import (
	"fmt"

	"ftsctl/cmd/utils"
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

var PrjName string

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Project configuration",
	Long:  `Shows the project configuration of the selected project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("projectName") {
			utils.DivdlnL()
			fmt.Printf("Warning: No --projectName (-n) specified.\n\nPlease set the Flag to retrieve the Project Configuration of a Project.\n")
			utils.DivdlnL()
			return
		}

		client := utils.NewClient()
		var pConfig Config

		if err := client.GetJSON("/api/v2/projects/"+PrjName, &pConfig); utils.LogHTTPError(err, "fetch project config") {
			return
		}

		// Formatted Output
		utils.DivdlnL()
		fmt.Printf("Project configuration with the Project Name: %s \n", PrjName)
		utils.DivdlnL()

		fmt.Printf("CohortSelector - TrustCenterAgent:\n")
		fmt.Printf("  Server BaseURL: %s\n", pConfig.CohortSelector.TrustCenterAgent.Server.BaseUrl)
		fmt.Printf("  Domain: %s\n", pConfig.CohortSelector.TrustCenterAgent.Domain)
		fmt.Printf("  Patient Identifier System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PatientIdentifierSystem)
		fmt.Printf("  Policy System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PolicySystem)
		fmt.Println("  Policies:")
		for _, policy := range pConfig.CohortSelector.TrustCenterAgent.Policies {
			fmt.Printf("    - %s\n", policy)
		}

		utils.DivdlnS()
		fmt.Printf("DataSelector - Everything:\n")
		fmt.Printf("  FHIR Server BaseURL: %s\n", pConfig.DataSelector.Everything.FhirServer.BaseUrl)
		fmt.Printf("  Resolve Patient Identifier System: %s\n", pConfig.DataSelector.Everything.Resolve.PatientIdentifierSystem)

		utils.DivdlnS()
		fmt.Printf("Deidentificator - Deidentifhir:\n")
		fmt.Printf("  TrustCenterAgent Server BaseURL: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Server.BaseUrl)
		fmt.Printf("  Max Date Shift: %s\n", pConfig.Deidentificator.Deidentifhir.MaxDateShift)
		fmt.Printf("  Deidentifhir Config: %s\n", pConfig.Deidentificator.Deidentifhir.DeidentifhirConfig)
		fmt.Printf("  Scraper Config: %s\n", pConfig.Deidentificator.Deidentifhir.ScraperConfig)
		fmt.Println("  Domains:")
		fmt.Printf("    Pseudonym: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Pseudonym)
		fmt.Printf("    Salt: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Salt)
		fmt.Printf("    DateShift: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.DateShift)

		utils.DivdlnS()
		fmt.Printf("BundleSender - ResearchDomainAgent:\n")
		fmt.Printf("  Server BaseURL: %s\n", pConfig.BundleSender.ResearchDomainAgent.Server.BaseUrl)
		fmt.Printf("  Project: %s\n", pConfig.BundleSender.ResearchDomainAgent.Project)
		utils.DivdlnS()
	},
}

func init() {
	ConfigCmd.Flags().StringVarP(&PrjName, "projectName", "n", "", "Please Enter the Name of the Project")
	Cmd.AddCommand(ConfigCmd)
}
