package cmd

import (
	"strings"
	"testing"
)

func TestGroupEmailConstruction(t *testing.T) {
	// Save original domain
	originalDomain := domain
	defer func() { domain = originalDomain }()

	tests := []struct {
		name           string
		input          string
		configDomain   string
		expectedOutput string
		description    string
	}{
		{
			name:           "group without @ gets domain appended",
			input:          "engineering",
			configDomain:   "example.com",
			expectedOutput: "engineering@example.com",
			description:    "Short group name should get domain appended",
		},
		{
			name:           "group with @ unchanged",
			input:          "engineering@custom.com",
			configDomain:   "example.com",
			expectedOutput: "engineering@custom.com",
			description:    "Full email should remain unchanged",
		},
		{
			name:           "empty domain with short name",
			input:          "engineering",
			configDomain:   "",
			expectedOutput: "engineering@",
			description:    "With empty domain, still appends @",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain = tt.configDomain

			groupEmail := tt.input
			if !strings.Contains(tt.input, "@") {
				groupEmail = tt.input + "@" + getDomain()
			}

			if groupEmail != tt.expectedOutput {
				t.Errorf("Group email construction = %q, want %q (%s)",
					groupEmail, tt.expectedOutput, tt.description)
			}
		})
	}
}

func TestGroupCommandExists(t *testing.T) {
	if groupCmd == nil {
		t.Fatal("groupCmd should not be nil")
	}

	if groupCmd.Use != "group" {
		t.Errorf("groupCmd.Use = %q, want %q", groupCmd.Use, "group")
	}
}

func TestGroupListCommandExists(t *testing.T) {
	found := false
	for _, cmd := range groupCmd.Commands() {
		if cmd.Use == "list" {
			found = true
			break
		}
	}

	if !found {
		t.Error("group list command should exist")
	}
}
