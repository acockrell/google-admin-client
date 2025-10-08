package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigCommandRegistration(t *testing.T) {
	// Check that config command is registered
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "config" {
			found = true

			// Check for validate subcommand
			validateFound := false
			for _, subcmd := range cmd.Commands() {
				if subcmd.Name() == "validate" {
					validateFound = true
					break
				}
			}

			if !validateFound {
				t.Error("config command missing validate subcommand")
			}

			break
		}
	}

	if !found {
		t.Error("config command not registered with root command")
	}
}

func TestValidateToken(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "gac-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // Cleanup, error not critical in tests
	}()

	tests := []struct {
		name          string
		tokenContent  string
		expectedValid bool
		description   string
	}{
		{
			name:          "valid json token",
			tokenContent:  `{"access_token":"test","token_type":"Bearer","expiry":"2024-12-31T23:59:59Z"}`,
			expectedValid: true,
			description:   "should accept valid JSON token structure",
		},
		{
			name:          "empty file",
			tokenContent:  "",
			expectedValid: false,
			description:   "should reject empty token file",
		},
		{
			name:          "invalid json",
			tokenContent:  "not a json file",
			expectedValid: false,
			description:   "should reject non-JSON content",
		},
		{
			name:          "json array instead of object",
			tokenContent:  `["not", "an", "object"]`,
			expectedValid: false,
			description:   "should reject JSON arrays",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test token file
			tokenFile := filepath.Join(tmpDir, tt.name+".json")
			if err := os.WriteFile(tokenFile, []byte(tt.tokenContent), 0600); err != nil {
				t.Fatalf("Failed to write test token file: %v", err)
			}

			// Validate token
			result := validateToken(tokenFile)

			if result != tt.expectedValid {
				t.Errorf("validateToken() = %v, want %v (%s)", result, tt.expectedValid, tt.description)
			}
		})
	}
}

func TestValidateTokenNonExistentFile(t *testing.T) {
	// Test with non-existent file
	result := validateToken("/nonexistent/path/token.json")
	if result {
		t.Error("validateToken() should return false for non-existent file")
	}
}

func TestConfigValidateSubcommandPresence(t *testing.T) {
	var configCmd *cobra.Command

	// Find the config command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "config" {
			configCmd = cmd
			break
		}
	}

	if configCmd == nil {
		t.Fatal("config command not found")
	}

	// Check for validate subcommand
	found := false
	for _, subcmd := range configCmd.Commands() {
		if subcmd.Name() == "validate" {
			found = true

			// Verify it has a Run function
			if subcmd.Run == nil && subcmd.RunE == nil {
				t.Error("config validate command has no Run function")
			}

			break
		}
	}

	if !found {
		t.Error("config validate subcommand not found")
	}
}
