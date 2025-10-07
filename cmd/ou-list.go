package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	ouListType string
)

// ouListCmd represents the ou list command
var ouListCmd = &cobra.Command{
	Use:   "list [ou-path]",
	Short: "List organizational units",
	Long: `List organizational units in your Google Workspace domain.

Usage
-----

$ gac ou list
$ gac ou list /Engineering
$ gac ou list --type all

Description
-----------

Lists organizational units in a hierarchical structure. If an OU path is
provided, lists only that OU and its children. Otherwise, lists all OUs
in the domain starting from the root.

The --type flag controls what to display:
  all      - Show all OUs (default)
  children - Show only direct children of the specified OU

`,
	RunE: ouListRunFunc,
}

func init() {
	ouCmd.AddCommand(ouListCmd)
	ouListCmd.Flags().StringVarP(&ouListType, "type", "t", "all", "list type: all or children")
}

func ouListRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	// Determine the OU path to list
	ouPath := ""
	if len(args) > 0 {
		ouPath = args[0]
	}

	// Get the customer ID (my_customer for the current domain)
	customerID := "my_customer"

	// List organizational units
	listCall := client.Orgunits.List(customerID)

	if ouPath != "" {
		// List specific OU and optionally its children
		listCall = listCall.OrgUnitPath(ouPath)
	}

	if ouListType == "children" {
		listCall = listCall.Type("children")
	} else {
		listCall = listCall.Type("all")
	}

	result, err := listCall.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing organizational units: %v\n", err)
		return err
	}

	if result.OrganizationUnits == nil || len(result.OrganizationUnits) == 0 {
		fmt.Println("No organizational units found.")
		return nil
	}

	// Display the organizational units
	fmt.Printf("Found %d organizational unit(s):\n\n", len(result.OrganizationUnits))

	for _, ou := range result.OrganizationUnits {
		displayOU(ou)
		fmt.Println()
	}

	return nil
}

func displayOU(ou *admin.OrgUnit) {
	// Calculate indentation based on path depth
	depth := strings.Count(ou.OrgUnitPath, "/") - 1
	if depth < 0 {
		depth = 0
	}
	indent := strings.Repeat("  ", depth)

	fmt.Printf("%s%s\n", indent, ou.Name)
	fmt.Printf("%s  Path: %s\n", indent, ou.OrgUnitPath)

	if ou.Description != "" {
		fmt.Printf("%s  Description: %s\n", indent, ou.Description)
	}

	if ou.ParentOrgUnitPath != "" {
		fmt.Printf("%s  Parent: %s\n", indent, ou.ParentOrgUnitPath)
	}

	fmt.Printf("%s  ID: %s\n", indent, ou.OrgUnitId)

	if ou.BlockInheritance {
		fmt.Printf("%s  Block Inheritance: Yes\n", indent)
	}
}
