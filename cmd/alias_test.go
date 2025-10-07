package cmd

import (
	"testing"
)

// TestAliasCommandExists verifies that the alias command is registered
func TestAliasCommandExists(t *testing.T) {
	if aliasCmd == nil {
		t.Fatal("alias command is nil")
	}

	if aliasCmd.Use != "alias" {
		t.Errorf("Expected alias command Use to be 'alias', got '%s'", aliasCmd.Use)
	}

	if aliasCmd.Short == "" {
		t.Error("alias command should have a Short description")
	}

	if aliasCmd.Long == "" {
		t.Error("alias command should have a Long description")
	}
}

// TestAliasSubcommands verifies all expected alias subcommands exist
func TestAliasSubcommands(t *testing.T) {
	expectedCommands := []string{"list", "add", "remove"}

	commands := aliasCmd.Commands()
	commandMap := make(map[string]bool)

	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, cmdName := range expectedCommands {
		if !commandMap[cmdName] {
			t.Errorf("Alias command missing subcommand: %s", cmdName)
		}
	}
}

// TestAliasListCommand verifies the alias list command structure
func TestAliasListCommand(t *testing.T) {
	if aliasListCmd == nil {
		t.Fatal("alias list command is nil")
	}

	if aliasListCmd.Short == "" {
		t.Error("alias list command should have a Short description")
	}

	if aliasListCmd.Long == "" {
		t.Error("alias list command should have a Long description")
	}

	if aliasListCmd.RunE == nil {
		t.Error("alias list command should have a RunE function")
	}
}

// TestAliasAddCommand verifies the alias add command structure
func TestAliasAddCommand(t *testing.T) {
	if aliasAddCmd == nil {
		t.Fatal("alias add command is nil")
	}

	if aliasAddCmd.Short == "" {
		t.Error("alias add command should have a Short description")
	}

	if aliasAddCmd.Long == "" {
		t.Error("alias add command should have a Long description")
	}

	if aliasAddCmd.RunE == nil {
		t.Error("alias add command should have a RunE function")
	}
}

// TestAliasRemoveCommand verifies the alias remove command structure
func TestAliasRemoveCommand(t *testing.T) {
	if aliasRemoveCmd == nil {
		t.Fatal("alias remove command is nil")
	}

	if aliasRemoveCmd.Short == "" {
		t.Error("alias remove command should have a Short description")
	}

	if aliasRemoveCmd.Long == "" {
		t.Error("alias remove command should have a Long description")
	}

	if aliasRemoveCmd.RunE == nil {
		t.Error("alias remove command should have a RunE function")
	}
}

// TestAliasRemoveCommandHasFlags verifies the alias remove command has expected flags
func TestAliasRemoveCommandHasFlags(t *testing.T) {
	expectedFlags := []string{"force"}

	for _, flagName := range expectedFlags {
		flag := aliasRemoveCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("alias remove command missing flag: %s", flagName)
		}
	}
}
