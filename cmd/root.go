package cmd

import (
	"ftsctl/cmd/config"
	"ftsctl/cmd/process"
	"ftsctl/cmd/project"
	"ftsctl/cmd/transfer"
	"ftsctl/cmd/utils"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose     bool
	baseURLFlag string
)

var rootCmd = &cobra.Command{
	Use:   "ftsctl",
	Short: "Control the FTSnext-system from the command line.",
	Long: utils.GenerateLogo() + utils.GenerateLogoV2() + `Manage and control the FTSnext-system efficiently 
and seamlessly directly from the command line interface, 
enabling advanced configuration and operation 
through text-based commands.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		level := slog.LevelError
		if verbose {
			level = slog.LevelDebug
		}

		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		}))
		slog.SetDefault(logger)

		slog.Debug("Verbose mode is active")

		// Base URL override
		if baseURLFlag != "" {
			slog.Debug("Using base URL from flag", "baseURL", baseURLFlag)
			viper.Set("api.base_url", baseURLFlag)
		} else {
			baseFromConfig := viper.GetString("api.base_url")
			slog.Debug("Using base URL from config file", "baseURL", baseFromConfig)
		}
	},
}

func Execute() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Use OS-specific config directory
	configDir, err := utils.GetConfigDir()
	if err != nil {
		slog.Warn("Could not determine config directory", "error", err)
	} else {
		viper.AddConfigPath(configDir)
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Debug("Could not read config file, continuing without it", "error", err)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(config.Cmd)
	rootCmd.AddCommand(project.Cmd)
	rootCmd.AddCommand(process.Cmd)
	rootCmd.AddCommand(transfer.Cmd)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&baseURLFlag, "base-url", "u", "", "Override base API URL")
}
