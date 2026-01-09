package project

import (
	"fmt"
	"io"

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

// PrjName is exported for backward compatibility with existing tests
var PrjName string

var ConfigCmd = &cobra.Command{
	Use:          "config",
	Short:        "Project configuration",
	Long:         `Shows the project configuration of the selected project.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, _ := cmd.Flags().GetString("projectName")
		return ExecuteProjectConfig(utils.NewClient(), cmd.OutOrStdout(), projectName)
	},
}

// ExecuteProjectConfig fetches and displays the project configuration.
// This function is exported for testing.
func ExecuteProjectConfig(client *utils.Client, w io.Writer, projectName string) error {
	var pConfig Config

	if err := client.GetJSON("/api/v2/projects/"+projectName, &pConfig); err != nil {
		return fmt.Errorf("failed to fetch configuration for project %q: %w", projectName, err)
	}

	// Formatted Output
	utils.FprintDivdlnL(w)
	fmt.Fprintf(w, "Project configuration with the Project Name: %s \n", projectName)
	utils.FprintDivdlnL(w)

	fmt.Fprintln(w, "CohortSelector - TrustCenterAgent:")
	fmt.Fprintf(w, "  Server BaseURL: %s\n", pConfig.CohortSelector.TrustCenterAgent.Server.BaseUrl)
	fmt.Fprintf(w, "  Domain: %s\n", pConfig.CohortSelector.TrustCenterAgent.Domain)
	fmt.Fprintf(w, "  Patient Identifier System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PatientIdentifierSystem)
	fmt.Fprintf(w, "  Policy System: %s\n", pConfig.CohortSelector.TrustCenterAgent.PolicySystem)
	fmt.Fprintln(w, "  Policies:")
	for _, policy := range pConfig.CohortSelector.TrustCenterAgent.Policies {
		fmt.Fprintf(w, "    - %s\n", policy)
	}

	utils.FprintDivdlnS(w)
	fmt.Fprintln(w, "DataSelector - Everything:")
	fmt.Fprintf(w, "  FHIR Server BaseURL: %s\n", pConfig.DataSelector.Everything.FhirServer.BaseUrl)
	fmt.Fprintf(w, "  Resolve Patient Identifier System: %s\n", pConfig.DataSelector.Everything.Resolve.PatientIdentifierSystem)

	utils.FprintDivdlnS(w)
	fmt.Fprintln(w, "Deidentificator - Deidentifhir:")
	fmt.Fprintf(w, "  TrustCenterAgent Server BaseURL: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Server.BaseUrl)
	fmt.Fprintf(w, "  Max Date Shift: %s\n", pConfig.Deidentificator.Deidentifhir.MaxDateShift)
	fmt.Fprintf(w, "  Deidentifhir Config: %s\n", pConfig.Deidentificator.Deidentifhir.DeidentifhirConfig)
	fmt.Fprintf(w, "  Scraper Config: %s\n", pConfig.Deidentificator.Deidentifhir.ScraperConfig)
	fmt.Fprintln(w, "  Domains:")
	fmt.Fprintf(w, "    Pseudonym: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Pseudonym)
	fmt.Fprintf(w, "    Salt: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.Salt)
	fmt.Fprintf(w, "    DateShift: %s\n", pConfig.Deidentificator.Deidentifhir.TrustCenterAgent.Domains.DateShift)

	utils.FprintDivdlnS(w)
	fmt.Fprintln(w, "BundleSender - ResearchDomainAgent:")
	fmt.Fprintf(w, "  Server BaseURL: %s\n", pConfig.BundleSender.ResearchDomainAgent.Server.BaseUrl)
	fmt.Fprintf(w, "  Project: %s\n", pConfig.BundleSender.ResearchDomainAgent.Project)
	utils.FprintDivdlnS(w)
	return nil
}

func init() {
	ConfigCmd.Flags().StringVarP(&PrjName, "projectName", "n", "", "Name of the project")
	_ = ConfigCmd.MarkFlagRequired("projectName")
	Cmd.AddCommand(ConfigCmd)
}
