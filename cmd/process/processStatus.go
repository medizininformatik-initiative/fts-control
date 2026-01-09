package process

import (
	"fmt"
	"io"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Show process status",
	Long:         `Shows the process status of the selected process ID`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		processId, _ := cmd.Flags().GetString("processId")
		return ExecuteProcessStatus(utils.NewClient(), cmd.OutOrStdout(), processId)
	},
}

// ExecuteProcessStatus fetches and displays the status of a specific process.
// This function is exported for testing.
func ExecuteProcessStatus(client *utils.Client, w io.Writer, procId string) error {
	var prcss utils.Process

	if err := client.GetJSON("/api/v2/process/status/"+procId, &prcss); err != nil {
		return fmt.Errorf("failed to fetch status for process %q: %w", procId, err)
	}

	utils.FprintDivdlnL(w)
	_, _ = fmt.Fprintf(w, "Representation of the Process with the ProcessID: %s \n", procId)
	utils.FprintDivdlnL(w)
	utils.FprintProcessStatus(w, prcss)
	utils.FprintDivdlnS(w)
	return nil
}

func init() {
	StatusCmd.Flags().StringP("processId", "i", "", "Process ID to query")
	_ = StatusCmd.MarkFlagRequired("processId")
	Cmd.AddCommand(StatusCmd)
}
