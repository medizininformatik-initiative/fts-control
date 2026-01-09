package process

import (
	"fmt"
	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var listProcessesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all transfer process statuses",
	Long:  `Lists all available transfer process statuses.`,
	Run: func(cmdCobra *cobra.Command, args []string) {
		client := utils.NewClient()
		var processes []utils.Process

		if err := client.GetJSON("/api/v2/process/statuses", &processes); utils.LogHTTPError(err, "fetch process statuses") {
			return
		}

		utils.DivdlnL()
		fmt.Println("List of all transfer process statuses:")
		utils.DivdlnL()

		for _, p := range processes {
			utils.PrintProcessStatus(p)
			utils.DivdlnS()
		}
	},
}

func init() {
	Cmd.AddCommand(listProcessesCmd)
}
