# Group Settings Management

Configure Google Workspace group settings including access control, posting permissions, moderation, and email preferences.

## Table of Contents

- [Overview](#overview)
- [View Group Settings](#view-group-settings)
- [Update Group Settings](#update-group-settings)
- [Configuration Scenarios](#configuration-scenarios)
- [Available Settings](#available-settings)
- [Best Practices](#best-practices)

## Overview

Group settings control how members interact with a Google Workspace group, including:

- **Access Control** - Who can join and view the group
- **Posting Permissions** - Who can post messages
- **Moderation** - Message approval and spam filtering
- **Email Settings** - Custom footers, reply-to behavior
- **Archive Settings** - Read-only groups for historical records
- **Member Management** - Invitation and approval workflows

## View Group Settings

Display current settings for a group in table or JSON format.

### Basic Usage

```bash
# View group settings in table format
gac group-settings list operations@example.com

# View group settings as JSON
gac group-settings list engineering --format json
```

### Flags

- `-f, --format` - Output format: `table` or `json` (default: table)

### Examples

```bash
# View all settings for a group
gac group-settings list team@example.com

# Export settings as JSON for backup
gac group-settings list announcements@example.com --format json > announcements-settings.json

# Quick check of support group settings
gac group-settings list support
```

## Update Group Settings

Modify group settings to control access, posting, and behavior.

### Basic Usage

```bash
# Allow anyone in domain to join
gac group-settings update operations@example.com \
  --who-can-join ALL_IN_DOMAIN_CAN_JOIN

# Configure posting permissions
gac group-settings update engineering@example.com \
  --who-can-post-message ALL_MEMBERS_CAN_POST \
  --allow-web-posting true

# Disable external members
gac group-settings update sales@example.com \
  --allow-external-members false

# Add custom footer to group emails
gac group-settings update support@example.com \
  --custom-footer-text "For assistance, contact support@example.com" \
  --include-custom-footer true

# Enable message moderation
gac group-settings update announcements@example.com \
  --message-moderation-level MODERATE_ALL_MESSAGES \
  --who-can-moderate-content ALL_MANAGERS_CAN_POST

# Configure reply-to settings
gac group-settings update team@example.com \
  --reply-to REPLY_TO_SENDER

# Make group archive-only (read-only)
gac group-settings update archive@example.com \
  --archive-only true
```

### Key Points

- **Partial updates** - Only specified settings are changed; others remain unchanged
- **Boolean values** - Use `true` or `false` for boolean settings
- **Domain appending** - Group name without `@` automatically gets domain appended
- **Validation** - Email addresses are validated before updates

## Configuration Scenarios

### Team Collaboration Group

Standard group for team collaboration with all members able to post.

```bash
gac group-settings update team@example.com \
  --who-can-join ALL_IN_DOMAIN_CAN_JOIN \
  --who-can-post-message ALL_MEMBERS_CAN_POST \
  --allow-web-posting true \
  --reply-to REPLY_TO_SENDER
```

### Moderated Announcements

Announcement-only group where only managers can post and all messages are moderated.

```bash
gac group-settings update announcements@example.com \
  --who-can-post-message ALL_MANAGERS_CAN_POST \
  --message-moderation-level MODERATE_ALL_MESSAGES \
  --who-can-view-group ALL_IN_DOMAIN_CAN_VIEW \
  --allow-external-members false
```

### External Partner Group

Group for collaborating with external partners with moderation for non-members.

```bash
gac group-settings update partners@example.com \
  --allow-external-members true \
  --who-can-join INVITED_CAN_JOIN \
  --who-can-post-message ALL_MEMBERS_CAN_POST \
  --message-moderation-level MODERATE_NON_MEMBERS
```

### Read-Only Archive

Historical group that's read-only, preserving past discussions.

```bash
gac group-settings update archive-2024@example.com \
  --archive-only true \
  --who-can-view-group ALL_IN_DOMAIN_CAN_VIEW \
  --who-can-post-message NONE_CAN_POST
```

### Support Group with Custom Footer

Customer support group with compliance footer and custom reply-to.

```bash
gac group-settings update support@example.com \
  --custom-footer-text "This email is for support purposes only. For urgent issues, call 1-800-SUPPORT." \
  --include-custom-footer true \
  --reply-to REPLY_TO_CUSTOM \
  --custom-reply-to support-team@example.com \
  --who-can-post-message ALL_IN_DOMAIN_CAN_POST
```

## Available Settings

### Access Control

Settings that control who can join and view the group:

- `--who-can-join` - Who can join the group
  - `CAN_REQUEST_TO_JOIN` - Users can request to join
  - `ALL_IN_DOMAIN_CAN_JOIN` - Anyone in domain can join
  - `ANYONE_CAN_JOIN` - Anyone on the internet can join
  - `INVITED_CAN_JOIN` - Only invited users can join

- `--who-can-view-group` - Who can view group messages
  - `ANYONE_CAN_VIEW` - Anyone can view
  - `ALL_IN_DOMAIN_CAN_VIEW` - Anyone in domain can view
  - `ALL_MEMBERS_CAN_VIEW` - Only members can view
  - `ALL_MANAGERS_CAN_VIEW` - Only managers can view

- `--who-can-view-membership` - Who can see the member list
  - Same values as `who-can-view-group`

- `--allow-external-members` - Allow external members (`true`/`false`)

### Posting Permissions

Settings that control who can post and how messages are moderated:

- `--who-can-post-message` - Who can post messages
  - `NONE_CAN_POST` - No one can post (read-only)
  - `ALL_MANAGERS_CAN_POST` - Only managers can post
  - `ALL_MEMBERS_CAN_POST` - All members can post
  - `ALL_IN_DOMAIN_CAN_POST` - Anyone in domain can post
  - `ANYONE_CAN_POST` - Anyone can post

- `--allow-web-posting` - Allow posting from web (`true`/`false`)

- `--message-moderation-level` - Message moderation level
  - `MODERATE_ALL_MESSAGES` - Moderate all messages
  - `MODERATE_NON_MEMBERS` - Moderate messages from non-members
  - `MODERATE_NEW_MEMBERS` - Moderate messages from new members
  - `MODERATE_NONE` - No moderation

- `--spam-moderation-level` - Spam filter sensitivity
  - `ALLOW` - Allow all messages
  - `MODERATE` - Moderate suspected spam
  - `SILENTLY_MODERATE` - Silently moderate spam
  - `REJECT` - Reject spam

### Email Settings

Settings for email customization and delivery:

- `--reply-to` - Reply-to behavior
  - `REPLY_TO_CUSTOM` - Use custom reply-to address
  - `REPLY_TO_SENDER` - Reply goes to sender
  - `REPLY_TO_LIST` - Reply goes to group
  - `REPLY_TO_OWNER` - Reply goes to group owner
  - `REPLY_TO_IGNORE` - Ignore reply-to

- `--custom-reply-to` - Custom reply-to email address

- `--custom-footer-text` - Custom footer text (max 1,000 characters)

- `--include-custom-footer` - Include custom footer (`true`/`false`)

- `--send-message-deny-notification` - Notify on message denial (`true`/`false`)

- `--include-in-global-address-list` - Show in GAL (`true`/`false`)

### Archive Settings

Settings for archiving and visibility:

- `--archive-only` - Make group read-only/archive-only (`true`/`false`)
  - When `true`, group rejects new messages

- `--show-in-group-directory` - Show in group directory (`true`/`false`)

### Member Management

Settings for managing group membership:

- `--who-can-leave-group` - Who can leave the group

- `--who-can-add` - Who can add members

- `--who-can-invite` - Who can invite members

- `--who-can-approve-members` - Who can approve membership requests

- `--allow-google-communication` - Allow Google communications (`true`/`false`)

- `--members-can-post-as-the-group` - Members can post as group (`true`/`false`)

### Moderation

Settings for moderation and content control:

- `--who-can-contact-owner` - Who can contact group owners

- `--who-can-moderate-members` - Who can moderate members

- `--who-can-moderate-content` - Who can moderate content/messages

- `--who-can-ban-users` - Who can ban users from group

## Best Practices

### Security

1. **Restrict external access** - Use `--allow-external-members false` for internal groups
2. **Use moderation** - Enable moderation for public-facing or announcement groups
3. **Control join permissions** - Set appropriate `--who-can-join` based on sensitivity
4. **Regular audits** - Periodically review group settings for compliance

### Organization

1. **Document purpose** - Use group descriptions to explain purpose and settings
2. **Consistent naming** - Use clear, descriptive group names
3. **Test before deploying** - Test settings with pilot group before rolling out
4. **Archive old groups** - Use `--archive-only` instead of deleting historical groups

### Communication

1. **Custom footers** - Add compliance/contact info with `--custom-footer-text`
2. **Reply-to control** - Set `--reply-to` to control conversation flow
3. **Spam protection** - Enable `--spam-moderation-level` for public groups
4. **Message moderation** - Use for announcement lists to prevent clutter

### Compliance

1. **External collaboration** - Moderate non-member messages with `MODERATE_NON_MEMBERS`
2. **Audit trails** - Document why groups have specific settings
3. **Access reviews** - Regularly review `--who-can-view-group` and `--allow-external-members`
4. **Footer requirements** - Add required compliance text to all groups

## Troubleshooting

### Common Issues

**Settings not updating**
- Verify you have Groups Settings Admin privileges
- Check group email is correct (try with full @domain.com)
- Ensure OAuth scope `https://www.googleapis.com/auth/apps.groups.settings` is enabled

**External members can't be enabled**
- Check your Google Workspace admin console allows external members
- Verify domain settings permit external sharing
- May require super admin to enable at organization level

**Custom footer not appearing**
- Ensure `--include-custom-footer true` is set
- Footer text must be under 1,000 characters
- May take a few minutes to propagate

**Archive-only not working**
- Verify you set `--archive-only true`
- Check that `--who-can-post-message NONE_CAN_POST` is also set
- May need to refresh group membership

### Getting Help

For more troubleshooting help, see:
- [Troubleshooting Guide](../reference/troubleshooting.md)
- [Command Reference](../reference/commands.md)
- [Group Management Guide](group-management.md)

## Examples

For complete working examples, see:
- [Group Settings Management Script](../../examples/group-settings-management.sh)
- [Examples README](../../examples/README.md#11-group-settings-management)

## Related Documentation

- [Group Management](group-management.md) - Creating and managing groups
- [User Management](user-management.md) - Managing group members
- [Command Reference](../reference/commands.md) - Complete command list
