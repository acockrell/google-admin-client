package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

var (
	calResourceName        string
	calResourceType        string
	calResourceDescription string
	calResourceCategory    string
	calResourceBuildingId  string
	calResourceFloorName   string
	calResourceFloorSection string
	calResourceCapacity    int64
	calResourceUserVisibleDesc string
)

// calResourceCreateCmd represents the cal-resource create command
var calResourceCreateCmd = &cobra.Command{
	Use:   "create <resource-id>",
	Short: "Create a new calendar resource",
	Long: `Create a new calendar resource in your Google Workspace domain.

Usage
-----

$ gac cal-resource create conf-room-1 --name "Conference Room 1" --type room --capacity 10
$ gac cal-resource create projector-1 --name "HD Projector" --type equipment --category "AV Equipment"
$ gac cal-resource create room-5a --name "Meeting Room 5A" --type room --building-id bld-001 --floor 5 --capacity 6

Description
-----------

Creates a new calendar resource such as a conference room, equipment, or
other bookable resource. The resource ID must be unique within your domain.

Resource types:
  room      - Conference rooms and meeting spaces
  equipment - Projectors, laptops, cameras, etc.
  other     - Other types of bookable resources

Examples:
  # Create a conference room
  gac cal-resource create conf-a --name "Conference Room A" --type room --capacity 12 --building-id main-bldg

  # Create equipment
  gac cal-resource create proj-hd-1 --name "HD Projector" --type equipment --category "AV Equipment"

  # Create a resource with location details
  gac cal-resource create room-exec --name "Executive Boardroom" --type room --capacity 20

`,
	Args: cobra.ExactArgs(1),
	RunE: calResourceCreateRunFunc,
}

func init() {
	calResourceCmd.AddCommand(calResourceCreateCmd)
	calResourceCreateCmd.Flags().StringVarP(&calResourceName, "name", "n", "", "resource name (required)")
	calResourceCreateCmd.Flags().StringVarP(&calResourceType, "type", "t", "room", "resource type: room, equipment, or other")
	calResourceCreateCmd.Flags().StringVarP(&calResourceDescription, "description", "d", "", "resource description")
	calResourceCreateCmd.Flags().StringVarP(&calResourceCategory, "category", "c", "", "resource category")
	calResourceCreateCmd.Flags().StringVarP(&calResourceBuildingId, "building-id", "b", "", "building ID where resource is located")
	calResourceCreateCmd.Flags().StringVarP(&calResourceFloorName, "floor", "f", "", "floor name/number")
	calResourceCreateCmd.Flags().StringVarP(&calResourceFloorSection, "section", "s", "", "floor section")
	calResourceCreateCmd.Flags().Int64Var(&calResourceCapacity, "capacity", 0, "resource capacity (for rooms)")
	calResourceCreateCmd.Flags().StringVar(&calResourceUserVisibleDesc, "user-description", "", "user-visible description")

	if err := calResourceCreateCmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}
}

func calResourceCreateRunFunc(cmd *cobra.Command, args []string) error {
	client, err := newAdminClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	resourceId := args[0]

	// Validate resource type
	validTypes := map[string]string{
		"room":      "ROOM",
		"equipment": "EQUIPMENT",
		"other":     "OTHER",
	}
	apiResourceType, ok := validTypes[calResourceType]
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: Invalid resource type '%s'. Must be: room, equipment, or other\n", calResourceType)
		return fmt.Errorf("invalid resource type")
	}

	// Create the calendar resource
	resource := &admin.CalendarResource{
		ResourceId:              resourceId,
		ResourceName:            calResourceName,
		ResourceType:            apiResourceType,
		ResourceDescription:     calResourceDescription,
		ResourceCategory:        calResourceCategory,
		BuildingId:              calResourceBuildingId,
		FloorName:               calResourceFloorName,
		FloorSection:            calResourceFloorSection,
		UserVisibleDescription:  calResourceUserVisibleDesc,
	}

	// Only set capacity if it's greater than 0
	if calResourceCapacity > 0 {
		resource.Capacity = calResourceCapacity
	}

	// Note: FeatureInstances is managed separately through the Features API
	// and cannot be set directly during resource creation
	// Features must be created first, then associated with resources

	customerID := "my_customer"

	result, err := client.Resources.Calendars.Insert(customerID, resource).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating calendar resource: %v\n", err)
		return err
	}

	fmt.Printf("Successfully created calendar resource:\n\n")
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

	if result.Capacity > 0 {
		fmt.Printf("  Capacity: %d\n", result.Capacity)
	}

	// Note: Features are displayed separately via the Features API
	// and are not included in the basic resource creation response

	return nil
}
