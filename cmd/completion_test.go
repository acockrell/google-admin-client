package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionCommandRegistration(t *testing.T) {
	// Check that completion command is registered
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			found = true

			// Check for subcommands
			subcommands := map[string]bool{
				"bash": false,
				"zsh":  false,
				"fish": false,
			}

			for _, subcmd := range cmd.Commands() {
				if _, exists := subcommands[subcmd.Name()]; exists {
					subcommands[subcmd.Name()] = true
				}
			}

			// Verify all expected subcommands are present
			for shell, foundSubcmd := range subcommands {
				if !foundSubcmd {
					t.Errorf("completion command missing %s subcommand", shell)
				}
			}

			break
		}
	}

	if !found {
		t.Error("completion command not registered with root command")
	}
}

func TestCompletionSubcommands(t *testing.T) {
	var completionCmd *cobra.Command

	// Find the completion command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			completionCmd = cmd
			break
		}
	}

	if completionCmd == nil {
		t.Fatal("completion command not found")
	}

	tests := []struct {
		name        string
		subcommand  string
		wantPresent bool
	}{
		{"bash completion", "bash", true},
		{"zsh completion", "zsh", true},
		{"fish completion", "fish", true},
		{"pwsh should not exist", "pwsh", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, subcmd := range completionCmd.Commands() {
				if subcmd.Name() == tt.subcommand {
					found = true
					break
				}
			}

			if found != tt.wantPresent {
				if tt.wantPresent {
					t.Errorf("Expected %s subcommand to be present, but it was not found", tt.subcommand)
				} else {
					t.Errorf("Expected %s subcommand to not be present, but it was found", tt.subcommand)
				}
			}
		})
	}
}

func TestCompletionValidArgs(t *testing.T) {
	var completionCmd *cobra.Command

	// Find the completion command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			completionCmd = cmd
			break
		}
	}

	if completionCmd == nil {
		t.Fatal("completion command not found")
		return
	}

	expectedValidArgs := []string{"bash", "zsh", "fish"}
	if len(completionCmd.ValidArgs) != len(expectedValidArgs) {
		t.Errorf("ValidArgs length = %d, want %d", len(completionCmd.ValidArgs), len(expectedValidArgs))
	}

	for _, expected := range expectedValidArgs {
		found := false
		for _, validArg := range completionCmd.ValidArgs {
			if validArg == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected valid arg %s not found in ValidArgs", expected)
		}
	}
}
