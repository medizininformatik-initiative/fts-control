package project

import (
	"context"
	"fmt"
	"io"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var ListProjectsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List of all available projects",
	Long:         `List of all available projects from the API.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteListProjects(utils.NewClient(), cmd.OutOrStdout())
	},
}

// ExecuteListProjects fetches and displays the list of projects.
// This function is exported for testing.
func ExecuteListProjects(client *utils.Client, w io.Writer) error {
	var projects []string

	if err := client.GetJSON(context.Background(), utils.EndpointProjects, &projects); err != nil {
		return fmt.Errorf("failed to fetch projects: %w", err)
	}

	utils.FprintDivdlnL(w)
	_, _ = fmt.Fprintln(w, "List of all available projects:")
	utils.FprintDivdlnL(w)
	for i, project := range projects {
		_, _ = fmt.Fprintf(w, "%d. %s\n", i+1, project)
	}
	utils.FprintDivdlnS(w)
	return nil
}

func init() {
	Cmd.AddCommand(ListProjectsCmd)
}
