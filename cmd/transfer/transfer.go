/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package transfer

import (
	"fmt"

	"github.com/spf13/cobra"
)

// transferCmd represents the transfer command

var Cmd = &cobra.Command{
	Use:   "transfer",
	Short: "Start transfer processes",
	Long: `Provides commands to start transfer processes, 
such as initiating a transfer of consented patients in a project.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a subcommand. Use 'ftsctl transfer --help' to see available options.")
			return
		}
		// unknown subcommand
		fmt.Printf("Error: unknown subcommand '%s'\n", args[0])
		fmt.Println("Use 'ftsctl transfer --help' to see available subcommands.")
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
	// transferCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transferCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
