package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	ouDescription      string
	ouParent           string
	ouBlockInheritance bool
)

// ouCreateCmd represents the ou create command
var ouCreateCmd = &cobra.Command{
	Use:   "create <ou-path>",
	Short: "Create a new organizational unit",
	Long: `Create a new organizational unit in your Google Workspace domain.

Usage
-----

$ gac ou create /Engineering
$ gac ou create /Engineering/Backend --description "Backend engineering team"
$ gac ou create /Sales --parent / --description "Sales team"
$ gac ou create /IT --block-inheritance

Description
-----------

Creates a new organizational unit at the specified path. The OU path must
start with a forward slash (/).

The OU name is derived from the last segment of the path. For example,
/Engineering/Backend creates an OU named "Backend" under /Engineering.

If the parent path doesn't exist, you must create it first.

Examples:
  # Create a top-level OU
  gac ou create /Engineering --description "Engineering department"

  # Create a nested OU
  gac ou create /Engineering/Backend --description "Backend team"

  # Create with inheritance blocking
  gac ou create /Contractors --block-inheritance

`,
	Args: cobra.ExactArgs(1),
	RunE: ouCreateRunFunc,
}

func init() {
	ouCmd.AddCommand(ouCreateCmd)
	ouCreateCmd.Flags().StringVarP(&ouDescription, "description", "d", "", "organizational unit description")
	ouCreateCmd.Flags().StringVarP(&ouParent, "parent", "p", "", "parent OU path (auto-detected from path if not specified)")
	ouCreateCmd.Flags().BoolVarP(&ouBlockInheritance, "block-inheritance", "b", false, "block policy inheritance from parent")
}

func ouCreateRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	ouPath := args[0]

	// Validate OU path format
	if !strings.HasPrefix(ouPath, "/") {
		fmt.Fprintf(os.Stderr, "Error: OU path must start with '/'\n")
		return fmt.Errorf("invalid OU path format")
	}

	// Extract OU name from path
	pathSegments := strings.Split(strings.Trim(ouPath, "/"), "/")
	ouName := pathSegments[len(pathSegments)-1]

	// Determine parent path
	parentPath := ouParent
	if parentPath == "" {
		if len(pathSegments) == 1 {
			// Top-level OU, parent is root
			parentPath = "/"
		} else {
			// Nested OU, extract parent from path
			parentSegments := pathSegments[:len(pathSegments)-1]
			parentPath = "/" + strings.Join(parentSegments, "/")
		}
	}

	// Create the organizational unit
	ou := &admin.OrgUnit{
		Name:              ouName,
		ParentOrgUnitPath: parentPath,
		Description:       ouDescription,
		BlockInheritance:  ouBlockInheritance,
	}

	customerID := "my_customer"

	result, err := client.Orgunits.Insert(customerID, ou).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating organizational unit: %v\n", err)
		return err
	}

	fmt.Printf("Successfully created organizational unit:\n\n")
	fmt.Printf("  Name: %s\n", result.Name)
	fmt.Printf("  Path: %s\n", result.OrgUnitPath)
	if result.Description != "" {
		fmt.Printf("  Description: %s\n", result.Description)
	}
	fmt.Printf("  Parent: %s\n", result.ParentOrgUnitPath)
	fmt.Printf("  ID: %s\n", result.OrgUnitId)
	if result.BlockInheritance {
		fmt.Printf("  Block Inheritance: Yes\n")
	}

	return nil
}
