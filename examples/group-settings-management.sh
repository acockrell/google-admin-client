#!/bin/bash
#
# Group Settings Management Example
#
# This script demonstrates how to view and manage Google Workspace group settings,
# including access control, posting permissions, moderation, and email preferences.
#
# Usage:
#   ./group-settings-management.sh

set -euo pipefail

echo "========================================="
echo "Group Settings Management Example"
echo "========================================="
echo ""
echo "This script demonstrates:"
echo "  1. Viewing group settings"
echo "  2. Configuring access control (who can join/view)"
echo "  3. Setting posting permissions"
echo "  4. Enabling message moderation"
echo "  5. Adding custom email footers"
echo "  6. Managing archive settings"
echo ""

read -p "Press ENTER to start the demo (or Ctrl+C to cancel)..."
echo ""

# Configuration - Update these with your actual values
GROUP_EMAIL="${GROUP_EMAIL:-team@example.com}"
ANNOUNCEMENTS_GROUP="${ANNOUNCEMENTS_GROUP:-announcements@example.com}"
SUPPORT_GROUP="${SUPPORT_GROUP:-support@example.com}"

echo "Demo Configuration:"
echo "  Team Group: $GROUP_EMAIL"
echo "  Announcements Group: $ANNOUNCEMENTS_GROUP"
echo "  Support Group: $SUPPORT_GROUP"
echo ""

read -p "Press ENTER to continue..."
echo ""

# Step 1: View current group settings
echo "[1/7] Viewing current settings for $GROUP_EMAIL..."
echo ""
gac group-settings list "$GROUP_EMAIL"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 2: Configure access control
echo "[2/7] Configuring access control settings..."
echo ""
echo "Allowing anyone in domain to join the group..."
gac group-settings update "$GROUP_EMAIL" \
  --who-can-join ALL_IN_DOMAIN_CAN_JOIN \
  --allow-external-members false
echo ""
echo "✓ Access control updated"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 3: Configure posting permissions
echo "[3/7] Configuring posting permissions..."
echo ""
echo "Allowing all members to post messages..."
gac group-settings update "$GROUP_EMAIL" \
  --who-can-post-message ALL_MEMBERS_CAN_POST \
  --allow-web-posting true
echo ""
echo "✓ Posting permissions updated"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 4: Set up moderated announcements group
echo "[4/7] Setting up moderated announcements group..."
echo ""
echo "Configuring $ANNOUNCEMENTS_GROUP for moderated announcements..."
gac group-settings update "$ANNOUNCEMENTS_GROUP" \
  --who-can-post-message ALL_MANAGERS_CAN_POST \
  --message-moderation-level MODERATE_ALL_MESSAGES \
  --who-can-view-group ALL_IN_DOMAIN_CAN_VIEW
echo ""
echo "✓ Moderated announcements group configured"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 5: Add custom footer for support group
echo "[5/7] Adding custom footer to support group..."
echo ""
echo "Adding custom footer to $SUPPORT_GROUP..."
gac group-settings update "$SUPPORT_GROUP" \
  --custom-footer-text "For urgent issues, call 1-800-SUPPORT or visit https://support.example.com" \
  --include-custom-footer true
echo ""
echo "✓ Custom footer added"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 6: Configure reply-to settings
echo "[6/7] Configuring reply-to settings..."
echo ""
echo "Setting reply-to behavior for $GROUP_EMAIL..."
gac group-settings update "$GROUP_EMAIL" \
  --reply-to REPLY_TO_SENDER
echo ""
echo "✓ Reply-to settings updated"
echo ""

read -p "Press ENTER to continue to next step..."
echo ""

# Step 7: View updated settings
echo "[7/7] Viewing updated settings..."
echo ""
echo "Viewing settings in JSON format for $GROUP_EMAIL:"
echo ""
gac group-settings list "$GROUP_EMAIL" --format json
echo ""

echo "========================================="
echo "Demo Complete!"
echo "========================================="
echo ""
echo "Summary of what we demonstrated:"
echo ""
echo "✓ Viewed group settings in table and JSON formats"
echo "✓ Configured access control (who can join/view)"
echo "✓ Set posting permissions (who can post)"
echo "✓ Created a moderated announcements group"
echo "✓ Added custom email footer for support group"
echo "✓ Configured reply-to behavior"
echo ""
echo "Additional settings you can configure:"
echo ""
echo "  Access Control:"
echo "    --who-can-join                  (CAN_REQUEST_TO_JOIN, ALL_IN_DOMAIN_CAN_JOIN, etc.)"
echo "    --who-can-view-group            (ANYONE_CAN_VIEW, ALL_IN_DOMAIN_CAN_VIEW, etc.)"
echo "    --who-can-view-membership       (who can see member list)"
echo "    --allow-external-members        (true/false)"
echo ""
echo "  Posting Permissions:"
echo "    --who-can-post-message          (NONE_CAN_POST, ALL_MEMBERS_CAN_POST, etc.)"
echo "    --allow-web-posting             (true/false)"
echo "    --message-moderation-level      (MODERATE_ALL_MESSAGES, MODERATE_NONE, etc.)"
echo "    --spam-moderation-level         (spam filter settings)"
echo ""
echo "  Email Settings:"
echo "    --reply-to                      (REPLY_TO_SENDER, REPLY_TO_LIST, etc.)"
echo "    --custom-reply-to               (custom email address)"
echo "    --custom-footer-text            (custom footer message)"
echo "    --include-custom-footer         (true/false)"
echo "    --include-in-global-address-list (true/false)"
echo ""
echo "  Archive Settings:"
echo "    --archive-only                  (true/false - makes group read-only)"
echo "    --show-in-group-directory       (true/false)"
echo ""
echo "  Member Management:"
echo "    --who-can-leave-group           (who can leave)"
echo "    --who-can-add                   (who can add members)"
echo "    --who-can-invite                (who can invite)"
echo "    --who-can-approve-members       (who can approve join requests)"
echo "    --members-can-post-as-the-group (true/false)"
echo ""
echo "  Moderation:"
echo "    --who-can-contact-owner         (who can contact group owners)"
echo "    --who-can-moderate-members      (who can moderate members)"
echo "    --who-can-moderate-content      (who can moderate posts)"
echo "    --who-can-ban-users             (who can ban users)"
echo ""
echo "For more information, run:"
echo "  gac group-settings --help"
echo "  gac group-settings list --help"
echo "  gac group-settings update --help"
echo ""
