# Audit Log Export Guide

## Overview

The audit log export feature allows you to retrieve and analyze Google Workspace activity logs for compliance, security monitoring, and troubleshooting. This guide covers how to export audit logs for different applications and use cases.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Application Types](#application-types)
- [Basic Usage](#basic-usage)
- [Filtering Options](#filtering-options)
- [Output Formats](#output-formats)
- [Common Use Cases](#common-use-cases)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required API Scopes

The audit export feature requires the Reports API scope, which is automatically included when you authenticate:
- `https://www.googleapis.com/auth/admin.reports.audit.readonly`

### Required Permissions

Your Google Workspace account must have one of the following roles:
- Super Admin
- Reports Reader (for read-only access to audit logs)

## Application Types

Google Workspace provides audit logs for the following applications:

| Application | Description | Common Events |
|------------|-------------|---------------|
| `admin` | Admin console activities | User creation, settings changes, role assignments |
| `login` | User authentication | Login, logout, 2FA events, suspicious activity |
| `drive` | Google Drive operations | File access, sharing, downloads, permissions |
| `calendar` | Calendar activities | Event creation, sharing, modifications |
| `groups` | Google Groups operations | Group creation, membership changes, settings |
| `mobile` | Mobile device management | Device enrollment, wipe, settings |
| `token` | OAuth token operations | Token grants, revocations |
| `groups_enterprise` | Groups for Enterprise | Enterprise group activities |
| `saml` | SAML authentication | SSO activities |
| `chrome` | Chrome browser management | Policy enforcement, extensions |
| `gcp` | Google Cloud Platform | GCP console activities |
| `chat` | Google Chat | Messages, rooms, settings |
| `meet` | Google Meet | Meeting joins, recordings, settings |

## Basic Usage

### Export Last 24 Hours

By default, the audit export retrieves the last 24 hours of activity:

```bash
# Export admin console activities (last 24h)
gac audit export --app admin

# Export login activities (last 24h)
gac audit export --app login

# Export drive activities (last 24h)
gac audit export --app drive
```

### Custom Time Range

Specify a custom time range using RFC3339 format:

```bash
# Export one week of admin activities
gac audit export --app admin \
  --start-time 2024-10-01T00:00:00Z \
  --end-time 2024-10-08T00:00:00Z

# Export specific date range for drive
gac audit export --app drive \
  --start-time 2024-09-01T00:00:00Z \
  --end-time 2024-09-30T23:59:59Z
```

**Time Format**: Use RFC3339 format with timezone (e.g., `2024-10-08T00:00:00Z` for UTC)

## Filtering Options

### Filter by User

Export activities for a specific user:

```bash
# All activities for user
gac audit export --app admin --user user@example.com

# Login activities for user
gac audit export --app login --user user@example.com
```

### Filter by Event Type

Filter by specific event names:

```bash
# Only user creation events
gac audit export --app admin --event-name USER_CREATED

# Multiple event types
gac audit export --app admin \
  --event-name USER_CREATED \
  --event-name USER_DELETED \
  --event-name ROLE_ASSIGNMENT
```

### Filter by IP Address

Filter activities from a specific IP address:

```bash
gac audit export --app login --actor-ip 192.168.1.100
```

### Limit Results

Limit the number of results returned:

```bash
# Get only the first 100 events
gac audit export --app admin --max-results 100
```

## Output Formats

### JSON Format (Default)

JSON output includes complete event details:

```bash
# Output to stdout
gac audit export --app admin

# Save to file
gac audit export --app admin --output-file admin-audit.json

# Explicit JSON format
gac audit export --app admin --output json --output-file audit.json
```

**JSON Structure:**
```json
[
  {
    "id": {
      "time": "2024-10-08T10:30:00.123Z",
      "uniqueQualifier": "12345",
      "applicationName": "admin"
    },
    "actor": {
      "email": "admin@example.com",
      "profileId": "123456789"
    },
    "events": [
      {
        "name": "USER_CREATED",
        "type": "USER_SETTINGS",
        "parameters": [...]
      }
    ],
    "ipAddress": "192.168.1.10"
  }
]
```

### CSV Format

CSV format provides a tabular view of key fields:

```bash
# Output CSV to stdout
gac audit export --app login --output csv

# Save CSV to file
gac audit export --app drive --output csv --output-file drive-audit.csv
```

**CSV Columns:**
- Timestamp
- Actor (user email)
- Event (event name)
- IP Address
- Application

## Common Use Cases

### 1. Security Monitoring

**Monitor Login Attempts:**
```bash
# Export all login activities
gac audit export --app login --output csv --output-file logins.csv

# Check for specific user's login history
gac audit export --app login --user user@example.com
```

**Detect Suspicious Activity:**
```bash
# Filter by suspicious IP
gac audit export --app login --actor-ip 203.0.113.42

# Check for password changes
gac audit export --app admin --event-name USER_PASSWORD_CHANGED
```

### 2. Compliance and Auditing

**User Activity Audit:**
```bash
# Export all activities for a user over the last month
gac audit export --app admin \
  --user user@example.com \
  --start-time 2024-09-01T00:00:00Z \
  --end-time 2024-09-30T23:59:59Z \
  --output csv --output-file user-audit.csv
```

**Admin Actions Audit:**
```bash
# Export all admin console changes
gac audit export --app admin \
  --start-time 2024-10-01T00:00:00Z \
  --output csv --output-file admin-changes.csv
```

### 3. File Access Monitoring

**Drive Activity Tracking:**
```bash
# Export drive activities for compliance
gac audit export --app drive \
  --start-time 2024-10-01T00:00:00Z \
  --end-time 2024-10-31T23:59:59Z \
  --output csv --output-file drive-access.csv

# Track file sharing events
gac audit export --app drive --event-name SHARE
```

### 4. Group Management Audit

**Track Group Changes:**
```bash
# Export group operations
gac audit export --app groups --output csv

# Filter by specific events
gac audit export --app groups \
  --event-name ADD_GROUP_MEMBER \
  --event-name REMOVE_GROUP_MEMBER
```

### 5. Automated Daily Reports

**Daily Audit Export Script:**
```bash
#!/bin/bash
# daily-audit-export.sh

DATE=$(date +%Y-%m-%d)
REPORT_DIR="/var/log/gac-audits"

mkdir -p "$REPORT_DIR"

# Export admin activities
gac audit export --app admin \
  --output csv \
  --output-file "$REPORT_DIR/admin-$DATE.csv"

# Export login activities
gac audit export --app login \
  --output csv \
  --output-file "$REPORT_DIR/login-$DATE.csv"

# Export drive activities
gac audit export --app drive \
  --output csv \
  --output-file "$REPORT_DIR/drive-$DATE.csv"

echo "Daily audit reports exported to $REPORT_DIR"
```

## Best Practices

### 1. Regular Exports

- **Schedule regular exports** using cron or scheduled tasks
- **Export incrementally** (e.g., daily or weekly) rather than large date ranges
- **Store exports securely** with appropriate access controls

### 2. Data Retention

- **Define retention policies** based on compliance requirements
- **Archive older logs** to long-term storage
- **Implement backup strategies** for audit data

### 3. Performance Optimization

- **Use time ranges** to limit result sets
- **Filter by user or event** when possible
- **Export during off-peak hours** for large queries
- **Use CSV format** for large exports (smaller file size)

### 4. Security Considerations

- **Limit access** to audit export functionality
- **Encrypt exported files** when storing or transferring
- **Monitor export activities** (admins exporting audit logs)
- **Review regularly** for anomalies or suspicious patterns

### 5. Analysis and Alerting

- **Import CSV to analysis tools** (Excel, database, SIEM)
- **Set up automated alerts** for critical events
- **Create dashboards** for key metrics
- **Document procedures** for investigating incidents

## Troubleshooting

### No Audit Logs Found

**Problem:** Command returns "No audit logs found"

**Solutions:**
1. **Check time range** - Ensure dates are within Google's retention period (typically 6 months)
2. **Verify application name** - Use exact app names (lowercase: `admin`, `login`, etc.)
3. **Check filters** - Overly restrictive filters may exclude all results
4. **Confirm activity** - Ensure there was actual activity in the time range

### Permission Errors

**Problem:** "Permission denied" or "403 Forbidden"

**Solutions:**
1. **Verify admin role** - Account must be Super Admin or have Reports Reader role
2. **Check API enablement** - Ensure Reports API is enabled in Google Cloud Console
3. **Re-authenticate** - Delete cached token and re-authenticate:
   ```bash
   rm ~/.credentials/gac.json
   gac audit export --app admin
   ```

### Large Result Sets

**Problem:** Export takes too long or times out

**Solutions:**
1. **Reduce time range** - Export smaller date ranges
2. **Use pagination** - Results are automatically paginated
3. **Add filters** - Filter by user, event, or IP to reduce results
4. **Use CSV format** - More efficient for large exports

### Invalid Time Format

**Problem:** "Invalid start-time format" error

**Solutions:**
1. **Use RFC3339 format** - `2024-10-08T00:00:00Z`
2. **Include timezone** - Use `Z` for UTC or specify timezone
3. **Check date validity** - Ensure dates are valid and properly formatted

Example:
```bash
# Correct
--start-time 2024-10-08T00:00:00Z

# Incorrect
--start-time 2024-10-08
--start-time "October 8, 2024"
```

## Event Reference

### Common Admin Events

- `USER_CREATED` - New user account created
- `USER_DELETED` - User account deleted
- `USER_SUSPENDED` - User account suspended
- `USER_UNSUSPENDED` - User account reactivated
- `USER_PASSWORD_CHANGED` - User password modified
- `ROLE_ASSIGNMENT` - Admin role assigned
- `ROLE_DELETION` - Admin role removed
- `ORG_UNIT_CREATED` - Organizational unit created
- `GROUP_CREATED` - Group created
- `GROUP_DELETED` - Group deleted

### Common Login Events

- `login` - Successful login
- `logout` - User logout
- `login_failure` - Failed login attempt
- `login_challenge` - 2-factor authentication challenge
- `suspicious_login` - Login flagged as suspicious
- `account_disabled_password_leak` - Account disabled due to leaked password

### Common Drive Events

- `create` - File or folder created
- `edit` - File edited
- `view` - File viewed
- `rename` - File or folder renamed
- `delete` - File or folder deleted
- `trash` - File moved to trash
- `untrash` - File restored from trash
- `upload` - File uploaded
- `download` - File downloaded
- `add_to_folder` - File added to folder
- `remove_from_folder` - File removed from folder
- `change_user_access` - Sharing permissions changed

## Related Resources

- [Command Reference](../reference/commands.md) - Complete command documentation
- [Authentication Guide](../authentication.md) - OAuth2 setup and configuration
- [Google Reports API Documentation](https://developers.google.com/admin-sdk/reports) - Official API docs
- [Audit Event Reference](https://developers.google.com/admin-sdk/reports/reference/rest/v1/activities/list) - Complete event reference

## Need Help?

- Check the [Troubleshooting Guide](../reference/troubleshooting.md)
- Report issues on [GitHub](https://github.com/acockrell/google-admin-client/issues)
- Review the [FAQ](../reference/troubleshooting.md#frequently-asked-questions)
