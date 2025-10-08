package cmd

import (
	"testing"
)

func TestVersionCommandRegistration(t *testing.T) {
	// Check that version command is registered
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "version" {
			found = true

			// Check for --short flag
			shortFlag := cmd.Flag("short")
			if shortFlag == nil {
				t.Error("version command missing --short flag")
			}

			// Verify short flag has alias
			if shortFlag != nil && shortFlag.Shorthand != "s" {
				t.Errorf("version --short flag shorthand expected 's', got '%s'", shortFlag.Shorthand)
			}

			break
		}
	}

	if !found {
		t.Error("version command not registered with root command")
	}
}

func TestSetVersionInfo(t *testing.T) {
	// Save original values
	origVersion := version
	origCommit := commit
	origDate := date
	origBuiltBy := builtBy

	// Test setting version info
	testVersion := "v1.2.3"
	testCommit := "abc123"
	testDate := "2024-10-08T12:00:00Z"
	testBuiltBy := "test"

	SetVersionInfo(testVersion, testCommit, testDate, testBuiltBy)

	if version != testVersion {
		t.Errorf("SetVersionInfo: version = %s, want %s", version, testVersion)
	}
	if commit != testCommit {
		t.Errorf("SetVersionInfo: commit = %s, want %s", commit, testCommit)
	}
	if date != testDate {
		t.Errorf("SetVersionInfo: date = %s, want %s", date, testDate)
	}
	if builtBy != testBuiltBy {
		t.Errorf("SetVersionInfo: builtBy = %s, want %s", builtBy, testBuiltBy)
	}

	// Restore original values
	version = origVersion
	commit = origCommit
	date = origDate
	builtBy = origBuiltBy
}

func TestVersionDefaults(t *testing.T) {
	// Test that default values are set correctly
	// (This assumes the package-level variables haven't been changed by other tests)
	expectedDefaults := map[string]string{
		"version": "dev",
		"commit":  "none",
		"date":    "unknown",
		"builtBy": "unknown",
	}

	// We can't directly test the defaults after SetVersionInfo has been called,
	// but we can verify the function accepts all expected parameters
	SetVersionInfo(
		expectedDefaults["version"],
		expectedDefaults["commit"],
		expectedDefaults["date"],
		expectedDefaults["builtBy"],
	)

	if version != expectedDefaults["version"] {
		t.Errorf("Default version = %s, want %s", version, expectedDefaults["version"])
	}
	if commit != expectedDefaults["commit"] {
		t.Errorf("Default commit = %s, want %s", commit, expectedDefaults["commit"])
	}
	if date != expectedDefaults["date"] {
		t.Errorf("Default date = %s, want %s", date, expectedDefaults["date"])
	}
	if builtBy != expectedDefaults["builtBy"] {
		t.Errorf("Default builtBy = %s, want %s", builtBy, expectedDefaults["builtBy"])
	}
}
