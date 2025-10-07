#!/bin/bash
#
# Onboarding Script for New Employees
#
# This script automates the process of creating and configuring
# a new user account in Google Workspace.
#
# Usage:
#   ./onboarding-example.sh jdoe@example.com "John" "Doe" \
#       "Engineering" "Software Engineer" "manager@example.com" \
#       "john.doe@personal.com"

set -euo pipefail

# Check for required arguments
if [ $# -lt 7 ]; then
    echo "Usage: $0 <email> <first-name> <last-name> <department> <title> <manager-email> <personal-email>"
    echo ""
    echo "Example:"
    echo "  $0 jdoe@example.com \"John\" \"Doe\" \"Engineering\" \\"
    echo "     \"Software Engineer\" \"manager@example.com\" \"john.doe@personal.com\""
    exit 1
fi

EMAIL="$1"
FIRST_NAME="$2"
LAST_NAME="$3"
DEPARTMENT="$4"
TITLE="$5"
MANAGER="$6"
PERSONAL_EMAIL="$7"

echo "========================================="
echo "Employee Onboarding Script"
echo "========================================="
echo "Email: $EMAIL"
echo "Name: $FIRST_NAME $LAST_NAME"
echo "Department: $DEPARTMENT"
echo "Title: $TITLE"
echo "Manager: $MANAGER"
echo "Personal Email: $PERSONAL_EMAIL"
echo "========================================="
echo ""

# Confirm before proceeding
read -p "Proceed with onboarding? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Onboarding cancelled"
    exit 0
fi

# Step 1: Create user
echo "[1/4] Creating user account..."
if gac user create \
    -f "$FIRST_NAME" \
    -l "$LAST_NAME" \
    -e "$PERSONAL_EMAIL" \
    -g all-staff \
    "$EMAIL"; then
    echo "✓ User created successfully"
else
    echo "✗ User creation failed"
    exit 1
fi
echo ""

# Step 2: Update user profile
echo "[2/4] Updating user profile..."
if gac user update \
    --dept "$DEPARTMENT" \
    --title "$TITLE" \
    --manager "$MANAGER" \
    --type staff \
    --id "$(uuidgen)" \
    "$EMAIL"; then
    echo "✓ Profile updated successfully"
else
    echo "✗ Profile update failed"
    exit 1
fi
echo ""

# Step 3: Add to department group
echo "[3/4] Adding to department group..."
DEPT_GROUP=$(echo "$DEPARTMENT" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
if gac user update -g "$DEPT_GROUP" "$EMAIL"; then
    echo "✓ Added to group: $DEPT_GROUP"
else
    echo "⚠ Warning: Failed to add to department group (may not exist)"
fi
echo ""

# Step 4: Summary
echo "[4/4] Onboarding Summary"
echo "========================================="
echo "User account created: $EMAIL"
echo "Temporary password sent to: $PERSONAL_EMAIL"
echo "User will be prompted to change password on first login"
echo ""
echo "Next steps:"
echo "  1. Notify $PERSONAL_EMAIL to check for welcome email"
echo "  2. Grant access to necessary applications"
echo "  3. Add to project-specific groups"
echo "  4. Schedule orientation meeting"
echo "========================================="
echo "✓ Onboarding complete for $EMAIL"
