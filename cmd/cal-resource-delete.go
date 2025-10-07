package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	calResourceDeleteForce bool
)

// calResourceDeleteCmd represents the cal-resource delete command
var calResourceDeleteCmd = &cobra.Command{
	Use:   "delete <resource-id>",
	Short: "Delete a calendar resource",
	Long: `Delete a calendar resource from your Google Workspace domain.

Usage
-----

$ gac cal-resource delete old-projector
$ gac cal-resource delete conf-room-1 --force

Description
-----------

Deletes a calendar resource. This will remove the resource from the directory
and users will no longer be able to book it.

WARNING: This operation cannot be undone. Use with caution.

Any existing calendar events associated with this resource will remain, but
the resource will no longer be bookable for future events.

Examples:
  # Delete a resource with confirmation prompt
  gac cal-resource delete old-projector

  # Delete without confirmation prompt
  gac cal-resource delete conf-room-archived --force

Before deleting a resource:
  1. Check if there are any upcoming events booked for this resource
  2. Notify users who regularly use the resource
  3. Consider updating the resource instead of deleting if it's being replaced

`,
	Args: cobra.ExactArgs(1),
	RunE: calResourceDeleteRunFunc,
}

func init() {
	calResourceCmd.AddCommand(calResourceDeleteCmd)
	calResourceDeleteCmd.Flags().BoolVarP(&calResourceDeleteForce, "force", "f", false, "skip confirmation prompt")
}

func calResourceDeleteRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	resourceId := args[0]
	customerID := "my_customer"

	// Show warning and prompt for confirmation unless --force is used
	if !calResourceDeleteForce {
		// Get the resource details to show the user what they're deleting
		resource, err := client.Resources.Calendars.Get(customerID, resourceId).Do()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retrieving calendar resource: %v\n", err)
			return err
		}

		fmt.Printf("WARNING: You are about to delete calendar resource:\n")
		fmt.Printf("  Name: %s\n", resource.ResourceName)
		fmt.Printf("  Email: %s\n", resource.ResourceEmail)
		fmt.Printf("  Type: %s\n", resource.ResourceType)
		fmt.Printf("\nThis operation cannot be undone.\n\n")
		fmt.Printf("Type 'yes' to confirm deletion: ")

		var confirmation string
		_, err = fmt.Scanln(&confirmation)
		if err != nil {
			// If there's an error reading input (e.g., EOF), treat as cancellation
			fmt.Fprintf(os.Stderr, "\nDeletion cancelled.\n")
			return nil
		}

		if confirmation != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the calendar resource
	err = client.Resources.Calendars.Delete(customerID, resourceId).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting calendar resource: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nCommon reasons for failure:\n")
		fmt.Fprintf(os.Stderr, "  - Resource ID is incorrect\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions\n")
		fmt.Fprintf(os.Stderr, "  - Resource doesn't exist\n")
		return err
	}

	fmt.Printf("Successfully deleted calendar resource: %s\n", resourceId)

	return nil
}
