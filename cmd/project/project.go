/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package project

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "project",
	Short: "Inspect and monitor projects",
	Long: `Provides commands to	inspect projects,
including listing all available projects and showing the
configuration of a selected project.`,
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
