package cmd

import (
	"testing"
)

func TestOUCommandExists(t *testing.T) {
	if ouCmd == nil {
		t.Fatal("OU command should not be nil")
	}

	if ouCmd.Use != "ou" {
		t.Errorf("OU command Use = %q, want %q", ouCmd.Use, "ou")
	}

	if ouCmd.Short == "" {
		t.Error("OU command should have a short description")
	}
}

func TestOUSubcommands(t *testing.T) {
	expectedCommands := []string{"list", "create", "update", "delete"}

	commands := ouCmd.Commands()
	commandMap := make(map[string]bool)

	for _, cmd := range commands {
		// cmd.Use might be "list [ou-path]" so extract just the first word
		use := cmd.Name()
		commandMap[use] = true
	}

	for _, cmdName := range expectedCommands {
		if !commandMap[cmdName] {
			t.Errorf("OU command missing subcommand: %s", cmdName)
		}
	}
}

func TestOUListCommandHasFlags(t *testing.T) {
	if ouListCmd == nil {
		t.Fatal("OU list command should not be nil")
	}

	// Check for --type flag
	typeFlag := ouListCmd.Flags().Lookup("type")
	if typeFlag == nil {
		t.Error("OU list command should have --type flag")
	}
}

func TestOUCreateCommandHasFlags(t *testing.T) {
	if ouCreateCmd == nil {
		t.Fatal("OU create command should not be nil")
	}

	expectedFlags := []string{"description", "parent", "block-inheritance"}

	for _, flagName := range expectedFlags {
		flag := ouCreateCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("OU create command should have --%s flag", flagName)
		}
	}
}

func TestOUUpdateCommandHasFlags(t *testing.T) {
	if ouUpdateCmd == nil {
		t.Fatal("OU update command should not be nil")
	}

	expectedFlags := []string{"name", "description", "parent", "block-inheritance"}

	for _, flagName := range expectedFlags {
		flag := ouUpdateCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("OU update command should have --%s flag", flagName)
		}
	}
}

func TestOUDeleteCommandHasFlags(t *testing.T) {
	if ouDeleteCmd == nil {
		t.Fatal("OU delete command should not be nil")
	}

	// Check for --force flag
	forceFlag := ouDeleteCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("OU delete command should have --force flag")
	}
}
