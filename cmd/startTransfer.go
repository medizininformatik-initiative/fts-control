/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startTransferCmd represents the startTransfer command
var startTransferCmd = &cobra.Command{
	Use:   "startTransfer", // /api/v2/process/{project}/start
	Short: "Start a transfer process",
	Long:  `Start a transfer of patients with IDs given in the request body or if empty start a transfer of all consented patients.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("startTransfer called")
	},
}

func init() {
	rootCmd.AddCommand(startTransferCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startTransferCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startTransferCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
