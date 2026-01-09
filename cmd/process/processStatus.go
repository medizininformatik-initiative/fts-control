package process

import (
	"fmt"
	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var processId string

var StatusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Show process status",
	Long:         `Shows the process status of the selected process ID`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := utils.NewClient()
		var prcss utils.Process

		if err := client.GetJSON("/api/v2/process/status/"+processId, &prcss); err != nil {
			return fmt.Errorf("failed to fetch status for process %q: %w", processId, err)
		}

		utils.DivdlnL()
		fmt.Printf("Representation of the Process with the ProcessID: %s \n", processId)
		utils.DivdlnL()
		utils.PrintProcessStatus(prcss)
		utils.DivdlnS()
		return nil
	},
}

func init() {
	StatusCmd.Flags().StringVarP(&processId, "processId", "i", "", "Process ID to query")
	_ = StatusCmd.MarkFlagRequired("processId")
	Cmd.AddCommand(StatusCmd)
}
