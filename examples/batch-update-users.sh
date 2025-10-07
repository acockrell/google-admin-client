#!/bin/bash
#
# Batch Update Users Script
#
# Update multiple users from a CSV file containing user information.
#
# CSV Format:
#   email,department,title,manager
#
# Usage:
#   ./batch-update-users.sh [csv-file]

set -euo pipefail

# Input file (default: users-to-update.csv)
INPUT_FILE="${1:-users-to-update.csv}"

# Check if file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: File not found: $INPUT_FILE"
    echo ""
    echo "Usage: $0 [csv-file]"
    echo ""
    echo "CSV format:"
    echo "  email,department,title,manager"
    echo "  jdoe@example.com,Engineering,Senior Engineer,manager@example.com"
    echo "  jsmith@example.com,Sales,Account Executive,sales-mgr@example.com"
    echo ""
    echo "Example CSV file: users-to-update.csv"
    exit 1
fi

# Count total users
TOTAL_USERS=$(($(wc -l < "$INPUT_FILE") - 1))

if [ "$TOTAL_USERS" -lt 1 ]; then
    echo "Error: CSV file contains no user records"
    exit 1
fi

echo "========================================="
echo "Batch User Update"
echo "========================================="
echo "Input file: $INPUT_FILE"
echo "Total users: $TOTAL_USERS"
echo "========================================="
echo ""

# Show preview
echo "Preview of first 3 users:"
head -n 4 "$INPUT_FILE" | column -t -s,
echo ""

# Confirm before proceeding
read -p "Proceed with batch update? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled"
    exit 0
fi
echo ""

# Initialize counters
SUCCESS_COUNT=0
FAIL_COUNT=0
CURRENT=0

# Process each user
tail -n +2 "$INPUT_FILE" | while IFS=, read -r EMAIL DEPT TITLE MANAGER; do
    CURRENT=$((CURRENT + 1))
    echo "[$CURRENT/$TOTAL_USERS] Processing: $EMAIL"
    echo "  Department: $DEPT"
    echo "  Title: $TITLE"
    echo "  Manager: $MANAGER"

    # Update the user
    if gac user update \
        --dept "$DEPT" \
        --title "$TITLE" \
        --manager "$MANAGER" \
        "$EMAIL" 2>&1 | grep -q "success\|updated"; then
        echo "  ✓ Updated successfully"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo "  ✗ Update failed"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    echo ""

    # Small delay to avoid rate limiting
    sleep 0.5
done

echo "========================================="
echo "Batch Update Complete"
echo "========================================="
echo "Total users processed: $TOTAL_USERS"
echo "Successful updates: $SUCCESS_COUNT"
echo "Failed updates: $FAIL_COUNT"
echo "========================================="

if [ "$FAIL_COUNT" -gt 0 ]; then
    echo "⚠ Some updates failed. Check error messages above."
    exit 1
else
    echo "✓ All users updated successfully"
fi
