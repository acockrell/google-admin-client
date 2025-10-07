#!/bin/bash
#
# Organizational Unit Restructure Example
#
# This script demonstrates how to reorganize organizational units,
# including creating a new structure and moving users between OUs.
#
# Usage:
#   ./ou-restructure.sh

set -euo pipefail

echo "========================================="
echo "Organizational Unit Restructure Example"
echo "========================================="
echo ""
echo "This script will:"
echo "  1. Create a new OU structure"
echo "  2. List the created OUs"
echo "  3. Update OU properties"
echo "  4. Show how to clean up"
echo ""

read -p "Press ENTER to start the demo (or Ctrl+C to cancel)..."
echo ""

# Step 1: Create organizational structure
echo "[1/5] Creating organizational structure..."
echo ""

echo "Creating top-level OUs..."
gac ou create /Engineering --description "Engineering department"
gac ou create /Sales --description "Sales department"
gac ou create /Operations --description "Operations department"
echo ""

echo "Creating Engineering sub-OUs..."
gac ou create /Engineering/Backend --description "Backend engineering team"
gac ou create /Engineering/Frontend --description "Frontend engineering team"
gac ou create /Engineering/DevOps --description "DevOps and infrastructure team"
echo ""

echo "Creating Sales sub-OUs..."
gac ou create /Sales/Enterprise --description "Enterprise sales team"
gac ou create /Sales/SMB --description "Small and medium business sales"
echo ""

echo "✓ Organizational structure created"
echo ""

# Step 2: List the structure
echo "[2/5] Listing organizational structure..."
echo ""
gac ou list
echo ""

# Step 3: Update an OU
echo "[3/5] Updating organizational unit..."
echo ""
echo "Updating /Engineering description..."
gac ou update /Engineering --description "Engineering department - All technical teams"
echo "✓ OU updated"
echo ""

# Step 4: Demonstrate moving an OU
echo "[4/5] Demonstrating OU reorganization..."
echo ""
echo "Moving /Engineering/DevOps to /Operations/DevOps..."
gac ou update /Engineering/DevOps --parent /Operations --name "DevOps"
echo "✓ OU moved"
echo ""

echo "New structure:"
gac ou list /Operations
echo ""

# Step 5: Cleanup instructions
echo "[5/5] Cleanup Instructions"
echo "========================================="
echo ""
echo "To clean up the demo OUs, run these commands:"
echo ""
echo "# Delete sub-OUs first (must be empty of users)"
echo "gac ou delete /Engineering/Backend"
echo "gac ou delete /Engineering/Frontend"
echo "gac ou delete /Sales/Enterprise"
echo "gac ou delete /Sales/SMB"
echo "gac ou delete /Operations/DevOps"
echo ""
echo "# Then delete top-level OUs"
echo "gac ou delete /Engineering"
echo "gac ou delete /Sales"
echo "gac ou delete /Operations"
echo ""
echo "========================================="
echo "Demo Complete!"
echo "========================================="
echo ""
echo "Key Takeaways:"
echo "  - OUs must be created from top-level to nested"
echo "  - OUs must be empty before deletion"
echo "  - Moving OUs affects all users in that OU"
echo "  - Use descriptions to clarify OU purposes"
echo ""
echo "For more information:"
echo "  gac ou --help"
echo "  gac ou <command> --help"
