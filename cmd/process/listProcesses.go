package process

import (
	"fmt"
	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var listProcessesCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all transfer process statuses",
	Long:         `Lists all available transfer process statuses.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := utils.NewClient()
		var processes []utils.Process

		if err := client.GetJSON("/api/v2/process/statuses", &processes); err != nil {
			return fmt.Errorf("failed to fetch process statuses: %w", err)
		}

		utils.DivdlnL()
		fmt.Println("List of all transfer process statuses:")
		utils.DivdlnL()

		for _, p := range processes {
			utils.PrintProcessStatus(p)
			utils.DivdlnS()
		}
		return nil
	},
}

func init() {
	Cmd.AddCommand(listProcessesCmd)
}
