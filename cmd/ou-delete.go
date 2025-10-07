package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ouDeleteForce bool
)

// ouDeleteCmd represents the ou delete command
var ouDeleteCmd = &cobra.Command{
	Use:   "delete <ou-path>",
	Short: "Delete an organizational unit",
	Long: `Delete an organizational unit from your Google Workspace domain.

Usage
-----

$ gac ou delete /Engineering/Archived
$ gac ou delete /TempOU --force

Description
-----------

Deletes an organizational unit. The OU must be empty (contain no users)
before it can be deleted. If the OU contains sub-OUs, those must also be
empty or deleted first.

WARNING: This operation cannot be undone. Use with caution.

By default, the command will fail if the OU is not empty. Use --force
to attempt deletion (though the API will still reject if users are present).

Examples:
  # Delete an empty OU
  gac ou delete /Engineering/Archived

  # Force delete (will still fail if OU contains users)
  gac ou delete /TempOU --force

Before deleting an OU:
  1. Move all users to a different OU
  2. Delete or move any child OUs
  3. Verify the OU is empty with: gac ou list /path/to/ou

`,
	Args: cobra.ExactArgs(1),
	RunE: ouDeleteRunFunc,
}

func init() {
	ouCmd.AddCommand(ouDeleteCmd)
	ouDeleteCmd.Flags().BoolVarP(&ouDeleteForce, "force", "f", false, "skip confirmation prompt")
}

func ouDeleteRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	ouPath := args[0]
	customerID := "my_customer"

	// Show warning and prompt for confirmation unless --force is used
	if !ouDeleteForce {
		fmt.Printf("WARNING: You are about to delete organizational unit: %s\n", ouPath)
		fmt.Printf("This operation cannot be undone.\n\n")
		fmt.Printf("The OU must be empty (no users) to be deleted.\n\n")
		fmt.Printf("Type 'yes' to confirm deletion: ")

		var confirmation string
		fmt.Scanln(&confirmation)

		if confirmation != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the organizational unit
	err = client.Orgunits.Delete(customerID, ouPath).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting organizational unit: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - OU contains users (move them first)\n")
		fmt.Fprintf(os.Stderr, "  - OU contains sub-OUs (delete them first)\n")
		fmt.Fprintf(os.Stderr, "  - OU path is incorrect\n")
		return err
	}

	fmt.Printf("Successfully deleted organizational unit: %s\n", ouPath)

	return nil
}
