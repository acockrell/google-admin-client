package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	ouUpdateName        string
	ouUpdateDescription string
	ouUpdateParent      string
	ouUpdateBlock       string
)

// ouUpdateCmd represents the ou update command
var ouUpdateCmd = &cobra.Command{
	Use:   "update <ou-path>",
	Short: "Update an organizational unit",
	Long: `Update an existing organizational unit in your Google Workspace domain.

Usage
-----

$ gac ou update /Engineering --description "Engineering department"
$ gac ou update /Engineering --name "Engineering-Team"
$ gac ou update /Engineering/Backend --parent /Engineering/Development
$ gac ou update /Contractors --block-inheritance true

Description
-----------

Updates properties of an existing organizational unit. You can update:
- Name (changes the OU name but not its path)
- Description
- Parent OU (moves the OU to a different location)
- Block inheritance setting (true/false)

Note: Moving an OU to a different parent will affect all users in that OU
and its children. Use with caution.

Examples:
  # Update description
  gac ou update /Engineering --description "Updated description"

  # Rename an OU
  gac ou update /Engineering --name "Engineering-Dept"

  # Move an OU to a different parent
  gac ou update /Engineering/QA --parent /Operations

  # Enable inheritance blocking
  gac ou update /Contractors --block-inheritance true

`,
	Args: cobra.ExactArgs(1),
	RunE: ouUpdateRunFunc,
}

func init() {
	ouCmd.AddCommand(ouUpdateCmd)
	ouUpdateCmd.Flags().StringVarP(&ouUpdateName, "name", "n", "", "new name for the organizational unit")
	ouUpdateCmd.Flags().StringVarP(&ouUpdateDescription, "description", "d", "", "new description")
	ouUpdateCmd.Flags().StringVarP(&ouUpdateParent, "parent", "p", "", "new parent OU path")
	ouUpdateCmd.Flags().StringVarP(&ouUpdateBlock, "block-inheritance", "b", "", "block policy inheritance (true/false)")
}

func ouUpdateRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	ouPath := args[0]
	customerID := "my_customer"

	// Create update request with only the fields that were specified
	ou := &admin.OrgUnit{}
	updated := false

	if ouUpdateName != "" {
		ou.Name = ouUpdateName
		updated = true
	}

	if ouUpdateDescription != "" {
		ou.Description = ouUpdateDescription
		updated = true
	}

	if ouUpdateParent != "" {
		ou.ParentOrgUnitPath = ouUpdateParent
		updated = true
	}

	if ouUpdateBlock != "" {
		if ouUpdateBlock == "true" {
			ou.BlockInheritance = true
			ou.ForceSendFields = append(ou.ForceSendFields, "BlockInheritance")
		} else if ouUpdateBlock == "false" {
			ou.BlockInheritance = false
			ou.ForceSendFields = append(ou.ForceSendFields, "BlockInheritance")
		} else {
			fmt.Fprintf(os.Stderr, "Error: --block-inheritance must be 'true' or 'false'\n")
			return fmt.Errorf("invalid block-inheritance value")
		}
		updated = true
	}

	if !updated {
		fmt.Fprintf(os.Stderr, "Error: No update fields specified. Use --name, --description, --parent, or --block-inheritance\n")
		return fmt.Errorf("no update fields specified")
	}

	// Extract OU path components for the API call
	// The API uses orgUnitPath in the form of OrgUnitId or the full path
	result, err := client.Orgunits.Update(customerID, ouPath, ou).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating organizational unit: %v\n", err)
		return err
	}

	fmt.Printf("Successfully updated organizational unit:\n\n")
	fmt.Printf("  Name: %s\n", result.Name)
	fmt.Printf("  Path: %s\n", result.OrgUnitPath)
	if result.Description != "" {
		fmt.Printf("  Description: %s\n", result.Description)
	}
	fmt.Printf("  Parent: %s\n", result.ParentOrgUnitPath)
	fmt.Printf("  ID: %s\n", result.OrgUnitId)
	if result.BlockInheritance {
		fmt.Printf("  Block Inheritance: Yes\n")
	} else {
		fmt.Printf("  Block Inheritance: No\n")
	}

	return nil
}
