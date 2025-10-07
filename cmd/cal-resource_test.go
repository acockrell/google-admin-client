package cmd

import (
	"testing"
)

func TestCalResourceCommandExists(t *testing.T) {
	if calResourceCmd == nil {
		t.Fatal("Cal-resource command should not be nil")
	}

	if calResourceCmd.Use != "cal-resource" {
		t.Errorf("Cal-resource command Use = %q, want %q", calResourceCmd.Use, "cal-resource")
	}

	if calResourceCmd.Short == "" {
		t.Error("Cal-resource command should have a short description")
	}
}

func TestCalResourceSubcommands(t *testing.T) {
	expectedCommands := []string{"list", "create", "update", "delete"}

	commands := calResourceCmd.Commands()
	commandMap := make(map[string]bool)

	for _, cmd := range commands {
		// cmd.Use might be "list [resource-id]" so extract just the first word
		use := cmd.Name()
		commandMap[use] = true
	}

	for _, cmdName := range expectedCommands {
		if !commandMap[cmdName] {
			t.Errorf("Cal-resource command missing subcommand: %s", cmdName)
		}
	}
}

func TestCalResourceListCommandHasFlags(t *testing.T) {
	if calResourceListCmd == nil {
		t.Fatal("Cal-resource list command should not be nil")
	}

	// Check for --type flag
	typeFlag := calResourceListCmd.Flags().Lookup("type")
	if typeFlag == nil {
		t.Error("Cal-resource list command should have --type flag")
	}
}

func TestCalResourceCreateCommandHasFlags(t *testing.T) {
	if calResourceCreateCmd == nil {
		t.Fatal("Cal-resource create command should not be nil")
	}

	expectedFlags := []string{
		"name", "type", "description", "category", "building-id",
		"floor", "section", "capacity", "user-description",
	}

	for _, flagName := range expectedFlags {
		flag := calResourceCreateCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Cal-resource create command should have --%s flag", flagName)
		}
	}
}

func TestCalResourceUpdateCommandHasFlags(t *testing.T) {
	if calResourceUpdateCmd == nil {
		t.Fatal("Cal-resource update command should not be nil")
	}

	expectedFlags := []string{
		"name", "description", "category", "building-id",
		"floor", "section", "capacity", "user-description",
	}

	for _, flagName := range expectedFlags {
		flag := calResourceUpdateCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Cal-resource update command should have --%s flag", flagName)
		}
	}
}

func TestCalResourceDeleteCommandHasFlags(t *testing.T) {
	if calResourceDeleteCmd == nil {
		t.Fatal("Cal-resource delete command should not be nil")
	}

	// Check for --force flag
	forceFlag := calResourceDeleteCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("Cal-resource delete command should have --force flag")
	}
}

func TestCalResourceCommandIsRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "cal-resource" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Cal-resource command should be registered with root command")
	}
}
