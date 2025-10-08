package cmd

import (
	"testing"
)

func TestConfirmActionWithGlobalFlag(t *testing.T) {
	// Save original value
	origSkipConfirmations := skipConfirmations

	// Test with global --yes flag
	skipConfirmations = true
	result := confirmAction("Test message", false)
	if !result {
		t.Error("confirmAction should return true when global --yes flag is set")
	}

	// Restore original value
	skipConfirmations = origSkipConfirmations
}

func TestConfirmActionWithForceFlag(t *testing.T) {
	// Save original value
	origSkipConfirmations := skipConfirmations
	skipConfirmations = false

	// Test with command-specific --force flag
	result := confirmAction("Test message", true)
	if !result {
		t.Error("confirmAction should return true when force flag is set")
	}

	// Restore original value
	skipConfirmations = origSkipConfirmations
}

func TestConfirmActionPriority(t *testing.T) {
	// Save original value
	origSkipConfirmations := skipConfirmations

	// Test that global --yes takes precedence
	skipConfirmations = true
	result := confirmAction("Test message", false)
	if !result {
		t.Error("confirmAction should return true when global --yes is set, regardless of force flag")
	}

	// Test that force flag works when global --yes is false
	skipConfirmations = false
	result = confirmAction("Test message", true)
	if !result {
		t.Error("confirmAction should return true when force is true, even if global --yes is false")
	}

	// Restore original value
	skipConfirmations = origSkipConfirmations
}

func TestConfirmDeletion(t *testing.T) {
	// Save original value
	origSkipConfirmations := skipConfirmations

	tests := []struct {
		name           string
		resourceType   string
		resourceName   string
		additionalInfo string
		force          bool
		globalYes      bool
		expectedResult bool
		description    string
	}{
		{
			name:           "skip with global yes",
			resourceType:   "test resource",
			resourceName:   "test-123",
			additionalInfo: "additional info",
			force:          false,
			globalYes:      true,
			expectedResult: true,
			description:    "should skip confirmation with global --yes",
		},
		{
			name:           "skip with force flag",
			resourceType:   "test resource",
			resourceName:   "test-456",
			additionalInfo: "",
			force:          true,
			globalYes:      false,
			expectedResult: true,
			description:    "should skip confirmation with --force",
		},
		{
			name:           "both flags set",
			resourceType:   "test resource",
			resourceName:   "test-789",
			additionalInfo: "info",
			force:          true,
			globalYes:      true,
			expectedResult: true,
			description:    "should skip confirmation when both flags are set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipConfirmations = tt.globalYes

			result := confirmDeletion(tt.resourceType, tt.resourceName, tt.additionalInfo, tt.force)

			if result != tt.expectedResult {
				t.Errorf("confirmDeletion() = %v, want %v (%s)", result, tt.expectedResult, tt.description)
			}
		})
	}

	// Restore original value
	skipConfirmations = origSkipConfirmations
}

func TestConfirmDeletionMessageFormats(t *testing.T) {
	// Save original value
	origSkipConfirmations := skipConfirmations

	tests := []struct {
		name           string
		resourceType   string
		resourceName   string
		additionalInfo string
	}{
		{
			name:           "with additional info",
			resourceType:   "organizational unit",
			resourceName:   "/Engineering/Test",
			additionalInfo: "This OU must be empty",
		},
		{
			name:           "without additional info",
			resourceType:   "calendar resource",
			resourceName:   "room-101",
			additionalInfo: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with force=true to skip actual user input
			skipConfirmations = false
			result := confirmDeletion(tt.resourceType, tt.resourceName, tt.additionalInfo, true)

			if !result {
				t.Errorf("confirmDeletion() with force=true should return true")
			}
		})
	}

	// Restore original value
	skipConfirmations = origSkipConfirmations
}
