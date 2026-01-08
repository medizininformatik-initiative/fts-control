package transfer

import (
	"fmt"
	"ftsctl/cmd/utils"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

var startTransferCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a transfer process",
	Long: `Start a transfer of patients with IDs provided in the request
body, or if none are provided, start a transfer of all consented
patients for the specified project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("projectName")
		if projectName == "" {
			fmt.Println("ERROR: The --projectName (-n) flag is required!")
			fmt.Println("Usage: ftsctl startTransfer --projectName exampleProject [--ids id1,id2,id3]")
			return
		}

		idsStr, _ := cmd.Flags().GetString("ids")
		var ids []string
		if idsStr != "" {
			ids = strings.Split(idsStr, ",")
		}

		client := utils.NewClient()
		endpoint := fmt.Sprintf("/api/v2/process/%s/start", projectName)

		var body interface{}
		if len(ids) > 0 {
			body = ids
		}

		if err := client.PostJSON(endpoint, body, nil); utils.PrintHTTPError(err) {
			return
		}

		utils.DivdlnL()
		if len(ids) > 0 {
			fmt.Printf("Transfer of project '%s' has started with the following patient IDs:\n", projectName)
			for _, id := range ids {
				fmt.Printf("   - %s\n", id)
			}
		} else {
			fmt.Printf("Transfer of project '%s' has started for all consented patients.\n", projectName)
		}
		slog.Debug("Transfer request completed", "project", projectName)
		utils.DivdlnS()
	},
}

func init() {
	startTransferCmd.Flags().StringP("projectName", "n", "", "Project name (required)")
	startTransferCmd.Flags().StringP("ids", "i", "", "Comma-separated list of patient IDs (optional)")

	Cmd.AddCommand(startTransferCmd)
}
