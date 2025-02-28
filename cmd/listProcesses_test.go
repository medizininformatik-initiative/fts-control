package cmd

import (
	"encoding/json"
	"testing"
	"time"
)

// JSON unmarshalling of a Process object
func TestProcess_JSONUnmarshal(t *testing.T) {
	jsonData := `{"processId":"123","phase":"running","createdAt":[2024,2,28,15,30,45,0],"totalPatients":100,"totalBundles":50,"deidentifiedBundles":25,"sentBundles":10,"skippedBundles":5}`
	var p Process

	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	if p.ProcessID != "123" || p.Phase != "running" {
		t.Errorf("Incorrect values after unmarshalling: %+v", p)
	}
}

// Verify the ToTime method

func TestTimeArray_ToTime(t *testing.T) {
	timeArray := TimeArray{2024, 2, 28, 15, 30, 45, 0}
	expectedTime := time.Date(2024, 2, 28, 15, 30, 45, 0, time.UTC)

	result := timeArray.ToTime()
	if !result.Equal(expectedTime) {
		t.Errorf("Expected: %v, but got: %v", expectedTime, result)
	}
}

// Check if FinishedAt can be nil
func TestProcess_FinishedAtNil(t *testing.T) {
	jsonData := `{"processId":"123","phase":"running","createdAt":[2024,2,28,15,30,45,0],"totalPatients":100,"totalBundles":50,"deidentifiedBundles":25,"sentBundles":10,"skippedBundles":5}`
	var p Process

	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	if p.FinishedAt != nil {
		t.Errorf("FinishedAt should be nil, but got: %v", p.FinishedAt)
	}
}
