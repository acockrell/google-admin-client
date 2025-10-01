package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestAllCommandsRegistered(t *testing.T) {
	commands := map[string]bool{
		"user":     false,
		"group":    false,
		"calendar": false,
		"transfer": false,
		"init":     false,
	}

	for _, cmd := range rootCmd.Commands() {
		if _, exists := commands[cmd.Use]; exists {
			commands[cmd.Use] = true
		}
	}

	for cmdName, registered := range commands {
		if !registered {
			t.Errorf("Command %q is not registered with rootCmd", cmdName)
		}
	}

	// Note: completion command is added by Cobra automatically
}

func TestUserSubcommands(t *testing.T) {
	expectedSubcommands := []string{"create", "list", "update"}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range userCmd.Commands() {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("user subcommand %q not found", expected)
		}
	}
}

func TestCalendarSubcommands(t *testing.T) {
	expectedSubcommands := []string{"create", "list", "update"}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range calendarCmd.Commands() {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("calendar subcommand %q not found", expected)
		}
	}
}

func TestGroupSubcommands(t *testing.T) {
	expectedSubcommands := []string{"list"}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range groupCmd.Commands() {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("group subcommand %q not found", expected)
		}
	}
}

func TestCommandsHaveShortDescriptions(t *testing.T) {
	commands := []*cobra.Command{
		rootCmd,
		userCmd,
		groupCmd,
		calendarCmd,
		transferCmd,
		initCmd,
	}

	for _, cmd := range commands {
		if cmd.Short == "" {
			t.Errorf("Command %q missing short description", cmd.Use)
		}
	}
}

func TestUserCreateCommandHasFlags(t *testing.T) {
	expectedFlags := []string{"groups", "email", "first-name", "last-name"}

	for _, flagName := range expectedFlags {
		flag := createUserCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("user create command missing flag: --%s", flagName)
		}
	}
}

func TestUserUpdateCommandHasFlags(t *testing.T) {
	expectedFlags := []string{
		"address", "dept", "group", "id", "force", "manager",
		"ou", "phone", "remove", "title", "type",
	}

	for _, flagName := range expectedFlags {
		flag := updateUserCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("user update command missing flag: --%s", flagName)
		}
	}
}

func TestUserListCommandHasFlags(t *testing.T) {
	expectedFlags := []string{"full", "csv", "disabled-only"}

	for _, flagName := range expectedFlags {
		flag := listUserCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("user list command missing flag: --%s", flagName)
		}
	}
}
