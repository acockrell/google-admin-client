package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	updateCalResourceName        string
	updateCalResourceDescription string
	updateCalResourceCategory    string
	updateCalResourceBuildingId  string
	updateCalResourceFloorName   string
	updateCalResourceFloorSection string
	updateCalResourceCapacity    int64
	updateCalResourceUserVisibleDesc string
)

// calResourceUpdateCmd represents the cal-resource update command
var calResourceUpdateCmd = &cobra.Command{
	Use:   "update <resource-id>",
	Short: "Update an existing calendar resource",
	Long: `Update an existing calendar resource in your Google Workspace domain.

Usage
-----

$ gac cal-resource update conf-room-1 --name "Conference Room 1 - Updated"
$ gac cal-resource update projector-1 --capacity 15
$ gac cal-resource update room-5a --building-id new-bldg --floor 3

Description
-----------

Updates properties of an existing calendar resource. You can update the name,
description, location, capacity, and other attributes.

Only the flags you specify will be updated. Other properties will remain
unchanged.

Examples:
  # Update room capacity
  gac cal-resource update conf-a --capacity 15

  # Update room name and description
  gac cal-resource update conf-a --name "Conference Room A (Renovated)" --description "Newly renovated conference room"

  # Update building and floor
  gac cal-resource update room-123 --building-id main-building --floor "3rd Floor"

  # Update user-visible description
  gac cal-resource update exec-room --user-description "Large executive meeting room with video conferencing"

`,
	Args: cobra.ExactArgs(1),
	RunE: calResourceUpdateRunFunc,
}

func init() {
	calResourceCmd.AddCommand(calResourceUpdateCmd)
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceName, "name", "n", "", "resource name")
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceDescription, "description", "d", "", "resource description")
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceCategory, "category", "c", "", "resource category")
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceBuildingId, "building-id", "b", "", "building ID where resource is located")
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceFloorName, "floor", "f", "", "floor name/number")
	calResourceUpdateCmd.Flags().StringVarP(&updateCalResourceFloorSection, "section", "s", "", "floor section")
	calResourceUpdateCmd.Flags().Int64Var(&updateCalResourceCapacity, "capacity", -1, "resource capacity (for rooms)")
	calResourceUpdateCmd.Flags().StringVar(&updateCalResourceUserVisibleDesc, "user-description", "", "user-visible description")
}

func calResourceUpdateRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	resourceId := args[0]
	customerID := "my_customer"

	// First, get the existing resource to preserve unchanged fields
	existing, err := client.Resources.Calendars.Get(customerID, resourceId).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving calendar resource: %v\n", err)
		return err
	}

	// Start with the existing resource
	// Note: FeatureInstances is omitted as it's managed separately
	resource := &admin.CalendarResource{
		ResourceId:              existing.ResourceId,
		ResourceName:            existing.ResourceName,
		ResourceType:            existing.ResourceType,
		ResourceDescription:     existing.ResourceDescription,
		ResourceCategory:        existing.ResourceCategory,
		BuildingId:              existing.BuildingId,
		FloorName:               existing.FloorName,
		FloorSection:            existing.FloorSection,
		UserVisibleDescription:  existing.UserVisibleDescription,
		Capacity:                existing.Capacity,
	}

	// Update only the fields that were specified
	if cmd.Flags().Changed("name") {
		resource.ResourceName = updateCalResourceName
	}
	if cmd.Flags().Changed("description") {
		resource.ResourceDescription = updateCalResourceDescription
	}
	if cmd.Flags().Changed("category") {
		resource.ResourceCategory = updateCalResourceCategory
	}
	if cmd.Flags().Changed("building-id") {
		resource.BuildingId = updateCalResourceBuildingId
	}
	if cmd.Flags().Changed("floor") {
		resource.FloorName = updateCalResourceFloorName
	}
	if cmd.Flags().Changed("section") {
		resource.FloorSection = updateCalResourceFloorSection
	}
	if cmd.Flags().Changed("capacity") {
		resource.Capacity = updateCalResourceCapacity
	}
	if cmd.Flags().Changed("user-description") {
		resource.UserVisibleDescription = updateCalResourceUserVisibleDesc
	}

	result, err := client.Resources.Calendars.Update(customerID, resourceId, resource).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating calendar resource: %v\n", err)
		return err
	}

	fmt.Printf("Successfully updated calendar resource:\n\n")
	fmt.Printf("  Name: %s\n", result.ResourceName)
	fmt.Printf("  Email: %s\n", result.ResourceEmail)
	fmt.Printf("  ID: %s\n", result.ResourceId)
	fmt.Printf("  Type: %s\n", result.ResourceType)

	if result.ResourceDescription != "" {
		fmt.Printf("  Description: %s\n", result.ResourceDescription)
	}

	if result.ResourceCategory != "" {
		fmt.Printf("  Category: %s\n", result.ResourceCategory)
	}

	if result.BuildingId != "" {
		fmt.Printf("  Building ID: %s\n", result.BuildingId)
	}

	if result.FloorName != "" {
		fmt.Printf("  Floor: %s\n", result.FloorName)
	}

	if result.FloorSection != "" {
		fmt.Printf("  Floor Section: %s\n", result.FloorSection)
	}

	if result.Capacity > 0 {
		fmt.Printf("  Capacity: %d\n", result.Capacity)
	}

	if result.UserVisibleDescription != "" {
		fmt.Printf("  User Visible Description: %s\n", result.UserVisibleDescription)
	}

	// Handle FeatureInstances (interface{} type requires type assertion)
	if result.FeatureInstances != nil {
		if features, ok := result.FeatureInstances.([]interface{}); ok && len(features) > 0 {
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

	return nil
}
