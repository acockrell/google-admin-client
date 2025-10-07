package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

// aliasAddCmd represents the alias add command
var aliasAddCmd = &cobra.Command{
	Use:   "add <user-email> <alias-email>",
	Short: "Add an alias to a user",
	Long: `Add an email alias to a Google Workspace user.

Usage
-----

$ gac alias add user@example.com support@example.com
$ gac alias add jdoe@example.com john@example.com
$ gac alias add admin@example.com administrator@example.com

Description
-----------

Adds an email alias to a user account. The alias must be in the same domain
or an alias domain of your Google Workspace organization.

After adding an alias, email sent to the alias address will be delivered to
the user's primary mailbox.

Examples:
  # Add a department alias
  gac alias add user@example.com support@example.com

  # Add an alternative name
  gac alias add john.doe@example.com jdoe@example.com

  # Add a role-based alias
  gac alias add admin@example.com administrator@example.com

Requirements:
  - The alias must be in a domain or alias domain managed by your organization
  - The alias cannot already be in use by another user or group
  - The user account must exist
`,
	Args: cobra.ExactArgs(2),
	RunE: aliasAddRunFunc,
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd)
}

func aliasAddRunFunc(cmd *cobra.Command, args []string) error {
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

	// Create the alias
	alias := &admin.Alias{
		Alias: aliasEmail,
	}

	result, err := client.Users.Aliases.Insert(userEmail, alias).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding alias: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - User does not exist\n")
		fmt.Fprintf(os.Stderr, "  - Alias already exists for another user or group\n")
		fmt.Fprintf(os.Stderr, "  - Alias domain is not managed by your organization\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		return err
	}

	fmt.Printf("Successfully added alias:\n\n")
	fmt.Printf("  User:  %s\n", userEmail)
	fmt.Printf("  Alias: %s\n", result.Alias)

	return nil
}
