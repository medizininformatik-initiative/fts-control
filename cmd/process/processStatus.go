package process

import (
	"fmt"
	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var processId string

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show process status",
	Long:  `Shows the process status of the selected process ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("processId") {
			utils.DivdlnL()
			fmt.Printf("Warning: No --processId (-i) specified.\n\nPlease set the flag to retrieve the status of a process.\n")
			utils.DivdlnL()
			return
		}

		client := utils.NewClient()
		var prcss utils.Process

		if err := client.GetJSON("/api/v2/process/status/"+processId, &prcss); utils.LogHTTPError(err, "fetch process status") {
			return
		}

		utils.DivdlnL()
		fmt.Printf("Representation of the Process with the ProcessID: %s \n", processId)
		utils.DivdlnL()
		utils.PrintProcessStatus(prcss)
		utils.DivdlnS()
	},
}

func init() {
	StatusCmd.Flags().StringVarP(&processId, "processId", "i", "", "Please Enter the processId")
	Cmd.AddCommand(StatusCmd)
}
