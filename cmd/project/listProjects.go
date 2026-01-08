package project

import (
	"fmt"

	"ftsctl/cmd/utils"
	"github.com/spf13/cobra"
)

var ListProjectsCmd = &cobra.Command{
	Use:   "list", // /api/v2/projects
	Short: "List of all available projects",
	Long:  `List of all available projects from the API.`,

	Run: func(cmd *cobra.Command, args []string) {
		client := utils.NewClient()
		var projects []string

		if err := client.GetJSON("/api/v2/projects", &projects); utils.LogHTTPError(err, "fetch projects") {
			return
		}

		utils.DivdlnL()
		fmt.Println("List of all available projects:")
		utils.DivdlnL()
		for i, project := range projects {
			fmt.Printf("%d. %s\n", i+1, project)
		}
		utils.DivdlnS()
	},
}

func init() {

	Cmd.AddCommand(ListProjectsCmd)
}
