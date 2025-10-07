package cmd

import (
	"github.com/spf13/cobra"
)

// ouCmd represents the ou command
var ouCmd = &cobra.Command{
	Use:   "ou",
	Short: "Organizational unit operations",
	Long: `Manage Google Workspace organizational units.

Organizational units (OUs) allow you to organize users and apply different
policies to different groups of users. This command provides operations to
list, create, update, and delete organizational units.

Examples:
  gac ou list
  gac ou create /Engineering
  gac ou update /Engineering --description "Engineering team"
  gac ou delete /Engineering/Archived
`,
}

func init() {
	rootCmd.AddCommand(ouCmd)
}
