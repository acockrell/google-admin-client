#!/bin/bash
#
# Calendar Resource Management Example
#
# This script demonstrates how to manage calendar resources like
# conference rooms and equipment in your Google Workspace domain.
#
# Usage:
#   ./resource-management.sh

set -euo pipefail

echo "========================================="
echo "Calendar Resource Management Example"
echo "========================================="
echo ""
echo "This script demonstrates:"
echo "  1. Creating calendar resources (rooms and equipment)"
echo "  2. Listing calendar resources"
echo "  3. Updating resource details"
echo "  4. Deleting resources"
echo ""

read -p "Press ENTER to start the demo (or Ctrl+C to cancel)..."
echo ""

# Configuration - Update these with your actual values
BUILDING_ID="${BUILDING_ID:-main-bldg}"
ROOM_ID="${ROOM_ID:-demo-conf-room-a}"
EQUIPMENT_ID="${EQUIPMENT_ID:-demo-projector-1}"

echo "Demo Configuration:"
echo "  Building ID: $BUILDING_ID"
echo "  Room Resource ID: $ROOM_ID"
echo "  Equipment Resource ID: $EQUIPMENT_ID"
echo ""

read -p "Press ENTER to continue..."
echo ""

# Step 1: List current resources
echo "[1/7] Listing current calendar resources..."
echo ""
gac cal-resource list || echo "No resources found yet"
echo ""

# Step 2: Create a conference room
echo "[2/7] Creating a conference room resource..."
echo ""
gac cal-resource create "$ROOM_ID" \
  --name "Demo Conference Room A" \
  --type room \
  --capacity 12 \
  --building-id "$BUILDING_ID" \
  --floor "3rd Floor" \
  --description "Demo room for testing resource management"
echo ""
echo "✓ Conference room created"
echo ""

# Step 3: Create equipment
echo "[3/7] Creating an equipment resource..."
echo ""
gac cal-resource create "$EQUIPMENT_ID" \
  --name "Demo HD Projector" \
  --type equipment \
  --category "AV Equipment" \
  --description "High-definition projector for presentations"
echo ""
echo "✓ Equipment resource created"
echo ""

# Step 4: List all resources
echo "[4/7] Listing all calendar resources..."
echo ""
gac cal-resource list
echo ""

# Step 5: List only rooms
echo "[5/7] Listing only conference rooms..."
echo ""
gac cal-resource list --type room
echo ""

# Step 6: Update a resource
echo "[6/7] Updating conference room capacity..."
echo ""
gac cal-resource update "$ROOM_ID" \
  --capacity 15 \
  --user-description "Updated demo room with increased capacity"
echo ""
echo "✓ Resource updated"
echo ""

# Step 7: Demonstrate deleting a resource
echo "[7/7] Deleting demo resources..."
echo ""

echo "Deleting equipment resource..."
gac cal-resource delete "$EQUIPMENT_ID" --force
echo "✓ Equipment deleted"
echo ""

echo "Deleting conference room..."
gac cal-resource delete "$ROOM_ID" --force
echo "✓ Conference room deleted"
echo ""

# Cleanup instructions
echo "========================================="
echo "Demo Complete!"
echo "========================================="
echo ""
echo "Key Takeaways:"
echo "  - Calendar resources are bookable items like rooms and equipment"
echo "  - Resources can be assigned to buildings and floors for organization"
echo "  - Each resource gets a unique email address for booking"
echo "  - Resources can be filtered by type (room, equipment, other)"
echo "  - Capacity helps prevent overbooking of shared spaces"
echo ""
echo "Common Use Cases:"
echo "  - Conference rooms and meeting spaces"
echo "  - Equipment (projectors, cameras, laptops)"
echo "  - Company vehicles"
echo "  - Shared workspaces and hot desks"
echo ""
echo "Best Practices:"
echo "  - Use consistent naming conventions (e.g., bldg-floor-room)"
echo "  - Set accurate capacity to prevent overbooking"
echo "  - Organize resources by building and floor"
echo "  - Use meaningful categories for equipment"
echo "  - Update resource descriptions when changes are made"
echo ""
echo "For more information:"
echo "  gac cal-resource --help"
echo "  gac cal-resource <command> --help"
echo ""
