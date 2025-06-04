/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package process

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "process",
	Short: "Inspect and monitor transfer processes",
	Long: `Provides commands to inspect and monitor transfer processes, 
including listing all available processes and showing the status 
of a specific process.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a subcommand. Use 'ftsctl process --help' to see available options.")
			return
		}
		// unknown subcommand
		fmt.Printf("Error: unknown subcommand '%s'\n", args[0])
		fmt.Println("Use 'ftsctl process --help' to see available subcommands.")
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

}
