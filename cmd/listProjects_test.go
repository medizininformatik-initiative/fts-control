package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"testing"
)

// TestListProjectsCommandExists checks if the processStatus command is properly registered in rootCmd.
func TestListProjectsCommandExists(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "listProjects" {
			found = c
			break
		}
	}
	if found == nil {
		t.Errorf("The processStatus command should exist in rootCmd")
	}
}

func TestJSONParsing(t *testing.T) {
	// simulated API response as JSON string
	jsonData := `["Project1", "Project2", "Project3"]`
	var projects []string

	err := json.Unmarshal([]byte(jsonData), &projects)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// expected number of projects
	expectedLength := 3
	if len(projects) != expectedLength {
		t.Errorf("Expected %d projects, but got %d", expectedLength, len(projects))
	}
	// check if frist project is correct

	expectedFirstProject := "Project1"
	if projects[0] != expectedFirstProject {
		t.Errorf("Expected first project to be %s, but got %s", expectedFirstProject, projects[0])
	}

	// check if second project is correct
	expectedSecondProject := "Project2"
	if projects[1] != expectedSecondProject {
		t.Errorf("Expected second project to be %s, but got %s", expectedSecondProject, projects[1])
	}

	// check if third project is correct
	expectedThirdProject := "Project3"
	if projects[2] != expectedThirdProject {
		t.Errorf("Expected third project to be %s, but got %s", expectedThirdProject, projects[2])
	}
}

func TestApiUrlGeneration(t *testing.T) {
	baseUrl := "http://localhost:8080"
	expectedUrl := "http://localhost:8080/api/v2/projects"

	apiUrl := fmt.Sprintf("%s/api/v2/projects", baseUrl)

	if apiUrl != expectedUrl {
		t.Errorf("Expected URL: %s, but got: %s", expectedUrl, apiUrl)
	}
}
