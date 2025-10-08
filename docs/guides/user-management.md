# User Management

Comprehensive guide for managing users in Google Workspace with `gac`.

## Table of Contents

- [Create a User](#create-a-user)
- [List Users](#list-users)
- [Update a User](#update-a-user)
- [Suspend User Account](#suspend-user-account)
- [Unsuspend User Account](#unsuspend-user-account)
- [Common Workflows](#common-workflows)

## Create a User

Create new users in your Google Workspace domain.

### Basic Usage

```bash
# Interactive creation (prompts for details)
gac user create newuser@example.com

# With optional details
gac user create \
  -f John \
  -l Doe \
  -e personal@email.com \
  -g engineering \
  -g all-staff \
  newuser@example.com
```

### Flags

- `-f, --first-name` - First name
- `-l, --last-name` - Last name
- `-e, --email` - Personal email address
- `-g, --groups` - Groups to add user to (can be repeated)

### Examples

```bash
# Create user with full details
gac user create \
  --first-name Sarah \
  --last-name Johnson \
  --email sarah.personal@gmail.com \
  --groups engineering \
  --groups team-platform \
  sjohnson@example.com

# Create contractor account
gac user create \
  --first-name Mike \
  --last-name Smith \
  msmith-contractor@example.com
```

## List Users

List and export user information from your domain.

### Basic Usage

```bash
# List all active users
gac user list

# List a specific user
gac user list jdoe@example.com

# List only disabled accounts
gac user list --disabled-only

# Export to CSV
gac user list --csv > users.csv

# Full export with all fields
gac user list --full
```

### Flags

- `-c, --csv` - Export as CSV
- `-d, --disabled-only` - Show only disabled accounts
- `-f, --full` - Include all user fields

### Examples

```bash
# Export all users to CSV
gac user list --csv > all-users.csv

# List only suspended accounts
gac user list --disabled-only

# Get detailed info for specific user
gac user list --full admin@example.com
```

## Update a User

Update user information, department, groups, and more.

### Basic Usage

```bash
# Update department and title
gac user update --dept Engineering --title "Senior Engineer" jdoe@example.com

# Add user to groups
gac user update -g developers -g team-leads jdoe@example.com

# Update phone numbers
gac user update --phone "mobile:555-555-5555" jdoe@example.com
gac user update --phone "mobile:555-555-5555; work:555-123-4567,555" jdoe@example.com

# Update address
gac user update --address "Columbus, OH" jdoe@example.com

# Set manager
gac user update --manager manager@example.com jdoe@example.com

# Update employee type
gac user update --type staff jdoe@example.com

# Set employee ID (UUID)
gac user update --id $(uuidgen) jdoe@example.com
gac user update --id $(uuidgen) --force jdoe@example.com

# Update organizational unit
gac user update --ou /Engineering/Backend jdoe@example.com

# Update custom fields
gac user update --github-profile doeMaker jdoe@example.com
gac user update --amazon-username john.doe jdoe@example.com
gac user update --vpn-role developer jdoe@example.com

# Disable user account
gac user update --remove jdoe@example.com

# Clear personal information
gac user update --clear-pii jdoe@example.com
```

### Flags

- `-t, --title` - Job title
- `-d, --dept` - Department
- `-e, --type` - Employee type (staff or contractor)
- `-g, --group` - Groups to add user to (can be repeated)
- `-m, --manager` - Manager's email address
- `-p, --phone` - Phone number(s) in format "type:number" or "type:number; type:number,ext"
- `-a, --address` - Work address
- `-o, --ou` - Organizational unit path
- `-i, --id` - Employee UUID
- `-f, --force` - Overwrite existing values (e.g., employee ID)
- `--github-profile` - GitHub username
- `--amazon-username` - Amazon username
- `--vpn-role` - VPN access role
- `-r, --remove` - Disable user account
- `--clear-pii` - Clear personal information

### Examples

```bash
# Promote user to senior role
gac user update \
  --title "Senior Software Engineer" \
  --dept Engineering \
  --group tech-leads \
  jdoe@example.com

# Update contact information
gac user update \
  --phone "mobile:555-1234; work:555-5678,101" \
  --address "123 Main St, Columbus, OH 43215" \
  jdoe@example.com

# Move user to different OU
gac user update --ou /Sales/West-Coast jdoe@example.com

# Update all custom fields
gac user update \
  --github-profile johndoe \
  --amazon-username john.doe \
  --vpn-role engineer \
  jdoe@example.com
```

## Suspend User Account

Suspend user accounts to prevent access while preserving data.

### Basic Usage

```bash
# Suspend a user with confirmation
gac user suspend user@example.com

# Suspend with a reason
gac user suspend user@example.com --reason "Left company"

# Suspend without confirmation
gac user suspend user@example.com --force
```

### What Happens When You Suspend

Suspending a user account prevents the user from:
- Signing in to their account
- Accessing any Google Workspace services (Gmail, Drive, Calendar, etc.)
- Receiving new emails (emails will bounce)

**Important:** The account data is preserved and can be restored by unsuspending the account.

### Flags

- `-r, --reason` - Reason for suspension (optional)
- `-f, --force` - Skip confirmation prompt

### Common Use Cases

- **Employee termination or departure** - Immediately disable access when someone leaves
- **Policy violations or security incidents** - Temporarily suspend while investigating
- **Account compromise or suspicious activity** - Protect account from unauthorized access
- **Extended leave or sabbatical** - Suspend during long absences

### Examples

```bash
# Suspend departing employee
gac user suspend \
  --reason "Terminated - Last day 2025-10-15" \
  --force \
  former-employee@example.com

# Suspend for security review
gac user suspend \
  --reason "Security review - suspicious login activity" \
  user@example.com

# Quick suspension (with prompt)
gac user suspend user@example.com
```

## Unsuspend User Account

Restore access to previously suspended accounts.

### Basic Usage

```bash
# Unsuspend a user with confirmation
gac user unsuspend user@example.com

# Unsuspend without confirmation
gac user unsuspend user@example.com --force
```

### What Happens When You Unsuspend

Unsuspending a user account restores the user's ability to:
- Sign in to their account
- Access all Google Workspace services (Gmail, Drive, Calendar, etc.)
- Send and receive emails
- Access their data and documents

**Important:** All account data and settings are preserved during suspension and will be available after unsuspending.

### Flags

- `-f, --force` - Skip confirmation prompt

### Common Use Cases

- **Restoring access after employee returns from leave** - Re-enable after sabbatical or extended leave
- **Correcting accidental suspensions** - Quick restoration if suspended in error
- **Restoring accounts after security incidents are resolved** - Re-enable after investigation complete
- **Re-enabling accounts after policy violations are addressed** - Restore access after remediation

### Examples

```bash
# Restore access for returning employee
gac user unsuspend user@example.com --force

# Restore after security review
gac user unsuspend \
  user@example.com

# Quick restore (with prompt)
gac user unsuspend user@example.com
```

## Common Workflows

### Employee Onboarding

Complete workflow for onboarding a new employee:

```bash
# 1. Create user account
gac user create \
  --first-name Jane \
  --last-name Smith \
  --email jane.personal@gmail.com \
  --groups all-staff \
  --groups engineering \
  jsmith@example.com

# 2. Set up user details
gac user update \
  --title "Software Engineer" \
  --dept Engineering \
  --ou /Engineering/Backend \
  --manager tech-lead@example.com \
  --phone "mobile:555-1234" \
  --address "San Francisco, CA" \
  jsmith@example.com

# 3. Add to team groups
gac user update \
  --group team-backend \
  --group eng-all-hands \
  jsmith@example.com

# 4. Set custom fields
gac user update \
  --github-profile janesmith \
  --vpn-role developer \
  --id $(uuidgen) \
  jsmith@example.com
```

### Employee Offboarding

Complete workflow for offboarding a departing employee:

```bash
# 1. Suspend user account
gac user suspend \
  --reason "Left company - Last day 2025-10-15" \
  --force \
  former-employee@example.com

# 2. Clear PII (if required by policy)
gac user update \
  --clear-pii \
  former-employee@example.com

# 3. Move to Former Employees OU
gac user update \
  --ou "/Former Employees" \
  former-employee@example.com
```

### Department Transfer

Transfer user to new department:

```bash
# Update department, OU, and manager
gac user update \
  --dept Sales \
  --ou /Sales/East-Coast \
  --manager sales-manager@example.com \
  --title "Account Executive" \
  user@example.com

# Update group memberships
# (Note: You may need to manually remove from old dept groups)
gac user update \
  --group sales-team \
  --group east-coast-sales \
  user@example.com
```

### Bulk User Export

Export all user data for reporting:

```bash
# Export all users with full details
gac user list --full --csv > full-user-export.csv

# Export only active users
gac user list --csv > active-users.csv

# Export only suspended users
gac user list --disabled-only --csv > suspended-users.csv
```

## Best Practices

### Security

1. **Use suspension instead of deletion** - Suspend accounts to preserve data for compliance
2. **Document suspension reasons** - Always use `--reason` flag for audit trail
3. **Use force flag sparingly** - Confirmations prevent accidental changes
4. **Set employee IDs with UUIDs** - Use `$(uuidgen)` for unique identifiers

### Organization

1. **Consistent naming** - Use standard email format (firstname.lastname@domain.com)
2. **Use organizational units** - Organize users by department or location
3. **Group memberships** - Add users to appropriate groups during creation
4. **Keep manager hierarchy** - Always set manager for proper org chart

### Maintenance

1. **Regular audits** - Periodically review suspended accounts
2. **Clean up PII** - Use `--clear-pii` for departed employees per policy
3. **Update contact info** - Keep phone and address current
4. **Monitor disabled accounts** - Track accounts with `--disabled-only` flag

## Troubleshooting

### Common Issues

**Error: "User already exists"**
- Check if user email is already in use
- Verify you have the correct domain

**Error: "Invalid email format"**
- Ensure email follows user@domain.com format
- Check for typos or extra spaces

**Error: "Insufficient permissions"**
- Verify your OAuth credentials have User Admin scope
- Check you're authenticated with admin account

**Suspension not working**
- Confirm user email is correct
- Check you have necessary permissions
- Verify user isn't a super admin (can't be suspended)

For more help, see the [Troubleshooting Guide](../reference/troubleshooting.md).

## Related Documentation

- [Group Management](group-management.md) - Manage groups and memberships
- [Organizational Units](ou-management.md) - Manage organizational structure
- [Alias Management](alias-management.md) - Email aliases for users
- [Command Reference](../reference/commands.md) - Complete command documentation
