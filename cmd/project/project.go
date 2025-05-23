/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package project

import (
	"fmt"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command

var Cmd = &cobra.Command{
	Use:   "project",
	Short: "Inspect and monitor projects",
	Long: `Provides commands to	inspect projects, 
including listing all available projects and showing the 
configuration of a selected project.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a subcommand. Use 'ftsctl project --help' to see available options.")
			return
		}
		// unknown subcommand
		fmt.Printf("Error: unknown subcommand '%s'\n", args[0])
		fmt.Println("Use 'ftsctl project --help' to see available subcommands.")
	},
}

func init() {
	Cmd.SetHelpTemplate(`
{{.Long}}

Usage:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Available Subcommands:
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}{{end}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
