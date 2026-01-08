package transfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ftsctl/cmd/utils"
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// RequestPayload defines the structure for the optional ID list
type RequestPayload struct {
	IDs []string `json:"ids,omitempty"`
}

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

		apiUrl := utils.BuildApiUrl(fmt.Sprintf("/api/v2/process/%s/start", projectName))

		var resp *http.Response
		var err error

		if len(ids) > 0 {
			jsonData, err := json.Marshal(ids)
			if err != nil {
				slog.Error("Error creating JSON payload", "error", err)
				os.Exit(1)
			}
			slog.Debug("JSON payload created", "payload", string(jsonData))

			resp, err = http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonData))
		} else {
			resp, err = http.Post(apiUrl, "application/json", nil)
		}

		if err != nil {
			slog.Error("Error sending request", "error", err)
			os.Exit(1)
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				slog.Warn("Error closing response body", "error", cerr)
			}
		}()

		utils.DivdlnL()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
			if len(ids) > 0 {
				fmt.Printf("Transfer of project '%s' has started with the following patient IDs:\n", projectName)
				for _, id := range ids {
					fmt.Printf("   - %s\n", id)
				}
			} else {
				fmt.Printf("Transfer of project '%s' has started for all consented patients.\n", projectName)
			}
			slog.Debug("Transfer request completed", "url", apiUrl, "status", resp.Status)
		} else {
			fmt.Printf("Request failed! Response Status: %s\n", resp.Status)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Error reading response body", "error", err)
			} else {
				fmt.Printf("Server response body:\n%s\n", string(bodyBytes))
			}
			utils.DivdlnL()
		}
	},
}

func init() {
	startTransferCmd.Flags().StringP("projectName", "n", "", "Project name (required)")
	startTransferCmd.Flags().StringP("ids", "i", "", "Comma-separated list of patient IDs (optional)")

	Cmd.AddCommand(startTransferCmd)
}
