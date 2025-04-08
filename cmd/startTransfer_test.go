package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"testing"
)

func TestStartTransferCommandExists(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "startTransfer" {
			found = c
			break
		}
	}
	if found == nil {
		t.Errorf("The processStatus command should exist in rootCmd")
	}
}

func TestRequestPayloadMarshaling_OmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		payload  requestPayload
		expected string
	}{
		{
			name:     "Nil IDs",
			payload:  requestPayload{IDs: nil},
			expected: `{}`,
		},
		{
			name:     "Empty slice IDs",
			payload:  requestPayload{IDs: []string{}},
			expected: `{}`,
		},
		{
			name:     "With IDs",
			payload:  requestPayload{IDs: []string{"id1", "id2"}},
			expected: `{"ids":["id1","id2"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("unexpected error while marshaling: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("expected JSON: %s, but got: %s", tt.expected, string(result))
			}
		})
	}
}
