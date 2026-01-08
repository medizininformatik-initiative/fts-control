package cmd

import (
	"encoding/json"
	"ftsctl/cmd/transfer"
	"testing"
)

func TestStartTransferCommandExists(t *testing.T) {
	cmd, _, err := transfer.Cmd.Find([]string{"StartTransfer"})
	if err != nil || cmd == nil {
		t.Errorf("The 'StartTransfer' subcommand should exist under 'process'")
	}
}

func TestIDsPayloadMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		expected string
	}{
		{
			name:     "Nil IDs",
			ids:      nil,
			expected: `null`,
		},
		{
			name:     "Empty slice IDs",
			ids:      []string{},
			expected: `[]`,
		},
		{
			name:     "With IDs",
			ids:      []string{"id1", "id2"},
			expected: `["id1","id2"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.ids)
			if err != nil {
				t.Fatalf("unexpected error while marshaling: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("expected JSON: %s, but got: %s", tt.expected, string(result))
			}
		})
	}
}
