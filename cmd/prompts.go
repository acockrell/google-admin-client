package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Global yes flag (set in root.go)
var skipConfirmations bool

// confirmAction prompts the user for confirmation of a destructive action.
// Returns true if the user confirmed (typed "yes"), false otherwise.
//
// The confirmation is automatically skipped if:
// - The global --yes flag is set (skipConfirmations = true)
// - The force parameter is true (command-specific --force flag)
//
// Parameters:
//   - message: The warning message to display to the user
//   - force: Whether to skip the prompt (command-specific --force flag)
//
// Returns:
//   - true if confirmed (user typed "yes", or confirmation was skipped)
//   - false if not confirmed (user typed anything else or cancelled)
func confirmAction(message string, force bool) bool {
	// Check global --yes flag first
	if skipConfirmations {
		LogDebug("Confirmation skipped via global --yes flag", map[string]interface{}{
			"message": message,
		})
		return true
	}

	// Check command-specific --force flag
	if force {
		LogDebug("Confirmation skipped via --force flag", map[string]interface{}{
			"message": message,
		})
		return true
	}

	// Display warning message
	fmt.Print(message)
	fmt.Print("\nType 'yes' to confirm: ")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		// If there's an error reading input (e.g., EOF, pipe closed),
		// treat as cancellation for safety
		LogDebug("Error reading confirmation input", map[string]interface{}{
			"error": err.Error(),
		})
		fmt.Fprintf(os.Stderr, "\nAction cancelled.\n")
		return false
	}

	// Trim whitespace and compare
	response = strings.TrimSpace(response)
	confirmed := response == "yes"

	if !confirmed {
		fmt.Println("Action cancelled.")
		LogDebug("User cancelled action", map[string]interface{}{
			"response": response,
		})
	} else {
		LogDebug("User confirmed action", nil)
	}

	return confirmed
}

// confirmDeletion is a specialized confirmation for deletion operations.
// It displays a standardized deletion warning with the resource details.
//
// Parameters:
//   - resourceType: The type of resource being deleted (e.g., "organizational unit", "calendar resource")
//   - resourceName: The name/identifier of the resource
//   - additionalInfo: Optional additional information to display (can be empty string)
//   - force: Whether to skip the prompt (command-specific --force flag)
//
// Returns:
//   - true if confirmed, false otherwise
func confirmDeletion(resourceType, resourceName, additionalInfo string, force bool) bool {
	var message strings.Builder

	message.WriteString(fmt.Sprintf("WARNING: You are about to delete %s: %s\n", resourceType, resourceName))
	message.WriteString("This operation cannot be undone.\n")

	if additionalInfo != "" {
		message.WriteString("\n")
		message.WriteString(additionalInfo)
		message.WriteString("\n")
	}

	return confirmAction(message.String(), force)
}
