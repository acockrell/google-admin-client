#!/bin/bash
#
# User Alias Management Example
#
# This script demonstrates how to manage email aliases for users,
# including adding, listing, and removing aliases.
#
# Usage:
#   ./alias-management.sh

set -euo pipefail

echo "========================================="
echo "User Alias Management Example"
echo "========================================="
echo ""
echo "This script demonstrates:"
echo "  1. Adding email aliases to users"
echo "  2. Listing user aliases"
echo "  3. Removing aliases"
echo ""

read -p "Press ENTER to start the demo (or Ctrl+C to cancel)..."
echo ""

# Configuration - Update these with your actual values
USER_EMAIL="${USER_EMAIL:-user@example.com}"
ALIAS1="${ALIAS1:-support@example.com}"
ALIAS2="${ALIAS2:-help@example.com}"
ALIAS3="${ALIAS3:-info@example.com}"

echo "Demo Configuration:"
echo "  User: $USER_EMAIL"
echo "  Aliases to add: $ALIAS1, $ALIAS2, $ALIAS3"
echo ""

read -p "Press ENTER to continue..."
echo ""

# Step 1: List current aliases
echo "[1/5] Listing current aliases for user..."
echo ""
gac alias list "$USER_EMAIL" || echo "User has no aliases yet"
echo ""

# Step 2: Add aliases
echo "[2/5] Adding aliases to user..."
echo ""

echo "Adding alias: $ALIAS1"
gac alias add "$USER_EMAIL" "$ALIAS1"
echo ""

echo "Adding alias: $ALIAS2"
gac alias add "$USER_EMAIL" "$ALIAS2"
echo ""

echo "Adding alias: $ALIAS3"
gac alias add "$USER_EMAIL" "$ALIAS3"
echo ""

echo "✓ Aliases added"
echo ""

# Step 3: List aliases after adding
echo "[3/5] Listing all aliases after adding..."
echo ""
gac alias list "$USER_EMAIL"
echo ""

# Step 4: Demonstrate removing an alias
echo "[4/5] Removing one alias..."
echo ""
echo "Removing alias: $ALIAS2"
gac alias remove "$USER_EMAIL" "$ALIAS2" --force
echo "✓ Alias removed"
echo ""

# Step 5: List final aliases
echo "[5/5] Final alias list..."
echo ""
gac alias list "$USER_EMAIL"
echo ""

# Cleanup instructions
echo "========================================="
echo "Cleanup Instructions"
echo "========================================="
echo ""
echo "To remove the remaining demo aliases, run:"
echo ""
echo "gac alias remove $USER_EMAIL $ALIAS1 --force"
echo "gac alias remove $USER_EMAIL $ALIAS3 --force"
echo ""
echo "========================================="
echo "Demo Complete!"
echo "========================================="
echo ""
echo "Key Takeaways:"
echo "  - Aliases allow users to receive mail at multiple addresses"
echo "  - Aliases must be in your organization's domains"
echo "  - Each alias can only be assigned to one user or group"
echo "  - Removing an alias stops mail delivery to that address"
echo ""
echo "Common Use Cases:"
echo "  - Department addresses (support@, sales@, info@)"
echo "  - Role-based addresses (admin@, webmaster@)"
echo "  - Alternative name formats (first.last@, firstlast@)"
echo "  - Legacy addresses when renaming users"
echo ""
echo "For more information:"
echo "  gac alias --help"
echo "  gac alias <command> --help"
