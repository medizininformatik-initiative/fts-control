package transfer

import (
	"fmt"
	"io"
	"log/slog"
	"strings"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var startTransferCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a transfer process",
	Long: `Start a transfer of patients with IDs provided in the request
body, or if none are provided, start a transfer of all consented
patients for the specified project.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName, _ := cmd.Flags().GetString("projectName")
		idsStr, _ := cmd.Flags().GetString("ids")

		var ids []string
		if idsStr != "" {
			ids = strings.Split(idsStr, ",")
		}

		return ExecuteStartTransfer(utils.NewClient(), cmd.OutOrStdout(), projectName, ids)
	},
}

// ExecuteStartTransfer starts a transfer for the given project.
// If ids is empty/nil, transfers all consented patients.
// This function is exported for testing.
func ExecuteStartTransfer(client *utils.Client, w io.Writer, projectName string, ids []string) error {
	endpoint := utils.EndpointTransferStart(projectName)

	var body interface{}
	if len(ids) > 0 {
		body = ids
	}

	if err := client.PostJSON(endpoint, body, nil); err != nil {
		return fmt.Errorf("failed to start transfer for project %q: %w", projectName, err)
	}

	utils.FprintDivdlnL(w)
	if len(ids) > 0 {
		_, _ = fmt.Fprintf(w, "Transfer of project '%s' has started with the following patient IDs:\n", projectName)
		for _, id := range ids {
			_, _ = fmt.Fprintf(w, "   - %s\n", id)
		}
	} else {
		_, _ = fmt.Fprintf(w, "Transfer of project '%s' has started for all consented patients.\n", projectName)
	}
	slog.Debug("Transfer request completed", "project", projectName)
	utils.FprintDivdlnS(w)
	return nil
}

func init() {
	startTransferCmd.Flags().StringP("projectName", "n", "", "Project name (required)")
	startTransferCmd.Flags().StringP("ids", "i", "", "Comma-separated list of patient IDs (optional)")
	_ = startTransferCmd.MarkFlagRequired("projectName")
	Cmd.AddCommand(startTransferCmd)
}
