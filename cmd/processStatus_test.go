package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

// TestProcessStatusCommandExists checks if the processStatus command is properly registered in rootCmd.
func TestProcessStatusCommandExists(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "processStatus" {
			found = c
			break
		}
	}
	if found == nil {
		t.Errorf("The processStatus command should exist in rootCmd")
	}
}

// TestHasProcessIDFlag verifies that the processStatus command includes the required 'processId' flag.
func TestHasProcessIDFlag(t *testing.T) {
	flag := processStatusCmd.Flags().Lookup("processId")
	if flag == nil {
		t.Errorf("The processStatus command should have a 'processId' flag")
		return
	}
	if flag.Changed {
		t.Errorf("The 'processId' flag should not be set by default")
	}
}
