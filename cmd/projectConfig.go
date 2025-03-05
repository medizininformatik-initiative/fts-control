package cmd

import (
	"fmt"

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

// projectConfigCmd represents the projectConfig command
var projectConfigCmd = &cobra.Command{
	Use:   "projectConfig",
	Short: "Project configuration",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("projectConfig called")
	},
}

func init() {
	rootCmd.AddCommand(projectConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
