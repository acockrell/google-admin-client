package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// groupSettingsListCmd represents the group-settings list command
var groupSettingsListCmd = &cobra.Command{
	Use:   "list [group-email]",
	Short: "View group settings",
	Long: `View settings for a Google Workspace group.

Usage
-----

$ gac group-settings list operations@example.com
$ gac group-settings list engineering
$ gac group-settings list sales@example.com --format json

Description
-----------

Displays the settings for a group, including:
- Who can join the group
- Who can view group messages
- Who can post messages
- Message moderation settings
- Email delivery preferences
- Archive settings
- Custom footer text

Use the global --format flag to control output format (json, csv, yaml, table, plain).

If you don't include the @domain part in the group email, the configured domain
will be automatically appended.
`,
	Args: cobra.ExactArgs(1),
	RunE: groupSettingsListRunFunc,
}

func init() {
	groupSettingsCmd.AddCommand(groupSettingsListCmd)
}

func groupSettingsListRunFunc(cmd *cobra.Command, args []string) error {
	groupEmail := args[0]
	if !strings.Contains(groupEmail, "@") {
		groupEmail = groupEmail + "@" + getDomain()
	}

	// Validate email
	if err := ValidateEmail(groupEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid group email: %v\n", err)
		return err
	}

	client, err := newGroupsSettingsClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	settings, err := client.Groups.Get(groupEmail).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting group settings for %s: %v\n", groupEmail, err)
		return err
	}

	// Use the unified formatter for output
	if err := FormatOutput(settings, nil); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
