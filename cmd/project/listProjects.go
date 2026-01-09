package project

import (
	"fmt"

	"ftsctl/cmd/utils"
	"github.com/spf13/cobra"
)

var ListProjectsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List of all available projects",
	Long:         `List of all available projects from the API.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := utils.NewClient()
		var projects []string

		if err := client.GetJSON("/api/v2/projects", &projects); err != nil {
			return fmt.Errorf("failed to fetch projects: %w", err)
		}

		utils.DivdlnL()
		fmt.Println("List of all available projects:")
		utils.DivdlnL()
		for i, project := range projects {
			fmt.Printf("%d. %s\n", i+1, project)
		}
		utils.DivdlnS()
		return nil
	},
}

func init() {

	Cmd.AddCommand(ListProjectsCmd)
}
