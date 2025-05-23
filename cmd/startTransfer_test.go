package cmd

import (
	"encoding/json"
	"ftsctl/cmd/process"
	"ftsctl/cmd/transfer"
	"testing"
)

func TestStartTransferCommandExists(t *testing.T) {
	cmd, _, err := process.Cmd.Find([]string{"StartTransfer"})
	if err != nil || cmd == nil {
		t.Errorf("The 'StartTransfer' subcommand should exist under 'process'")
	}
}

func TestRequestPayloadMarshaling_OmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		payload  transfer.RequestPayload
		expected string
	}{
		{
			name:     "Nil IDs",
			payload:  transfer.RequestPayload{IDs: nil},
			expected: `{}`,
		},
		{
			name:     "Empty slice IDs",
			payload:  transfer.RequestPayload{IDs: []string{}},
			expected: `{}`,
		},
		{
			name:     "With IDs",
			payload:  transfer.RequestPayload{IDs: []string{"id1", "id2"}},
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
