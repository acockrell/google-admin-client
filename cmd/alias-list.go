package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// aliasListCmd represents the alias list command
var aliasListCmd = &cobra.Command{
	Use:   "list <user-email>",
	Short: "List aliases for a user",
	Long: `List all email aliases for a Google Workspace user.

Usage
-----

$ gac alias list user@example.com
$ gac alias list jdoe@example.com

Description
-----------

Lists all email aliases associated with a user account. Each alias is an
alternative email address that delivers to the same mailbox as the primary
user email.

Examples:
  # List all aliases for a user
  gac alias list user@example.com

  # List aliases for a specific user
  gac alias list john.doe@example.com

The output shows:
  - The primary user email
  - All configured aliases
  - Total count of aliases
`,
	Args: cobra.ExactArgs(1),
	RunE: aliasListRunFunc,
}

func init() {
	aliasCmd.AddCommand(aliasListCmd)
}

func aliasListRunFunc(cmd *cobra.Command, args []string) error {
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

	// List aliases for the user
	result, err := client.Users.Aliases.List(userEmail).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing aliases: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - User does not exist\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		fmt.Fprintf(os.Stderr, "  - Invalid user email\n")
		return err
	}

	// Display results
	QuietPrintf("Aliases for %s:\n\n", userEmail)

	if len(result.Aliases) == 0 {
		QuietPrintln("No aliases found.")
		return nil
	}

	// Convert aliases to simple struct for formatting
	type aliasItem struct {
		Alias string `json:"alias"`
	}

	var aliases []aliasItem
	for _, aliasInterface := range result.Aliases {
		// The Aliases field is []interface{}, so we need to type assert
		if aliasMap, ok := aliasInterface.(map[string]interface{}); ok {
			if aliasEmail, ok := aliasMap["alias"].(string); ok {
				aliases = append(aliases, aliasItem{Alias: aliasEmail})
			}
		}
	}

	headers := []string{"Alias"}
	if err := FormatOutput(aliases, headers); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	QuietPrintf("\nTotal: %d alias(es)\n", len(result.Aliases))

	return nil
}
