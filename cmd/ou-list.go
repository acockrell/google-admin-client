package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

// ouListItem represents a simplified OU for list output
type ouListItem struct {
	Name             string `json:"name"`
	Path             string `json:"path"`
	Description      string `json:"description"`
	ParentPath       string `json:"parentPath"`
	ID               string `json:"id"`
	BlockInheritance string `json:"blockInheritance"`
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

	if len(result.OrganizationUnits) == 0 {
		QuietPrintln("No organizational units found.")
		return nil
	}

	QuietPrintf("Found %d organizational unit(s):\n\n", len(result.OrganizationUnits))

	// Convert to simplified list items
	var items []ouListItem
	for _, ou := range result.OrganizationUnits {
		blockInheritance := "No"
		if ou.BlockInheritance {
			blockInheritance = "Yes"
		}
		item := ouListItem{
			Name:             ou.Name,
			Path:             ou.OrgUnitPath,
			Description:      ou.Description,
			ParentPath:       ou.ParentOrgUnitPath,
			ID:               ou.OrgUnitId,
			BlockInheritance: blockInheritance,
		}
		items = append(items, item)
	}

	headers := []string{"Name", "Path", "Description", "ParentPath", "ID", "BlockInheritance"}

	// For JSON/YAML, output full OU data
	var outputData interface{}
	if outputFormat == OutputFormatJSON || outputFormat == OutputFormatYAML {
		outputData = result.OrganizationUnits
	} else {
		outputData = items
	}

	if err := FormatOutput(outputData, headers); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
