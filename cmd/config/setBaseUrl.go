package config

import (
	"fmt"
	"log/slog"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
)

var SetBaseUrlCmd = &cobra.Command{
	Use:   "set-base-url <url>",
	Short: "Set the base API URL",
	Long: `Set the base API URL in the config file. Creates the file if it doesn't exist.

Example:
  ftsctl config set-base-url http://localhost:8080
  ftsctl config set-base-url https://api.example.com`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		urlArg := args[0]

		validatedUrl, err := utils.ValidateURL(urlArg)
		if err != nil {
			return fmt.Errorf("URL validation failed: %w", err)
		}

		slog.Debug("Validated URL", "url", validatedUrl)

		if err := utils.WriteConfig(validatedUrl); err != nil {
			return fmt.Errorf("failed to write configuration: %w", err)
		}

		configPath, err := utils.GetConfigPath()
		if err != nil {
			slog.Warn("Could not determine config path for display", "error", err)
			configPath = "<unknown>"
		}

		utils.DivdlnL()
		fmt.Println("Configuration updated successfully!")
		utils.DivdlnS()
		fmt.Printf("Base URL: %s\n", validatedUrl)
		fmt.Printf("Config file: %s\n", configPath)
		utils.DivdlnS()
		return nil
	},
}

func init() {
	Cmd.AddCommand(SetBaseUrlCmd)
}
