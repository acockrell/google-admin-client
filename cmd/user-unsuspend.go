package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	unsuspendForce bool
)

// userUnsuspendCmd represents the user unsuspend command
var userUnsuspendCmd = &cobra.Command{
	Use:   "unsuspend <user-email>",
	Short: "Unsuspend (restore) a user account",
	Long: `Unsuspend a user account, restoring access to Google Workspace services.

Usage
-----

$ gac user unsuspend user@example.com
$ gac user unsuspend user@example.com --force

Description
-----------

Unsuspends a previously suspended user account, restoring the user's ability to:
- Sign in to their account
- Access all Google Workspace services (Gmail, Drive, Calendar, etc.)
- Send and receive emails
- Access their data and documents

All account data and settings are preserved during suspension and will be
available after unsuspending.

By default, the command will prompt for confirmation before unsuspending.
Use --force to skip the confirmation prompt.

Examples:
  # Unsuspend a user with confirmation
  gac user unsuspend user@example.com

  # Unsuspend without confirmation
  gac user unsuspend user@example.com --force

Common use cases:
  - Restoring access after employee returns from leave
  - Correcting accidental suspensions
  - Restoring accounts after security incidents are resolved
  - Re-enabling accounts after policy violations are addressed
`,
	Args: cobra.ExactArgs(1),
	RunE: userUnsuspendRunFunc,
}

func init() {
	userCmd.AddCommand(userUnsuspendCmd)
	userUnsuspendCmd.Flags().BoolVarP(&unsuspendForce, "force", "f", false, "skip confirmation prompt")
}

func userUnsuspendRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	userEmail := args[0]

	// Validate email format
	if err := ValidateEmail(userEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid email format: %v\n", err)
		return err
	}

	// Show warning and prompt for confirmation unless --force is used
	if !unsuspendForce {
		fmt.Printf("You are about to unsuspend user account: %s\n", userEmail)
		fmt.Printf("\nThe user will regain:\n")
		fmt.Printf("  - Ability to sign in to their account\n")
		fmt.Printf("  - Access to all Google Workspace services\n")
		fmt.Printf("  - Ability to send and receive emails\n")
		fmt.Printf("\nType 'yes' to confirm unsuspension: ")

		var confirmation string
		_, err := fmt.Scanln(&confirmation)
		if err != nil {
			// If there's an error reading input (e.g., EOF), treat as cancellation
			fmt.Fprintf(os.Stderr, "\nUnsuspension cancelled.\n")
			return nil
		}

		if confirmation != "yes" {
			fmt.Println("Unsuspension cancelled.")
			return nil
		}
	}

	// Unsuspend the user
	user := &admin.User{
		Suspended: false,
	}

	// Force send the Suspended field
	user.ForceSendFields = []string{"Suspended"}

	result, err := client.Users.Update(userEmail, user).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unsuspending user: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - User does not exist\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		fmt.Fprintf(os.Stderr, "  - User is already active (not suspended)\n")
		return err
	}

	fmt.Printf("Successfully unsuspended user account:\n\n")
	fmt.Printf("  Email: %s\n", result.PrimaryEmail)
	fmt.Printf("  Name: %s %s\n", result.Name.GivenName, result.Name.FamilyName)
	fmt.Printf("  Suspended: %v\n", result.Suspended)
	fmt.Printf("\nThe user can now sign in and access Google Workspace services.\n")

	return nil
}
