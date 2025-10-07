package cmd

import (
	"github.com/spf13/cobra"
)

// calResourceCmd represents the cal-resource command
var calResourceCmd = &cobra.Command{
	Use:   "cal-resource",
	Short: "Manage calendar resources",
	Long: `Manage calendar resources in your Google Workspace domain.

Calendar resources are bookable items such as conference rooms, equipment,
and other shared resources. This command allows you to list, create, update,
and delete calendar resources.

Available Commands:
  list    - List calendar resources
  create  - Create a new calendar resource
  update  - Update an existing calendar resource
  delete  - Delete a calendar resource

Use "gac cal-resource [command] --help" for more information about a command.
`,
}

func init() {
	rootCmd.AddCommand(calResourceCmd)
}
