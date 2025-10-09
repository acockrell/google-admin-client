package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	calResourceListType string
)

// calResourceListCmd represents the cal-resource list command
var calResourceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List calendar resources",
	Long: `List calendar resources in your Google Workspace domain.

Usage
-----

$ gac cal-resource list
$ gac cal-resource list --type room
$ gac cal-resource list --type equipment

Description
-----------

Lists calendar resources such as conference rooms, equipment, and other
bookable resources in your domain. You can filter by resource type.

The --type flag controls what to display:
  all       - Show all resources (default)
  room      - Show only rooms/conference rooms
  equipment - Show only equipment resources
  other     - Show other types of resources

`,
	RunE: calResourceListRunFunc,
}

func init() {
	calResourceCmd.AddCommand(calResourceListCmd)
	calResourceListCmd.Flags().StringVarP(&calResourceListType, "type", "t", "all", "resource type filter: all, room, equipment, or other")
}

func calResourceListRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	// Get the customer ID (my_customer for the current domain)
	customerID := "my_customer"

	// List all buildings first (needed for resource context)
	buildingsCall := client.Resources.Buildings.List(customerID)
	buildingsResult, err := buildingsCall.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not retrieve buildings: %v\n", err)
	}

	// Create a map of building IDs to names for display
	buildingMap := make(map[string]string)
	if buildingsResult != nil {
		for _, building := range buildingsResult.Buildings {
			buildingMap[building.BuildingId] = building.BuildingName
		}
	}

	// List calendar resources
	listCall := client.Resources.Calendars.List(customerID)

	result, err := listCall.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing calendar resources: %v\n", err)
		return err
	}

	if len(result.Items) == 0 {
		fmt.Println("No calendar resources found.")
		return nil
	}

	// Filter resources by type if specified
	var filteredResources []*admin.CalendarResource
	for _, resource := range result.Items {
		if calResourceListType == "all" {
			filteredResources = append(filteredResources, resource)
		} else if calResourceListType == "room" && resource.ResourceType == "ROOM" {
			filteredResources = append(filteredResources, resource)
		} else if calResourceListType == "equipment" && resource.ResourceType == "EQUIPMENT" {
			filteredResources = append(filteredResources, resource)
		} else if calResourceListType == "other" && resource.ResourceType != "ROOM" && resource.ResourceType != "EQUIPMENT" {
			filteredResources = append(filteredResources, resource)
		}
	}

	if len(filteredResources) == 0 {
		QuietPrintf("No calendar resources found matching type: %s\n", calResourceListType)
		return nil
	}

	// Display the calendar resources
	QuietPrintf("Found %d calendar resource(s):\n\n", len(filteredResources))

	// Convert to simplified list items
	type calResourceItem struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		ID          string `json:"id"`
		Type        string `json:"type"`
		Building    string `json:"building"`
		Floor       string `json:"floor"`
		Capacity    int64  `json:"capacity"`
		Description string `json:"description"`
	}

	var items []calResourceItem
	for _, resource := range filteredResources {
		building := ""
		if resource.BuildingId != "" {
			buildingName := buildingMap[resource.BuildingId]
			if buildingName != "" {
				building = fmt.Sprintf("%s (%s)", buildingName, resource.BuildingId)
			} else {
				building = resource.BuildingId
			}
		}

		item := calResourceItem{
			Name:        resource.ResourceName,
			Email:       resource.ResourceEmail,
			ID:          resource.ResourceId,
			Type:        resource.ResourceType,
			Building:    building,
			Floor:       resource.FloorName,
			Capacity:    resource.Capacity,
			Description: resource.ResourceDescription,
		}
		items = append(items, item)
	}

	headers := []string{"Name", "Email", "ID", "Type", "Building", "Floor", "Capacity", "Description"}

	// For JSON/YAML, output full resource data
	var outputData interface{}
	if outputFormat == OutputFormatJSON || outputFormat == OutputFormatYAML {
		outputData = filteredResources
	} else {
		outputData = items
	}

	if err := FormatOutput(outputData, headers); err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
