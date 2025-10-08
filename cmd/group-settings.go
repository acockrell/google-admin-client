package cmd

import (
	"github.com/spf13/cobra"
)

// groupSettingsCmd represents the group-settings command
var groupSettingsCmd = &cobra.Command{
	Use:   "group-settings",
	Short: "Manage group settings",
	Long: `Manage group settings in your Google Workspace domain.

Group settings control various aspects of how a group operates, including:
- Who can join and view the group
- Who can post messages
- Message moderation settings
- Email delivery preferences
- Group archiving settings
- Custom footer text

Available Commands:
  list    - View group settings
  update  - Update group settings

Use "gac group-settings [command] --help" for more information about a command.
`,
}

func init() {
	rootCmd.AddCommand(groupSettingsCmd)
}
