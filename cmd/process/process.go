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
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify a subcommand (use '%s --help' to see available options)", cmd.CommandPath())
		}
		return fmt.Errorf("unknown subcommand '%s' (use '%s --help' to see available subcommands)", args[0], cmd.CommandPath())
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
