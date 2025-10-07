package cmd

import (
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "User alias operations",
	Long: `Manage email aliases for Google Workspace users.

Email aliases allow users to receive mail at multiple addresses that all
deliver to the same mailbox. This is useful for department addresses,
role-based addresses, or alternative names.

Available Commands:
  list    - List aliases for a user
  add     - Add an alias to a user
  remove  - Remove an alias from a user

Examples:
  # List all aliases for a user
  gac alias list user@example.com

  # Add an alias to a user
  gac alias add user@example.com support@example.com

  # Remove an alias from a user
  gac alias remove user@example.com old-alias@example.com

For more information on a specific command, use:
  gac alias [command] --help
`,
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
