package config

import (
	"fmt"
	"log/slog"
	"os"

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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		urlArg := args[0]

		validatedUrl, err := utils.ValidateURL(urlArg)
		if err != nil {
			slog.Error("URL validation failed", "url", urlArg, "error", err)
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		slog.Debug("Validated URL", "url", validatedUrl)

		if err := utils.WriteConfig(validatedUrl); err != nil {
			slog.Error("Failed to write config", "error", err)
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
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
	},
}

func init() {
	Cmd.AddCommand(SetBaseUrlCmd)
}
