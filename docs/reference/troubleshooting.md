# Troubleshooting

Common issues and solutions for `gac`.

## Table of Contents

- [Authentication Issues](#authentication-issues)
- [Permission Issues](#permission-issues)
- [File Permission Warnings](#file-permission-warnings)
- [Domain Configuration Issues](#domain-configuration-issues)
- [Input Validation Errors](#input-validation-errors)
- [API Rate Limiting](#api-rate-limiting)
- [Getting Help](#getting-help)

## Authentication Issues

### Problem: "Invalid client secret path" or "Access denied"

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
   - Groups Settings API

### Problem: "Token expired" or "Invalid token"

**Solution**:
```bash
# Delete the cached token and re-authenticate
rm ~/.credentials/gac.json
gac user list
```

### Problem: OAuth consent screen issues

**Solution**:
1. Check that your OAuth consent screen is configured in Google Cloud Console
2. Ensure your app is set to "Internal" for Google Workspace domains
3. Verify all required scopes are added to the consent screen
4. Try deleting and recreating the OAuth client ID

## Permission Issues

### Problem: "Permission denied" or "Insufficient permissions"

**Solution**:
1. Verify the authenticated account has Google Workspace admin privileges
2. Check that all required OAuth2 scopes are granted (see [Authentication Guide](../authentication.md))
3. Try deleting and re-creating your OAuth2 credentials in Google Cloud Console
4. Ensure you have the specific admin role needed:
   - **User Admin** - For user operations
   - **Groups Admin** - For group operations
   - **Organizational Units Admin** - For OU operations

### Problem: "Groups Settings API access denied"

**Solution**:
1. Verify the Groups Settings API is enabled in Google Cloud Console
2. Check that the scope `https://www.googleapis.com/auth/apps.groups.settings` is included
3. Delete cached token and re-authenticate:
   ```bash
   rm ~/.credentials/gac.json
   gac group-settings list <group-email>
   ```

## File Permission Warnings

### Problem: Warnings about insecure file permissions

**Solution**:
```bash
# Fix credential directory permissions
chmod 700 ~/.credentials

# Fix credential file permissions
chmod 600 ~/.credentials/*.json
```

### Problem: "Cannot write to token file"

**Solution**:
```bash
# Ensure credentials directory exists
mkdir -p ~/.credentials

# Fix directory permissions
chmod 700 ~/.credentials

# Verify you have write access
ls -la ~/.credentials
```

## Domain Configuration Issues

### Problem: Commands fail with domain-related errors

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

2. Verify domain is correct:
   ```bash
   # Check what domain is being used
   gac user list | head -1
   ```

### Problem: "Group not found" when using short names

**Solution**:
- Use full email format: `team@example.com` instead of just `team`
- Or set the domain globally so short names work:
  ```bash
  export GAC_DOMAIN=example.com
  gac group list team  # Now works
  ```

## Input Validation Errors

### Problem: "Invalid email address"

**Solution**:
- Emails must be valid RFC 5322 format: `user@example.com`
- No spaces, special characters in local part
- Must include `@` and valid domain

**Valid examples**:
```
john.doe@example.com
jane-smith@example.com
user+tag@example.com
```

**Invalid examples**:
```
john doe@example.com    # No spaces
@example.com            # No local part
user@                   # No domain
```

### Problem: "Invalid phone format"

**Solution**:
- Format as `type:number`
- Supported types: `mobile`, `work`, `home`
- Extensions: `work:555-123-4567,555`
- Multiple phones: `mobile:555-555-5555; work:555-123-4567,555`

**Valid examples**:
```
mobile:555-555-5555
work:555-123-4567,555
mobile:555-555-5555; work:555-123-4567
```

### Problem: "Invalid UUID format"

**Solution**:
- Use `uuidgen` to generate valid UUIDs:
  ```bash
  gac user update --id $(uuidgen) user@example.com
  ```
- UUID format: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

## API Rate Limiting

### Problem: "Rate limit exceeded" or "Quota exceeded"

**Solution**:
1. **Implement delays** between bulk operations:
   ```bash
   for user in $(cat users.txt); do
     gac user create $user
     sleep 1  # Wait 1 second between calls
   done
   ```

2. **Reduce frequency** of API calls:
   - Use batch operations when possible
   - Cache results locally instead of re-fetching

3. **Check quotas** in Google Cloud Console:
   - Go to APIs & Services > Quotas
   - Look for Admin SDK API quotas
   - Consider requesting quota increases

4. **Use exponential backoff**:
   ```bash
   # Retry with exponential backoff
   retry_count=0
   max_retries=3
   while [ $retry_count -lt $max_retries ]; do
     if gac user create user@example.com; then
       break
     fi
     retry_count=$((retry_count + 1))
     sleep $((2 ** retry_count))
   done
   ```

### Problem: "Concurrent requests limit exceeded"

**Solution**:
- Reduce concurrency in batch operations
- Run operations sequentially instead of in parallel
- Add delays between requests

## Common Command Errors

### Problem: User creation fails

**Possible causes**:
- Email already exists
- Invalid email format
- Missing required permissions
- Domain not allowed

**Solution**:
```bash
# Check if user exists
gac user list user@example.com

# Verify email format
# Must be: user@example.com

# Check your permissions
# Need User Admin role
```

### Problem: Group settings not updating

**Possible causes**:
- Wrong group email
- Missing Groups Settings API scope
- Insufficient permissions

**Solution**:
```bash
# Use full email format
gac group-settings list team@example.com

# Check settings were applied
gac group-settings list team@example.com --format json

# Re-authenticate with correct scopes
rm ~/.credentials/gac.json
gac group-settings list team@example.com
```

### Problem: Calendar operations fail

**Possible causes**:
- Calendar API not enabled
- Wrong calendar email
- Invalid date format

**Solution**:
1. Enable Calendar API in Google Cloud Console
2. Use correct date formats:
   - All-day: `2025-10-15`
   - Timed: `2025-10-15T09:00:00-04:00` (RFC3339)
3. Verify calendar email exists

## Getting Help

### Self-Service

1. **Check documentation**:
   - [Authentication Guide](../authentication.md) - OAuth setup
   - [User Guide](../guides/user-management.md) - User operations
   - [Command Reference](commands.md) - All commands

2. **Review examples**:
   - [Examples directory](../../examples/) - Working scripts
   - [Examples README](../../examples/README.md) - Scenario walkthroughs

3. **Verify setup**:
   ```bash
   # Test authentication
   gac user list

   # Check configuration
   cat ~/.google-admin.yaml

   # Verify credentials
   ls -la ~/.credentials/
   ```

### Getting Support

If you encounter issues not covered here:

1. **Search existing issues**:
   - Check [GitHub Issues](https://github.com/acockrell/google-admin-client/issues)
   - Look for similar problems and solutions

2. **Open a new issue** with:
   - **Command you ran** - Exact command that failed
   - **Error message** - Full error output
   - **Your configuration** - Redact sensitive info (tokens, emails)
   - **Steps to reproduce** - Minimal example that shows the issue
   - **Environment** - OS, Go version, gac version

3. **Include diagnostic info**:
   ```bash
   # Version info
   gac version

   # Check file permissions
   ls -la ~/.credentials/

   # Verify domain config
   cat ~/.google-admin.yaml
   ```

### Debug Mode

For verbose output, use environment variables:

```bash
# Enable debug logging (if implemented)
export GAC_DEBUG=true
gac user list

# Check what API calls are being made
# (Requires modifying code to add logging)
```

## Related Documentation

- [Authentication Guide](../authentication.md)
- [Command Reference](commands.md)
- [Examples](../../examples/README.md)
- [Contributing Guide](../development/contributing.md)
