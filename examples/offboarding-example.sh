#!/bin/bash
#
# Offboarding Script for Departing Employees
#
# This script automates the process of transferring data and
# disabling a user account in Google Workspace.
#
# Usage:
#   ./offboarding-example.sh olduser@example.com newowner@example.com

set -euo pipefail

# Check for required arguments
if [ $# -lt 2 ]; then
    echo "Usage: $0 <source-email> <destination-email>"
    echo ""
    echo "Example:"
    echo "  $0 olduser@example.com newowner@example.com"
    exit 1
fi

SOURCE_EMAIL="$1"
DESTINATION_EMAIL="$2"

echo "========================================="
echo "Employee Offboarding Script"
echo "========================================="
echo "Source User: $SOURCE_EMAIL"
echo "Destination User: $DESTINATION_EMAIL"
echo "========================================="
echo ""
echo "⚠ WARNING: This will:"
echo "  1. Transfer all documents and resources"
echo "  2. Disable the user account"
echo "  3. Clear personal information"
echo ""

# Confirm before proceeding
read -p "Are you sure you want to proceed? (type 'yes' to confirm): " -r
echo ""
if [[ ! $REPLY == "yes" ]]; then
    echo "Offboarding cancelled"
    exit 0
fi

# Step 1: Transfer data
echo "[1/4] Transferring data ownership..."
echo "This may take several minutes depending on the amount of data..."
if gac transfer --from "$SOURCE_EMAIL" --to "$DESTINATION_EMAIL"; then
    echo "✓ Transfer initiated successfully"
else
    echo "✗ Transfer failed"
    exit 1
fi
echo ""

# Step 2: Wait for transfer completion
echo "[2/4] Waiting for data transfer to complete..."
echo ""
echo "Please check the Google Admin Console for transfer status:"
echo "  Admin Console > Account > Data Migration > Transfer Tool for Users"
echo ""
read -p "Press ENTER when the transfer is complete and verified..." -r
echo ""

# Step 3: Disable account
echo "[3/4] Disabling user account..."
if gac user update --remove "$SOURCE_EMAIL"; then
    echo "✓ Account disabled successfully"
else
    echo "✗ Account disable failed"
    exit 1
fi
echo ""

# Step 4: Clear PII
echo "[4/4] Clearing personal information..."
if gac user update --clear-pii "$SOURCE_EMAIL"; then
    echo "✓ Personal information cleared"
else
    echo "⚠ Warning: Failed to clear personal information"
fi
echo ""

# Summary
echo "========================================="
echo "Offboarding Summary"
echo "========================================="
echo "✓ Data transferred from $SOURCE_EMAIL to $DESTINATION_EMAIL"
echo "✓ Account $SOURCE_EMAIL has been disabled"
echo "✓ Personal information cleared"
echo ""
echo "Additional manual steps:"
echo "  1. Remove from all groups (if needed)"
echo "  2. Revoke access to third-party applications"
echo "  3. Collect company equipment"
echo "  4. Archive email if required for compliance"
echo "  5. Document offboarding in HR system"
echo "========================================="
echo "Offboarding process completed for $SOURCE_EMAIL"
