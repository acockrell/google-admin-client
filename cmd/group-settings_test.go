package cmd

import (
	"testing"
)

func TestGroupSettingsCommandExists(t *testing.T) {
	if groupSettingsCmd == nil {
		t.Fatal("Group-settings command should not be nil")
	}

	if groupSettingsCmd.Use != "group-settings" {
		t.Errorf("Group-settings command Use = %q, want %q", groupSettingsCmd.Use, "group-settings")
	}

	if groupSettingsCmd.Short == "" {
		t.Error("Group-settings command should have a short description")
	}
}

func TestGroupSettingsSubcommands(t *testing.T) {
	expectedCommands := []string{"list", "update"}

	commands := groupSettingsCmd.Commands()
	commandMap := make(map[string]bool)

	for _, cmd := range commands {
		// cmd.Use might be "list [group-email]" so extract just the first word
		use := cmd.Name()
		commandMap[use] = true
	}

	for _, cmdName := range expectedCommands {
		if !commandMap[cmdName] {
			t.Errorf("Group-settings command missing subcommand: %s", cmdName)
		}
	}
}

func TestGroupSettingsListCommandHasFlags(t *testing.T) {
	if groupSettingsListCmd == nil {
		t.Fatal("Group-settings list command should not be nil")
	}

	// Check for --format flag
	formatFlag := groupSettingsListCmd.Flags().Lookup("format")
	if formatFlag == nil {
		t.Error("Group-settings list command should have --format flag")
	}
}

func TestGroupSettingsUpdateCommandHasFlags(t *testing.T) {
	if groupSettingsUpdateCmd == nil {
		t.Fatal("Group-settings update command should not be nil")
	}

	expectedFlags := []string{
		"who-can-join",
		"who-can-view-group",
		"who-can-view-membership",
		"who-can-post-message",
		"allow-external-members",
		"allow-web-posting",
		"archive-only",
		"message-moderation-level",
		"spam-moderation-level",
		"reply-to",
		"custom-reply-to",
		"custom-footer-text",
		"include-custom-footer",
		"send-message-deny-notification",
		"include-in-global-address-list",
		"show-in-group-directory",
		"who-can-leave-group",
		"who-can-add",
		"who-can-invite",
		"who-can-approve-members",
		"allow-google-communication",
		"members-can-post-as-the-group",
		"who-can-contact-owner",
		"who-can-moderate-members",
		"who-can-moderate-content",
		"who-can-ban-users",
	}

	for _, flagName := range expectedFlags {
		flag := groupSettingsUpdateCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Group-settings update command should have --%s flag", flagName)
		}
	}
}

func TestGroupSettingsCommandIsRegistered(t *testing.T) {
	// Check if group-settings command is registered with root command
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "group-settings" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Group-settings command should be registered with root command")
	}
}

func TestGroupSettingsListCommandStructure(t *testing.T) {
	if groupSettingsListCmd.Use != "list [group-email]" {
		t.Errorf("Group-settings list command Use = %q, want %q", groupSettingsListCmd.Use, "list [group-email]")
	}

	if groupSettingsListCmd.Short == "" {
		t.Error("Group-settings list command should have a short description")
	}

	if groupSettingsListCmd.Long == "" {
		t.Error("Group-settings list command should have a long description")
	}

	if groupSettingsListCmd.RunE == nil {
		t.Error("Group-settings list command should have a RunE function")
	}
}

func TestGroupSettingsUpdateCommandStructure(t *testing.T) {
	if groupSettingsUpdateCmd.Use != "update [group-email]" {
		t.Errorf("Group-settings update command Use = %q, want %q", groupSettingsUpdateCmd.Use, "update [group-email]")
	}

	if groupSettingsUpdateCmd.Short == "" {
		t.Error("Group-settings update command should have a short description")
	}

	if groupSettingsUpdateCmd.Long == "" {
		t.Error("Group-settings update command should have a long description")
	}

	if groupSettingsUpdateCmd.RunE == nil {
		t.Error("Group-settings update command should have a RunE function")
	}
}
