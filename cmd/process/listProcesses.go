package process

import (
	"context"
	"fmt"
	"io"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var listProcessesCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all transfer process statuses",
	Long:         `Lists all available transfer process statuses.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := utils.GetAuthenticatedClient()
		if err != nil {
			return fmt.Errorf("authentication error: %w", err)
		}
		return ExecuteListProcesses(client, cmd.OutOrStdout())
	},
}

// ExecuteListProcesses fetches and displays all process statuses.
// This function is exported for testing.
func ExecuteListProcesses(client *utils.Client, w io.Writer) error {
	var processes []utils.Process

	if err := client.GetJSON(context.Background(), utils.EndpointProcessStatuses, &processes); err != nil {
		return fmt.Errorf("failed to fetch process statuses: %w", err)
	}

	utils.FprintDivdlnL(w)
	_, _ = fmt.Fprintln(w, "List of all transfer process statuses:")
	utils.FprintDivdlnL(w)

	for _, p := range processes {
		utils.FprintProcessStatus(w, p)
		utils.FprintDivdlnS(w)
	}
	return nil
}

func init() {
	Cmd.AddCommand(listProcessesCmd)
}
