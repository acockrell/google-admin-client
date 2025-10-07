# gac - Google Admin Client

A command-line tool for managing Google Workspace (formerly Google Apps) users, groups, calendars, and data transfers.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage](#usage)
  - [User Management](#user-management)
  - [Group Management](#group-management)
  - [Calendar Operations](#calendar-operations)
  - [Organizational Unit Management](#organizational-unit-management)
  - [Data Transfers](#data-transfers)
- [Command Reference](#command-reference)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Architecture](#architecture)
- [License](#license)

## Overview

**gac** (Google Admin Client) is a powerful CLI tool for automating Google Workspace administrative tasks. Built with Go and the Cobra framework, it provides a simple interface for managing users, groups, calendars, and data transfers through the Google Admin SDK APIs.

### Key Features

- **User Management**: Create, list, and update users with comprehensive profile support
- **Group Management**: List groups and manage memberships
- **Calendar Operations**: Create, list, and update calendar events
- **Data Transfers**: Transfer ownership of documents and resources between users
- **Secure Authentication**: OAuth2 authentication with automatic token refresh
- **Flexible Configuration**: Support for config files, environment variables, and CLI flags
- **Input Validation**: Comprehensive validation of emails, phone numbers, UUIDs, and other inputs
- **Cross-Platform**: Builds for Linux and macOS (amd64 and arm64)

## Installation

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/acockrell/google-admin-client/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_darwin_arm64.tar.gz
tar xzf gac_darwin_arm64.tar.gz
sudo mv gac /usr/local/bin/

# macOS (Intel)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_darwin_amd64.tar.gz
tar xzf gac_darwin_amd64.tar.gz
sudo mv gac /usr/local/bin/

# Linux (amd64)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_linux_amd64.tar.gz
tar xzf gac_linux_amd64.tar.gz
sudo mv gac /usr/local/bin/

# Linux (arm64)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_linux_arm64.tar.gz
tar xzf gac_linux_arm64.tar.gz
sudo mv gac /usr/local/bin/
```

### Build from Source

Requirements:
- Go 1.25 or later
- Git

```bash
git clone https://github.com/acockrell/google-admin-client.git
cd google-admin-client
make build
sudo mv build/gac /usr/local/bin/
```

### Docker

```bash
docker pull ghcr.io/acockrell/google-admin-client:latest
docker run --rm -it ghcr.io/acockrell/google-admin-client:latest --help
```

### Verify Installation

```bash
gac --help
```

## Quick Start

### 1. Set Up Google Cloud Credentials

Before using `gac`, you need to set up OAuth2 credentials:

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create or select a project
3. Enable the Admin SDK API and Google Calendar API
4. Create OAuth2 credentials (Desktop app type)
5. Download the JSON credentials file

See [CREDENTIALS.md](CREDENTIALS.md) for detailed setup instructions.

### 2. Configure Credentials

Place your credentials in the default location:

```bash
mkdir -p ~/.credentials
mv ~/Downloads/client_secret_*.json ~/.credentials/client_secret.json
chmod 600 ~/.credentials/client_secret.json
```

### 3. Authenticate

Run any command to trigger the OAuth2 authentication flow:

```bash
gac user list
```

Follow the browser prompts to authenticate and grant permissions. Your token will be saved to `~/.credentials/gac.json`.

### 4. Run Commands

```bash
# List all users
gac user list

# Create a new user
gac user create newuser@example.com

# Update user information
gac user update --dept Engineering --title "Software Engineer" jdoe@example.com

# List groups
gac group list
```

## Configuration

`gac` supports multiple configuration methods with the following priority (highest to lowest):

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration file**
4. **Default values** (lowest priority)

### Configuration File

Create `~/.google-admin.yaml` with your settings:

```yaml
# Google Workspace domain
domain: example.com

# OAuth2 credential paths
client-secret: /path/to/client_secret.json
cache-file: /path/to/token.json
```

### Environment Variables

```bash
# Primary environment variables
export GAC_DOMAIN=example.com
export GAC_CLIENT_SECRET=/path/to/client_secret.json
export GAC_CACHE_FILE=/path/to/token.json

# Alternate environment variables (also supported)
export GOOGLE_ADMIN_DOMAIN=example.com
export GOOGLE_ADMIN_CLIENT_SECRET=/path/to/client_secret.json
export GOOGLE_ADMIN_CACHE_FILE=/path/to/token.json
```

### Global Flags

All commands support these global flags:

- `--domain string` - Domain for email addresses (e.g., example.com)
- `--client-secret string` - Path to OAuth2 client secret JSON file
- `--cache-file string` - Path to OAuth2 token cache file
- `--config string` - Path to config file (default: `$HOME/.google-admin.yaml`)

### Examples

See the [examples/](examples/) directory for sample configuration files and use cases.

## Usage

### User Management

#### Create a User

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

**Flags:**
- `-f, --first-name` - First name
- `-l, --last-name` - Last name
- `-e, --email` - Personal email address
- `-g, --groups` - Groups to add user to (can be repeated)

#### List Users

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

**Flags:**
- `-c, --csv` - Export as CSV
- `-d, --disabled-only` - Show only disabled accounts
- `-f, --full` - Include all user fields

#### Update a User

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

**Flags:**
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

### Group Management

#### List Groups

```bash
# List all groups
gac group list

# List a specific group
gac group list operations@example.com

# List group with members
gac group list operations@example.com --get-members

# Show only groups with inactive members
gac group list --contains-former-employees
```

**Flags:**
- `-m, --get-members` - List group members
- `-i, --contains-former-employees` - Show only groups with inactive members

### Calendar Operations

#### Create Calendar Event

```bash
# All-day event
gac calendar create user@example.com \
  -s "Team Building Day" \
  -b 2025-10-15 \
  -e 2025-10-16

# Timed event with description
gac calendar create user@example.com \
  -s "Sprint Planning" \
  -d "Q4 2025 Sprint Planning Meeting" \
  -b 2025-10-15T09:00:00-04:00 \
  -e 2025-10-15T11:00:00-04:00 \
  -l "Conference Room A" \
  -a attendee1@example.com \
  -a attendee2@example.com

# Recurring event
gac calendar create user@example.com \
  -s "Daily Standup" \
  -b 2025-10-15T09:00:00-04:00 \
  -e 2025-10-15T09:15:00-04:00 \
  -f daily \
  -c 30
```

**Flags:**
- `-s, --summary` - Event title (required)
- `-d, --description` - Event description
- `-b, --begin` - Event start time in RFC3339 format (required)
- `-e, --end` - Event end time in RFC3339 format (required)
- `-l, --location` - Event location (default: "The Matrix")
- `-a, --attendee` - Event attendees (can be repeated)
- `-f, --frequency` - Recurrence frequency: daily, weekly, monthly (default: "daily")
- `-c, --count` - Number of recurrences (default: 1)

#### List Calendar Events

```bash
# List next 10 events
gac calendar list user@example.com

# List next 50 events
gac calendar list user@example.com -n 50

# List events in date range
gac calendar list user@example.com \
  --time-min 2025-10-01T00:00:00-04:00 \
  --time-max 2025-10-31T23:59:59-04:00
```

**Flags:**
- `-n, --num-events` - Number of events to return (default: 10)
- `--time-min` - Minimum event start time (RFC3339 format)
- `--time-max` - Maximum event start time (RFC3339 format)

#### Update Calendar Event

```bash
gac calendar update user@example.com \
  -i event_id_here \
  -s "Updated Meeting Title" \
  -b 2025-10-15T10:00:00-04:00 \
  -e 2025-10-15T11:00:00-04:00
```

**Flags:**
- `-i, --event-id` - ID of event to update (required)
- `-s, --summary` - Updated event title
- `-d, --description` - Updated description
- `-b, --begin` - Updated start time
- `-e, --end` - Updated end time
- `-l, --location` - Updated location
- `-a, --attendee` - Updated attendees

### Organizational Unit Management

Organizational units (OUs) allow you to organize users and apply different policies to different groups.

#### List Organizational Units

```bash
# List all organizational units
gac ou list

# List specific OU and its children
gac ou list /Engineering

# List only direct children
gac ou list /Engineering --type children
```

**Flags:**
- `-t, --type` - List type: all or children (default: "all")

#### Create Organizational Unit

```bash
# Create a top-level OU
gac ou create /Engineering --description "Engineering department"

# Create a nested OU
gac ou create /Engineering/Backend --description "Backend engineering team"

# Create with inheritance blocking
gac ou create /Contractors --block-inheritance
```

**Flags:**
- `-d, --description` - Organizational unit description
- `-p, --parent` - Parent OU path (auto-detected from path if not specified)
- `-b, --block-inheritance` - Block policy inheritance from parent

#### Update Organizational Unit

```bash
# Update description
gac ou update /Engineering --description "Updated description"

# Rename an OU
gac ou update /Engineering --name "Engineering-Dept"

# Move an OU to a different parent
gac ou update /Engineering/QA --parent /Operations

# Enable inheritance blocking
gac ou update /Contractors --block-inheritance true
```

**Flags:**
- `-n, --name` - New name for the organizational unit
- `-d, --description` - New description
- `-p, --parent` - New parent OU path
- `-b, --block-inheritance` - Block policy inheritance (true/false)

#### Delete Organizational Unit

```bash
# Delete an empty OU (with confirmation prompt)
gac ou delete /Engineering/Archived

# Force delete without confirmation
gac ou delete /TempOU --force
```

**Flags:**
- `-f, --force` - Skip confirmation prompt

**Note:** The OU must be empty (no users or sub-OUs) before it can be deleted.

### Data Transfers

Transfer document ownership from one user to another:

```bash
gac transfer --from olduser@example.com --to newuser@example.com
```

**Flags:**
- `-f, --from` - Source user email address (required)
- `-t, --to` - Destination user email address (required)

## Command Reference

### Global Commands

| Command | Description |
|---------|-------------|
| `gac --help` | Show help for gac |
| `gac version` | Show version information |
| `gac completion` | Generate shell completion scripts |

### User Commands

| Command | Description |
|---------|-------------|
| `gac user create [email]` | Create a new user |
| `gac user list [email]` | List users or get details for specific user |
| `gac user update [email]` | Update user information |

### Group Commands

| Command | Description |
|---------|-------------|
| `gac group list [email]` | List groups or get details for specific group |

### Calendar Commands

| Command | Description |
|---------|-------------|
| `gac calendar create [email]` | Create a calendar event |
| `gac calendar list [email]` | List calendar events |
| `gac calendar update [email]` | Update a calendar event |

### Organizational Unit Commands

| Command | Description |
|---------|-------------|
| `gac ou list [ou-path]` | List organizational units |
| `gac ou create <ou-path>` | Create a new organizational unit |
| `gac ou update <ou-path>` | Update an organizational unit |
| `gac ou delete <ou-path>` | Delete an organizational unit |

### Transfer Commands

| Command | Description |
|---------|-------------|
| `gac transfer --from [email] --to [email]` | Transfer data ownership |

## Troubleshooting

### Authentication Issues

**Problem**: "Invalid client secret path" or "Access denied"

**Solution**:
1. Verify your OAuth2 credentials are correctly set up in Google Cloud Console
2. Check that the client secret file exists and has correct permissions:
   ```bash
   ls -l ~/.credentials/client_secret.json
   chmod 600 ~/.credentials/client_secret.json
   ```
3. Ensure you're authenticating with a Google Workspace admin account
4. Verify the required APIs are enabled in Google Cloud Console:
   - Admin SDK API
   - Google Calendar API

**Problem**: "Token expired" or "Invalid token"

**Solution**:
```bash
# Delete the cached token and re-authenticate
rm ~/.credentials/gac.json
gac user list
```

### Permission Issues

**Problem**: "Permission denied" or "Insufficient permissions"

**Solution**:
1. Verify the authenticated account has Google Workspace admin privileges
2. Check that all required OAuth2 scopes are granted (see [CREDENTIALS.md](CREDENTIALS.md))
3. Try deleting and re-creating your OAuth2 credentials in Google Cloud Console

### File Permission Warnings

**Problem**: Warnings about insecure file permissions

**Solution**:
```bash
# Fix credential directory permissions
chmod 700 ~/.credentials

# Fix credential file permissions
chmod 600 ~/.credentials/*.json
```

### Domain Configuration Issues

**Problem**: Commands fail with domain-related errors

**Solution**:
1. Set your domain in configuration:
   ```bash
   # Via environment variable
   export GAC_DOMAIN=example.com

   # Via config file
   echo "domain: example.com" > ~/.google-admin.yaml

   # Via command-line flag
   gac --domain example.com user list
   ```

### Input Validation Errors

**Problem**: "Invalid email address" or "Invalid phone format"

**Solution**:
- **Emails**: Must be valid RFC 5322 format (e.g., `user@example.com`)
- **Phone numbers**: Format as `type:number` (e.g., `mobile:555-555-5555`)
  - Multiple phones: `mobile:555-555-5555; work:555-123-4567,555`
- **UUIDs**: Must be valid UUID format (use `uuidgen` on macOS/Linux)

### API Rate Limiting

**Problem**: "Rate limit exceeded" or "Quota exceeded"

**Solution**:
1. Implement delays between bulk operations
2. Reduce the frequency of API calls
3. Check your quota limits in Google Cloud Console
4. Consider requesting a quota increase if needed

### Getting Help

If you encounter issues not covered here:

1. Check the [CREDENTIALS.md](CREDENTIALS.md) for authentication details
2. Review the [ARCHITECTURE.md](ARCHITECTURE.md) for technical details
3. Check existing [GitHub Issues](https://github.com/acockrell/google-admin-client/issues)
4. Open a new issue with:
   - Command you ran
   - Error message
   - Your configuration (redact sensitive info)
   - Steps to reproduce

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:

- Reporting bugs
- Suggesting features
- Submitting pull requests
- Development setup
- Code style and testing requirements

## Architecture

For technical details about the project structure, design decisions, and extension points, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Related Documentation

- [CREDENTIALS.md](CREDENTIALS.md) - OAuth2 setup and security practices
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development and contribution guidelines
- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture documentation
- [DEBUGGING.md](DEBUGGING.md) - Debugging guide for developers
- [RELEASE.md](RELEASE.md) - Release process documentation

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration management
- Powered by [Google Admin SDK](https://developers.google.com/admin-sdk)
