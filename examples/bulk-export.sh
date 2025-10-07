#!/bin/bash
#
# Bulk User Export Script
#
# Export all Google Workspace users to CSV for backup, audit, or analysis.
#
# Usage:
#   ./bulk-export.sh [output-directory]

set -euo pipefail

# Output directory (default: current directory)
OUTPUT_DIR="${1:-.}"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/users_export_${TIMESTAMP}.csv"

echo "========================================="
echo "Bulk User Export"
echo "========================================="
echo "Output file: $OUTPUT_FILE"
echo "========================================="
echo ""

# Export users
echo "Exporting all users..."
if gac user list --csv --full > "$OUTPUT_FILE"; then
    echo "✓ Export complete"
else
    echo "✗ Export failed"
    exit 1
fi

# Calculate statistics
TOTAL_USERS=$(tail -n +2 "$OUTPUT_FILE" | wc -l | tr -d ' ')

echo ""
echo "========================================="
echo "Export Summary"
echo "========================================="
echo "Output file: $OUTPUT_FILE"
echo "Total users: $TOTAL_USERS"
echo "File size: $(ls -lh "$OUTPUT_FILE" | awk '{print $5}')"
echo "========================================="
echo ""
echo "You can now:"
echo "  - Open in spreadsheet: open '$OUTPUT_FILE'"
echo "  - View in terminal: cat '$OUTPUT_FILE' | column -t -s,"
echo "  - Search users: grep 'search-term' '$OUTPUT_FILE'"
