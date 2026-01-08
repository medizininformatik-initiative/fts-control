package config

import (
	"fmt"
	"log/slog"
	"os"

	"ftsctl/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display the current configuration settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := utils.GetConfigPath()
		if err != nil {
			slog.Error("Failed to get config path", "error", err)
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if !utils.ConfigExists() {
			utils.DivdlnL()
			fmt.Println("No configuration file found.")
			fmt.Printf("Expected location: %s\n", configPath)
			fmt.Println("\nRun 'ftsctl config set-base-url <url>' to create one.")
			utils.DivdlnL()
			return
		}

		baseUrl := viper.GetString("api.base_url")

		utils.DivdlnL()
		fmt.Println("Current Configuration:")
		utils.DivdlnL()
		fmt.Printf("Base URL: %s\n", baseUrl)
		fmt.Printf("Config file: %s\n", configPath)
		utils.DivdlnS()

		slog.Debug("Displayed config", "base_url", baseUrl, "path", configPath)
	},
}

func init() {
	Cmd.AddCommand(ShowCmd)
}
