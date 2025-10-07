package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	aliasRemoveForce bool
)

// aliasRemoveCmd represents the alias remove command
var aliasRemoveCmd = &cobra.Command{
	Use:   "remove <user-email> <alias-email>",
	Short: "Remove an alias from a user",
	Long: `Remove an email alias from a Google Workspace user.

Usage
-----

$ gac alias remove user@example.com old-alias@example.com
$ gac alias remove jdoe@example.com john@example.com --force

Description
-----------

Removes an email alias from a user account. After removal, email sent to
the alias address will no longer be delivered to the user's mailbox.

By default, the command will prompt for confirmation before removing the
alias. Use --force to skip the confirmation prompt.

Examples:
  # Remove an alias with confirmation
  gac alias remove user@example.com old-alias@example.com

  # Remove an alias without confirmation
  gac alias remove user@example.com support@example.com --force

WARNING: After removal, the alias address will no longer deliver mail to
this user. Make sure this is intentional before proceeding.
`,
	Args: cobra.ExactArgs(2),
	RunE: aliasRemoveRunFunc,
}

func init() {
	aliasCmd.AddCommand(aliasRemoveCmd)
	aliasRemoveCmd.Flags().BoolVarP(&aliasRemoveForce, "force", "f", false, "skip confirmation prompt")
}

func aliasRemoveRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	userEmail := args[0]
	aliasEmail := args[1]

	// Validate email formats
	if err := ValidateEmail(userEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid user email format: %v\n", err)
		return err
	}

	if err := ValidateEmail(aliasEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid alias email format: %v\n", err)
		return err
	}

	// Show warning and prompt for confirmation unless --force is used
	if !aliasRemoveForce {
		fmt.Printf("WARNING: You are about to remove alias:\n")
		fmt.Printf("  User:  %s\n", userEmail)
		fmt.Printf("  Alias: %s\n\n", aliasEmail)
		fmt.Printf("After removal, email to this alias will no longer be delivered.\n\n")
		fmt.Printf("Type 'yes' to confirm removal: ")

		var confirmation string
		_, err := fmt.Scanln(&confirmation)
		if err != nil {
			// If there's an error reading input (e.g., EOF), treat as cancellation
			fmt.Fprintf(os.Stderr, "\nRemoval cancelled.\n")
			return nil
		}

		if confirmation != "yes" {
			fmt.Println("Removal cancelled.")
			return nil
		}
	}

	// Remove the alias
	err = client.Users.Aliases.Delete(userEmail, aliasEmail).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error removing alias: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - User does not exist\n")
		fmt.Fprintf(os.Stderr, "  - Alias does not exist for this user\n")
		fmt.Fprintf(os.Stderr, "  - Alias email is incorrect\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		return err
	}

	fmt.Printf("Successfully removed alias:\n")
	fmt.Printf("  User:  %s\n", userEmail)
	fmt.Printf("  Alias: %s\n", aliasEmail)

	return nil
}
