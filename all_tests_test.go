package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestRunCmdTests(t *testing.T) {
	cmd := exec.Command("go", "test", "./cmd/...", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Dir = "."

	err := cmd.Run()
	if err != nil {
		t.Fatalf("cmd tests failed: %v", err)
	}
}
