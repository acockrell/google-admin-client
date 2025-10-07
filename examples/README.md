# Configuration Examples

This directory contains example configuration files and usage scenarios for the Google Admin Client (gac).

## Table of Contents

- [Configuration Files](#configuration-files)
- [Usage Scenarios](#usage-scenarios)
- [Best Practices](#best-practices)

## Configuration Files

### Basic Configuration

**File**: [`basic-config.yaml`](basic-config.yaml)

Minimal configuration for getting started:

```yaml
# Google Workspace domain
domain: example.com

# OAuth2 credentials (optional - defaults to ~/.credentials/)
client-secret: ~/.credentials/client_secret.json
cache-file: ~/.credentials/gac.json
```

**Usage**:
```bash
cp examples/basic-config.yaml ~/.google-admin.yaml
# Edit with your domain
gac user list
```

---

### Production Configuration

**File**: [`production-config.yaml`](production-config.yaml)

Configuration for production use with custom paths:

```yaml
# Production Google Workspace domain
domain: company.com

# OAuth2 credentials in custom locations
client-secret: /etc/gac/client_secret.json
cache-file: /var/lib/gac/token.json

# Additional security: restrict to specific paths
```

**Usage**:
```bash
# Use with --config flag
gac --config /etc/gac/production-config.yaml user list
```

---

### Development Configuration

**File**: [`development-config.yaml`](development-config.yaml)

Configuration for development/testing:

```yaml
# Development domain
domain: dev.example.com

# Development credentials
client-secret: ~/.credentials/dev_client_secret.json
cache-file: ~/.credentials/dev_token.json
```

**Usage**:
```bash
# Keep separate configs for dev and prod
export GAC_CONFIG=~/.google-admin-dev.yaml
gac user list
```

## Usage Scenarios

### 1. Onboarding New Employee

**File**: [`onboarding-example.sh`](onboarding-example.sh)

Complete script for onboarding a new employee:

```bash
#!/bin/bash
# Onboard a new employee with all necessary setup

EMAIL="$1"
FIRST_NAME="$2"
LAST_NAME="$3"
DEPARTMENT="$4"
TITLE="$5"
MANAGER="$6"
PERSONAL_EMAIL="$7"

# Create user
echo "Creating user: $EMAIL"
gac user create \
    -f "$FIRST_NAME" \
    -l "$LAST_NAME" \
    -e "$PERSONAL_EMAIL" \
    -g all-staff \
    "$EMAIL"

# Update user profile
echo "Updating user profile..."
gac user update \
    --dept "$DEPARTMENT" \
    --title "$TITLE" \
    --manager "$MANAGER" \
    --type staff \
    --id "$(uuidgen)" \
    "$EMAIL"

# Add to department group
DEPT_GROUP=$(echo "$DEPARTMENT" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
echo "Adding to department group: $DEPT_GROUP"
gac user update -g "$DEPT_GROUP" "$EMAIL"

echo "Onboarding complete for $EMAIL"
```

**Usage**:
```bash
./examples/onboarding-example.sh \
    jdoe@example.com \
    "John" \
    "Doe" \
    "Engineering" \
    "Software Engineer" \
    "manager@example.com" \
    "john.doe@personal.com"
```

---

### 2. Bulk User Export

**File**: [`bulk-export.sh`](bulk-export.sh)

Export all users to CSV for backup or analysis:

```bash
#!/bin/bash
# Export all users to CSV with timestamp

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="users_export_${TIMESTAMP}.csv"

echo "Exporting all users to $OUTPUT_FILE..."
gac user list --csv --full > "$OUTPUT_FILE"

echo "Export complete: $OUTPUT_FILE"
echo "Total users: $(tail -n +2 "$OUTPUT_FILE" | wc -l)"
```

**Usage**:
```bash
./examples/bulk-export.sh
```

---

### 3. Offboarding Employee

**File**: [`offboarding-example.sh`](offboarding-example.sh)

Script for offboarding an employee:

```bash
#!/bin/bash
# Offboard an employee: transfer data and disable account

SOURCE_EMAIL="$1"
DESTINATION_EMAIL="$2"

if [ -z "$SOURCE_EMAIL" ] || [ -z "$DESTINATION_EMAIL" ]; then
    echo "Usage: $0 <source-email> <destination-email>"
    exit 1
fi

# Transfer documents and resources
echo "Transferring data from $SOURCE_EMAIL to $DESTINATION_EMAIL..."
gac transfer --from "$SOURCE_EMAIL" --to "$DESTINATION_EMAIL"

# Wait for transfer to complete
echo "Waiting for transfer to complete (check Google Admin console for status)"
read -p "Press enter when transfer is complete..."

# Disable the user account
echo "Disabling user account: $SOURCE_EMAIL"
gac user update --remove "$SOURCE_EMAIL"

# Clear PII
echo "Clearing personal information..."
gac user update --clear-pii "$SOURCE_EMAIL"

echo "Offboarding complete for $SOURCE_EMAIL"
```

**Usage**:
```bash
./examples/offboarding-example.sh olduser@example.com newowner@example.com
```

---

### 4. Group Audit

**File**: [`group-audit.sh`](group-audit.sh)

Audit groups for inactive members:

```bash
#!/bin/bash
# Find groups containing former employees

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_FILE="group_audit_${TIMESTAMP}.txt"

echo "Auditing groups for inactive members..."
echo "Group Audit Report - $TIMESTAMP" > "$OUTPUT_FILE"
echo "======================================" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

gac group list --contains-former-employees >> "$OUTPUT_FILE"

echo "Audit complete: $OUTPUT_FILE"
cat "$OUTPUT_FILE"
```

**Usage**:
```bash
./examples/group-audit.sh
```

---

### 5. Calendar Event Creation

**File**: [`create-recurring-meeting.sh`](create-recurring-meeting.sh)

Create recurring team meetings:

```bash
#!/bin/bash
# Create a recurring team meeting

CALENDAR_EMAIL="$1"
MEETING_TITLE="$2"
START_DATE="$3"  # Format: 2025-10-15T09:00:00-04:00
END_DATE="$4"    # Format: 2025-10-15T10:00:00-04:00
FREQUENCY="${5:-weekly}"  # daily, weekly, or monthly
COUNT="${6:-52}"  # Number of occurrences

if [ -z "$CALENDAR_EMAIL" ] || [ -z "$MEETING_TITLE" ] || [ -z "$START_DATE" ] || [ -z "$END_DATE" ]; then
    echo "Usage: $0 <calendar-email> <meeting-title> <start-date> <end-date> [frequency] [count]"
    echo "Example: $0 team@example.com 'Team Standup' '2025-10-15T09:00:00-04:00' '2025-10-15T09:30:00-04:00' daily 90"
    exit 1
fi

echo "Creating recurring meeting: $MEETING_TITLE"
gac calendar create "$CALENDAR_EMAIL" \
    -s "$MEETING_TITLE" \
    -b "$START_DATE" \
    -e "$END_DATE" \
    -f "$FREQUENCY" \
    -c "$COUNT" \
    -l "Virtual - Zoom"

echo "Meeting created successfully"
```

**Usage**:
```bash
./examples/create-recurring-meeting.sh \
    team@example.com \
    "Daily Standup" \
    "2025-10-15T09:00:00-04:00" \
    "2025-10-15T09:15:00-04:00" \
    daily \
    90
```

---

### 6. User Update Batch Script

**File**: [`batch-update-users.sh`](batch-update-users.sh)

Update multiple users from a CSV file:

```bash
#!/bin/bash
# Batch update users from CSV
# CSV format: email,department,title,manager

INPUT_FILE="${1:-users-to-update.csv}"

if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: File not found: $INPUT_FILE"
    echo "Usage: $0 [csv-file]"
    exit 1
fi

echo "Updating users from: $INPUT_FILE"
echo ""

# Skip header line and process each user
tail -n +2 "$INPUT_FILE" | while IFS=, read -r EMAIL DEPT TITLE MANAGER; do
    echo "Updating: $EMAIL"
    echo "  Department: $DEPT"
    echo "  Title: $TITLE"
    echo "  Manager: $MANAGER"

    gac user update \
        --dept "$DEPT" \
        --title "$TITLE" \
        --manager "$MANAGER" \
        "$EMAIL"

    if [ $? -eq 0 ]; then
        echo "  ✓ Updated successfully"
    else
        echo "  ✗ Update failed"
    fi
    echo ""
done

echo "Batch update complete"
```

**Sample CSV** (`users-to-update.csv`):
```csv
email,department,title,manager
jdoe@example.com,Engineering,Senior Engineer,manager@example.com
jsmith@example.com,Sales,Account Executive,sales-mgr@example.com
```

**Usage**:
```bash
./examples/batch-update-users.sh users-to-update.csv
```

## Best Practices

### Security

1. **Protect credentials**:
   ```bash
   chmod 600 ~/.credentials/*.json
   chmod 700 ~/.credentials
   ```

2. **Use environment-specific configs**:
   - Keep separate credentials for dev, staging, and production
   - Never commit credentials to version control
   - Use `.gitignore` to exclude credential files

3. **Rotate credentials regularly**:
   - Create new OAuth2 clients periodically
   - Revoke old tokens in Google Cloud Console

### Configuration Management

1. **Use configuration files** for persistent settings:
   ```yaml
   # ~/.google-admin.yaml
   domain: example.com
   ```

2. **Use environment variables** for CI/CD:
   ```bash
   export GAC_DOMAIN=example.com
   export GAC_CLIENT_SECRET=/path/to/credentials.json
   ```

3. **Use command-line flags** for one-off overrides:
   ```bash
   gac --domain dev.example.com user list
   ```

### Scripting

1. **Check exit codes**:
   ```bash
   if gac user create newuser@example.com; then
       echo "Success"
   else
       echo "Failed"
       exit 1
   fi
   ```

2. **Add error handling**:
   ```bash
   set -euo pipefail  # Exit on error, undefined vars, pipe failures
   ```

3. **Log operations**:
   ```bash
   gac user create newuser@example.com 2>&1 | tee -a operations.log
   ```

4. **Use dry-run when available** (future feature):
   ```bash
   gac user update --dry-run --dept Engineering user@example.com
   ```

### Performance

1. **Batch operations** when possible to reduce API calls

2. **Add delays** between bulk operations to avoid rate limits:
   ```bash
   for user in $(cat users.txt); do
       gac user update --dept Engineering "$user"
       sleep 1  # Delay between requests
   done
   ```

3. **Monitor API quotas** in Google Cloud Console

### Testing

1. **Test on non-production** domains first

2. **Validate inputs** before bulk operations:
   ```bash
   # Validate email format
   if [[ ! "$EMAIL" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
       echo "Invalid email: $EMAIL"
       exit 1
   fi
   ```

3. **Keep backups** before destructive operations:
   ```bash
   gac user list --csv --full > backup_$(date +%Y%m%d).csv
   ```

## Additional Resources

- [Main README](../README.md) - Full documentation
- [CREDENTIALS.md](../CREDENTIALS.md) - OAuth2 setup guide
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Development guidelines
- [ARCHITECTURE.md](../ARCHITECTURE.md) - Technical details

## Contributing Examples

Have a useful script or configuration? We welcome contributions!

1. Add your example to this directory
2. Update this README with documentation
3. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.
