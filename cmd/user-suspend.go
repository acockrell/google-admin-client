package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	suspendReason string
	suspendForce  bool
)

// userSuspendCmd represents the user suspend command
var userSuspendCmd = &cobra.Command{
	Use:   "suspend <user-email>",
	Short: "Suspend a user account",
	Long: `Suspend a user account, preventing the user from logging in.

Usage
-----

$ gac user suspend user@example.com
$ gac user suspend user@example.com --reason "Policy violation"
$ gac user suspend user@example.com --force

Description
-----------

Suspends a user account, which prevents the user from:
- Signing in to their account
- Accessing any Google Workspace services (Gmail, Drive, Calendar, etc.)
- Receiving new emails (emails will bounce)

The account data is preserved and can be restored by unsuspending the account.

By default, the command will prompt for confirmation before suspending.
Use --force to skip the confirmation prompt.

Examples:
  # Suspend a user with confirmation
  gac user suspend user@example.com

  # Suspend with a reason
  gac user suspend user@example.com --reason "Left company"

  # Suspend without confirmation
  gac user suspend user@example.com --force

WARNING: Suspended users cannot access any Google Workspace services.
Make sure this is intentional before proceeding.

Common use cases:
  - Employee termination or departure
  - Policy violations or security incidents
  - Account compromise or suspicious activity
  - Extended leave or sabbatical
`,
	Args: cobra.ExactArgs(1),
	RunE: userSuspendRunFunc,
}

func init() {
	userCmd.AddCommand(userSuspendCmd)
	userSuspendCmd.Flags().StringVarP(&suspendReason, "reason", "r", "", "reason for suspension")
	userSuspendCmd.Flags().BoolVarP(&suspendForce, "force", "f", false, "skip confirmation prompt")
}

func userSuspendRunFunc(cmd *cobra.Command, args []string) error {
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
	if !suspendForce {
		fmt.Printf("WARNING: You are about to suspend user account: %s\n", userEmail)
		fmt.Printf("\nThe user will:\n")
		fmt.Printf("  - Be unable to sign in to their account\n")
		fmt.Printf("  - Lose access to all Google Workspace services\n")
		fmt.Printf("  - Not receive new emails (emails will bounce)\n")
		if suspendReason != "" {
			fmt.Printf("\nReason: %s\n", suspendReason)
		}
		fmt.Printf("\nType 'yes' to confirm suspension: ")

		var confirmation string
		_, err := fmt.Scanln(&confirmation)
		if err != nil {
			// If there's an error reading input (e.g., EOF), treat as cancellation
			fmt.Fprintf(os.Stderr, "\nSuspension cancelled.\n")
			return nil
		}

		if confirmation != "yes" {
			fmt.Println("Suspension cancelled.")
			return nil
		}
	}

	// Suspend the user
	user := &admin.User{
		Suspended: true,
	}

	// Add suspension reason if provided
	if suspendReason != "" {
		user.SuspensionReason = suspendReason
	}

	// Force send the Suspended field
	user.ForceSendFields = []string{"Suspended"}

	result, err := client.Users.Update(userEmail, user).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error suspending user: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - User does not exist\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		fmt.Fprintf(os.Stderr, "  - User is already suspended\n")
		fmt.Fprintf(os.Stderr, "  - Super admin accounts may have restrictions\n")
		return err
	}

	fmt.Printf("Successfully suspended user account:\n\n")
	fmt.Printf("  Email: %s\n", result.PrimaryEmail)
	fmt.Printf("  Name: %s %s\n", result.Name.GivenName, result.Name.FamilyName)
	fmt.Printf("  Suspended: %v\n", result.Suspended)
	if result.SuspensionReason != "" {
		fmt.Printf("  Reason: %s\n", result.SuspensionReason)
	}

	return nil
}
