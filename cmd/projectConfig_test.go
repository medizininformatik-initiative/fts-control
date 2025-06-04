package cmd

import (
	"ftsctl/cmd/project"
	"testing"
)

func TestProjectConfingCommandExists(t *testing.T) {
	cmd, _, err := project.Cmd.Find([]string{"ProjectConfing"})
	if err != nil || cmd == nil {
		t.Errorf("The 'ProjectConfing' subcommand should exist under 'process'")
	}
}

func TestProjectConfigCmd_MissingProjectName(t *testing.T) {
	// Temporarily clear the project name
	project.PrjName = ""

	// Capture log output
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected an error due to missing project name, but command ran successfully")
		}
	}()

	// Run the command without a project name
	project.ConfigCmd.Run(nil, nil)
}
