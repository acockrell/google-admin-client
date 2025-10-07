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
		fmt.Printf("No calendar resources found matching type: %s\n", calResourceListType)
		return nil
	}

	// Display the calendar resources
	fmt.Printf("Found %d calendar resource(s):\n\n", len(filteredResources))

	for _, resource := range filteredResources {
		displayResource(resource, buildingMap)
		fmt.Println()
	}

	return nil
}

func displayResource(resource *admin.CalendarResource, buildingMap map[string]string) {
	fmt.Printf("%s\n", resource.ResourceName)
	fmt.Printf("  Email: %s\n", resource.ResourceEmail)
	fmt.Printf("  ID: %s\n", resource.ResourceId)
	fmt.Printf("  Type: %s\n", resource.ResourceType)

	if resource.ResourceDescription != "" {
		fmt.Printf("  Description: %s\n", resource.ResourceDescription)
	}

	if resource.ResourceCategory != "" {
		fmt.Printf("  Category: %s\n", resource.ResourceCategory)
	}

	if resource.BuildingId != "" {
		buildingName := buildingMap[resource.BuildingId]
		if buildingName != "" {
			fmt.Printf("  Building: %s (%s)\n", buildingName, resource.BuildingId)
		} else {
			fmt.Printf("  Building ID: %s\n", resource.BuildingId)
		}
	}

	if resource.FloorName != "" {
		fmt.Printf("  Floor: %s\n", resource.FloorName)
	}

	if resource.FloorSection != "" {
		fmt.Printf("  Floor Section: %s\n", resource.FloorSection)
	}

	if resource.Capacity > 0 {
		fmt.Printf("  Capacity: %d\n", resource.Capacity)
	}

	if resource.UserVisibleDescription != "" {
		fmt.Printf("  User Visible Description: %s\n", resource.UserVisibleDescription)
	}

	// Handle FeatureInstances (interface{} type requires type assertion)
	if resource.FeatureInstances != nil {
		if features, ok := resource.FeatureInstances.([]interface{}); ok && len(features) > 0 {
			fmt.Printf("  Features: ")
			for i, f := range features {
				if i > 0 {
					fmt.Printf(", ")
				}
				if featureMap, ok := f.(map[string]interface{}); ok {
					if feature, ok := featureMap["feature"].(map[string]interface{}); ok {
						if name, ok := feature["name"].(string); ok {
							fmt.Printf("%s", name)
						}
					}
				}
			}
			fmt.Println()
		}
	}
}
