package cmd

import (
	"github.com/spf13/cobra"
)

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Export and manage audit logs",
	Long: `Export and manage Google Workspace audit logs.

The audit command provides access to Google Workspace audit logs through the Admin Reports API.
You can export audit logs for various applications including admin console, login, drive, calendar,
groups, and more.

Examples:
  # Export admin console audit logs
  gac audit export --app admin

  # Export login activities for a specific user
  gac audit export --app login --user user@example.com

  # Export drive activities to CSV
  gac audit export --app drive --output csv --output-file drive-audit.csv
`,
}

func init() {
	rootCmd.AddCommand(auditCmd)
}
