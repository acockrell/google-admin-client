#!/bin/bash
#
# Group Audit Script
#
# Audit Google Workspace groups to find groups containing inactive/former employees.
# This helps maintain security and clean up group memberships.
#
# Usage:
#   ./group-audit.sh [output-directory]

set -euo pipefail

# Output directory (default: current directory)
OUTPUT_DIR="${1:-.}"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Generate filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="${OUTPUT_DIR}/group_audit_${TIMESTAMP}.txt"

echo "========================================="
echo "Group Audit Report"
echo "========================================="
echo "Scanning for groups with inactive members..."
echo "Output file: $OUTPUT_FILE"
echo "========================================="
echo ""

# Create report header
cat > "$OUTPUT_FILE" << EOF
Google Workspace Group Audit Report
Generated: $(date)
================================================================================

Groups Containing Inactive/Former Employees
--------------------------------------------------------------------------------

EOF

# Run the audit
if gac group list --contains-former-employees >> "$OUTPUT_FILE" 2>&1; then
    echo "✓ Audit complete"
else
    echo "✗ Audit failed"
    exit 1
fi

# Add footer
cat >> "$OUTPUT_FILE" << EOF

--------------------------------------------------------------------------------
End of Report
================================================================================

Recommended Actions:
1. Review each group listed above
2. Remove inactive members from groups
3. Consider archiving unused groups
4. Update group ownership if needed
5. Document any groups that legitimately include former employees
   (e.g., alumni groups)

For more details on a specific group:
  gac group list <group-email> --get-members
EOF

# Count groups with issues
PROBLEM_GROUPS=$(grep -c "@" "$OUTPUT_FILE" || echo "0")

echo ""
echo "========================================="
echo "Audit Summary"
echo "========================================="
echo "Groups with inactive members: $PROBLEM_GROUPS"
echo "Full report: $OUTPUT_FILE"
echo "========================================="
echo ""

if [ "$PROBLEM_GROUPS" -gt 0 ]; then
    echo "⚠ Action required: Review and clean up group memberships"
    echo ""
    echo "View report:"
    echo "  cat '$OUTPUT_FILE'"
    echo ""
    echo "Get details for a specific group:"
    echo "  gac group list <group-email> --get-members"
else
    echo "✓ No groups with inactive members found"
fi
