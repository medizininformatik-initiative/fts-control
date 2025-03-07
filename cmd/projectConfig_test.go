package cmd

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestProjectConfigCommandExists(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "projectConfig" {
			found = c
			break
		}
	}
	if found == nil {
		t.Errorf("The processStatus command should exist in rootCmd")
	}
}

func TestProjectConfigCmd_MissingProjectName(t *testing.T) {
	// Temporarily clear the project name
	prjName = ""

	// Capture log output
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected an error due to missing project name, but command ran successfully")
		}
	}()

	// Run the command without a project name
	projectConfigCmd.Run(nil, nil)
}
