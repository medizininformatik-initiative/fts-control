package cmd

import (
	"ftsctl/cmd/process"
	"testing"
)

// TestProcessStatusCommandExists checks if the processStatus command is properly registered in rootCmd.
func TestProcessStatusCommandExists(t *testing.T) {
	cmd, _, err := process.Cmd.Find([]string{"processStatus"})
	if err != nil || cmd == nil {
		t.Errorf("The 'processStatus' subcommand should exist under 'process'")
	}
}

// TestHasProcessIDFlag verifies that the processStatus command includes the required 'processId' flag.
func TestHasProcessIDFlag(t *testing.T) {
	flag := process.StatusCmd.Flags().Lookup("processId")
	if flag == nil {
		t.Errorf("The processStatus command should have a 'processId' flag")
		return
	}
	if flag.Changed {
		t.Errorf("The 'processId' flag should not be set by default")
	}
}
