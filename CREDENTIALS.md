# Credential Management Guide

This document describes secure practices for managing credentials when using the Google Admin Client (gac).

## Overview

The Google Admin Client uses OAuth2 for authentication with Google Workspace APIs. This requires two types of credentials:

1. **Client Secret**: OAuth2 client credentials from Google Cloud Console
2. **Access Token**: OAuth2 token stored locally after initial authentication

## OAuth2 Scope Requirements

The tool requires the following Google Workspace API scopes:

### Admin Directory API
- `https://www.googleapis.com/auth/admin.directory.user.readonly` - Read user information
- `https://www.googleapis.com/auth/admin.directory.user` - Manage users
- `https://www.googleapis.com/auth/admin.directory.group.readonly` - Read group information
- `https://www.googleapis.com/auth/admin.directory.group.member.readonly` - Read group membership
- `https://www.googleapis.com/auth/admin.directory.group.member` - Manage group membership

### Calendar API
- `https://www.googleapis.com/auth/calendar` - Full calendar access
- `https://www.googleapis.com/auth/calendar.readonly` - Read-only calendar access
- `https://www.googleapis.com/auth/calendar.events` - Manage calendar events
- `https://www.googleapis.com/auth/calendar.events.readonly` - Read calendar events

### Data Transfer API
- `https://www.googleapis.com/auth/admin.datatransfer` - Manage data transfers

These scopes are configured in `cmd/client.go:28-39`. If you modify these scopes, you must delete your previously saved token at `~/.credentials/gac.json` to re-authenticate.

## Setting Up OAuth2 Credentials

### 1. Create OAuth2 Client in Google Cloud Console

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create or select a project
3. Enable the following APIs:
   - Admin SDK API
   - Google Calendar API
4. Go to "APIs & Services" > "Credentials"
5. Click "Create Credentials" > "OAuth client ID"
6. Select "Desktop app" as application type
7. Download the JSON file

### 2. Configure Credentials

The tool supports multiple methods for providing credentials:

#### Option 1: Default Location (Recommended)
Place the client secret file at:
```
~/.credentials/client_secret.json
```

The token will be automatically saved to:
```
~/.credentials/gac.json
```

#### Option 2: Custom Location via Configuration
Set the credential path in your `.google-admin.yaml`:
```yaml
client-secret: /path/to/client_secret.json
cache-file: /path/to/token.json
```

#### Option 3: Environment Variables
```bash
export GAC_CLIENT_SECRET=/path/to/client_secret.json
export GAC_CACHE_FILE=/path/to/token.json
```

#### Option 4: Command-line Flags
```bash
gac --client-secret /path/to/client_secret.json user list
```

## Secure Credential Storage Practices

### File Permissions

**Critical**: Credential files should have restrictive permissions to prevent unauthorized access.

#### Recommended Permissions
```bash
# Set permissions on credential directory
chmod 700 ~/.credentials

# Set permissions on credential files
chmod 600 ~/.credentials/client_secret.json
chmod 600 ~/.credentials/gac.json
```

The tool will automatically:
- Create the `~/.credentials` directory with `0700` permissions (owner read/write/execute only)
- Warn if credential files have overly permissive permissions (world-readable or group-readable)

#### Checking Current Permissions
```bash
ls -la ~/.credentials
```

Expected output:
```
drwx------  2 user user 4096 Jan 01 12:00 .
-rw-------  1 user user  623 Jan 01 12:00 client_secret.json
-rw-------  1 user user  418 Jan 01 12:00 gac.json
```

### Security Measures

The tool implements several security measures for credential handling:

1. **Path Validation**: All credential file paths are validated to prevent directory traversal attacks
   - Paths must be within the user's home directory or temp directory
   - Paths containing `..` sequences are rejected

2. **File Permission Checks**: Warnings are issued for insecure file permissions
   - Files readable by group or world will trigger warnings
   - Recommendation to use `chmod 600` on credential files

3. **Secure Storage Location**: Credentials are stored in user-specific locations
   - Default: `~/.credentials/` directory with `0700` permissions
   - Only the file owner can access credential files

### What NOT to Do

- **Never commit credentials to version control** (add `.credentials/` to `.gitignore`)
- **Never share credential files** via email, chat, or file sharing services
- **Never store credentials in world-readable locations** like `/tmp/` with default permissions
- **Never use production credentials in development** without proper safeguards
- **Never grant broader API scopes than necessary** - use the minimum required scopes

### Best Practices

1. **Use separate credentials for different environments** (development, staging, production)
2. **Rotate credentials periodically** by creating new OAuth2 clients
3. **Revoke unused tokens** in Google Cloud Console under "APIs & Services" > "Credentials"
4. **Enable Google Workspace domain-wide delegation** for service accounts when appropriate
5. **Monitor API usage** in Google Cloud Console to detect unusual activity
6. **Use service accounts** for automated/server-side operations instead of user credentials
7. **Store credentials in encrypted filesystems** when possible
8. **Use environment-specific configuration files** that are not checked into version control

## Initial Authentication

On first run, the tool will:

1. Check for existing token at `~/.credentials/gac.json`
2. If no token exists, prompt for OAuth2 authorization:
   ```
   Go to the following link in your browser then type the authorization code:
   https://accounts.google.com/o/oauth2/auth?...
   ```
3. Open the link in your browser
4. Authenticate with a Google Workspace admin account
5. Grant the requested permissions
6. Copy the authorization code
7. Paste the code into the terminal
8. Token is saved to `~/.credentials/gac.json` with `0600` permissions

The token will be automatically refreshed when it expires (typically 1 hour for access tokens, but refresh tokens are long-lived).

## Troubleshooting

### "Invalid client secret path" Error
- Ensure the client secret file exists at the specified location
- Verify file permissions (must be readable by current user)
- Check that path doesn't contain directory traversal sequences (`..`)

### "Access denied" Errors
- Verify the authenticated user has Google Workspace admin privileges
- Check that all required API scopes are granted
- Confirm the APIs are enabled in Google Cloud Console

### "Token expired" Errors
- Delete the token file: `rm ~/.credentials/gac.json`
- Re-authenticate by running any command

### Permission Warnings
If you see warnings about insecure file permissions:
```bash
# Fix credential directory permissions
chmod 700 ~/.credentials

# Fix individual file permissions
chmod 600 ~/.credentials/*.json
```

