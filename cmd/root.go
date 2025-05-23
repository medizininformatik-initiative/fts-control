/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"ftsctl/cmd/process"
	"ftsctl/cmd/project"
	"ftsctl/cmd/transfer"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ftsctl",
	Short: "Control the FTSnext-system from the command line.",
	Long: `███████╗████████╗███████╗ ██████╗████████╗██╗     
██╔════╝╚══██╔══╝██╔════╝██╔════╝╚══██╔══╝██║     
█████╗     ██║   ███████╗██║        ██║   ██║     
██╔══╝     ██║   ╚════██║██║        ██║   ██║     
██║        ██║   ███████║╚██████╗   ██║   ███████╗
╚═╝        ╚═╝   ╚══════╝ ╚═════╝   ╚═╝   ╚══════╝
                                                  
Manage and control the FTSnext-system efficiently 
and seamlessly directly from the command line interface, 
enabling advanced configuration and operation 
through text-based commands.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(project.Cmd)
	rootCmd.AddCommand(process.Cmd)
	rootCmd.AddCommand(transfer.Cmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.v1.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
