package cmd

import (
	"testing"
)

// TestUserSuspendCommandExists verifies that the user suspend command is registered
func TestUserSuspendCommandExists(t *testing.T) {
	if userSuspendCmd == nil {
		t.Fatal("user suspend command is nil")
	}

	if userSuspendCmd.Use != "suspend <user-email>" {
		t.Errorf("Expected user suspend command Use to be 'suspend <user-email>', got '%s'", userSuspendCmd.Use)
	}

	if userSuspendCmd.Short == "" {
		t.Error("user suspend command should have a Short description")
	}

	if userSuspendCmd.Long == "" {
		t.Error("user suspend command should have a Long description")
	}

	if userSuspendCmd.RunE == nil {
		t.Error("user suspend command should have a RunE function")
	}
}

// TestUserUnsuspendCommandExists verifies that the user unsuspend command is registered
func TestUserUnsuspendCommandExists(t *testing.T) {
	if userUnsuspendCmd == nil {
		t.Fatal("user unsuspend command is nil")
	}

	if userUnsuspendCmd.Use != "unsuspend <user-email>" {
		t.Errorf("Expected user unsuspend command Use to be 'unsuspend <user-email>', got '%s'", userUnsuspendCmd.Use)
	}

	if userUnsuspendCmd.Short == "" {
		t.Error("user unsuspend command should have a Short description")
	}

	if userUnsuspendCmd.Long == "" {
		t.Error("user unsuspend command should have a Long description")
	}

	if userUnsuspendCmd.RunE == nil {
		t.Error("user unsuspend command should have a RunE function")
	}
}

// TestUserSuspendCommandHasFlags verifies the user suspend command has expected flags
func TestUserSuspendCommandHasFlags(t *testing.T) {
	expectedFlags := []string{"reason", "force"}

	for _, flagName := range expectedFlags {
		flag := userSuspendCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("user suspend command missing flag: %s", flagName)
		}
	}
}

// TestUserUnsuspendCommandHasFlags verifies the user unsuspend command has expected flags
func TestUserUnsuspendCommandHasFlags(t *testing.T) {
	expectedFlags := []string{"force"}

	for _, flagName := range expectedFlags {
		flag := userUnsuspendCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("user unsuspend command missing flag: %s", flagName)
		}
	}
}

// TestUserSuspendCommandRegistered verifies suspend command is in user subcommands
func TestUserSuspendCommandRegistered(t *testing.T) {
	commands := userCmd.Commands()
	found := false

	for _, cmd := range commands {
		if cmd.Name() == "suspend" {
			found = true
			break
		}
	}

	if !found {
		t.Error("suspend command not registered as user subcommand")
	}
}

// TestUserUnsuspendCommandRegistered verifies unsuspend command is in user subcommands
func TestUserUnsuspendCommandRegistered(t *testing.T) {
	commands := userCmd.Commands()
	found := false

	for _, cmd := range commands {
		if cmd.Name() == "unsuspend" {
			found = true
			break
		}
	}

	if !found {
		t.Error("unsuspend command not registered as user subcommand")
	}
}
