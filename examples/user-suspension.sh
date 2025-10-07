#!/bin/bash
#
# User Account Suspension Example
#
# This script demonstrates how to suspend and unsuspend user accounts,
# which is useful for employee departures, security incidents, or leave.
#
# Usage:
#   ./user-suspension.sh

set -euo pipefail

echo "========================================="
echo "User Account Suspension Example"
echo "========================================="
echo ""
echo "This script demonstrates:"
echo "  1. Suspending a user account"
echo "  2. Viewing suspended user status"
echo "  3. Unsuspending a user account"
echo ""

read -p "Press ENTER to start the demo (or Ctrl+C to cancel)..."
echo ""

# Configuration - Update these with your actual values
USER_EMAIL="${USER_EMAIL:-user@example.com}"
SUSPENSION_REASON="${SUSPENSION_REASON:-Demo suspension for testing}"

echo "Demo Configuration:"
echo "  User: $USER_EMAIL"
echo "  Reason: $SUSPENSION_REASON"
echo ""

read -p "Press ENTER to continue..."
echo ""

# Step 1: Check current user status
echo "[1/5] Checking current user status..."
echo ""
gac user list "$USER_EMAIL" 2>/dev/null || echo "Note: This will show user details if they exist"
echo ""

# Step 2: Suspend the user
echo "[2/5] Suspending user account..."
echo ""
echo "Suspending user: $USER_EMAIL"
echo "Reason: $SUSPENSION_REASON"
gac user suspend "$USER_EMAIL" --reason "$SUSPENSION_REASON" --force
echo ""

echo "✓ User account suspended"
echo ""

# Step 3: Verify suspension status
echo "[3/5] Verifying suspension status..."
echo ""
gac user list "$USER_EMAIL" 2>/dev/null || echo "User is now suspended"
echo ""

# Step 4: Wait a moment
echo "[4/5] Waiting before unsuspending..."
echo ""
echo "In a real scenario, the account might remain suspended for:"
echo "  - Employee departure: Indefinitely"
echo "  - Security incident: Until resolved"
echo "  - Extended leave: Duration of leave"
echo "  - Policy violation: Until corrective action"
echo ""

read -p "Press ENTER to unsuspend the account..."
echo ""

# Step 5: Unsuspend the user
echo "[5/5] Unsuspending user account..."
echo ""
gac user unsuspend "$USER_EMAIL" --force
echo ""

echo "✓ User account restored"
echo ""

# Verify unsuspension
echo "Verifying unsuspension..."
gac user list "$USER_EMAIL" 2>/dev/null || echo "User is now active"
echo ""

# Summary
echo "========================================="
echo "Demo Complete!"
echo "========================================="
echo ""
echo "Key Takeaways:"
echo "  - Suspended users cannot sign in or access services"
echo "  - Emails to suspended accounts bounce"
echo "  - All data is preserved during suspension"
echo "  - Unsuspending immediately restores full access"
echo "  - Always document suspension reasons for audit purposes"
echo ""
echo "Common Workflows:"
echo ""
echo "1. Employee Departure:"
echo "   - Suspend account immediately upon departure"
echo "   - Transfer data ownership to manager"
echo "   - Review account access after retention period"
echo ""
echo "2. Security Incident:"
echo "   - Suspend account immediately"
echo "   - Investigate the incident"
echo "   - Unsuspend after confirming resolution"
echo "   - Review and update security policies"
echo ""
echo "3. Extended Leave:"
echo "   - Suspend account during leave period"
echo "   - Set calendar auto-responder"
echo "   - Unsuspend on return date"
echo ""
echo "Best Practices:"
echo "  - Always provide a clear suspension reason"
echo "  - Document suspension in your ticketing system"
echo "  - Review suspended accounts monthly"
echo "  - Have a clear unsuspension approval process"
echo "  - Communicate with affected users when appropriate"
echo ""
echo "For more information:"
echo "  gac user suspend --help"
echo "  gac user unsuspend --help"
