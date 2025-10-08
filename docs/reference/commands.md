# Command Reference

Complete reference for all `gac` commands.

## Global Commands

| Command | Description |
|---------|-------------|
| `gac --help` | Show help for gac |
| `gac version` | Show version information |
| `gac completion` | Generate shell completion scripts |

## User Commands

| Command | Description |
|---------|-------------|
| `gac user create [email]` | Create a new user |
| `gac user list [email]` | List users or get details for specific user |
| `gac user update [email]` | Update user information |
| `gac user suspend <user-email>` | Suspend a user account |
| `gac user unsuspend <user-email>` | Unsuspend (restore) a user account |

See: [User Management Guide](../guides/user-management.md)

## Group Commands

| Command | Description |
|---------|-------------|
| `gac group list [email]` | List groups or get details for specific group |

See: [Group Management Guide](../guides/group-management.md)

## Group Settings Commands

| Command | Description |
|---------|-------------|
| `gac group-settings list <group-email>` | View group settings |
| `gac group-settings update <group-email>` | Update group settings |

See: [Group Settings Guide](../guides/group-settings.md)

## Calendar Commands

| Command | Description |
|---------|-------------|
| `gac calendar create [email]` | Create a calendar event |
| `gac calendar list [email]` | List calendar events |
| `gac calendar update [email]` | Update a calendar event |

See: [Calendar Operations Guide](../guides/calendar-operations.md)

## Calendar Resource Commands

| Command | Description |
|---------|-------------|
| `gac cal-resource list` | List calendar resources (rooms, equipment) |
| `gac cal-resource create <resource-id>` | Create a new calendar resource |
| `gac cal-resource update <resource-id>` | Update a calendar resource |
| `gac cal-resource delete <resource-id>` | Delete a calendar resource |

See: [Calendar Resources Guide](../guides/calendar-resources.md)

## Organizational Unit Commands

| Command | Description |
|---------|-------------|
| `gac ou list [ou-path]` | List organizational units |
| `gac ou create <ou-path>` | Create a new organizational unit |
| `gac ou update <ou-path>` | Update an organizational unit |
| `gac ou delete <ou-path>` | Delete an organizational unit |

See: [Organizational Units Guide](../guides/ou-management.md)

## Alias Commands

| Command | Description |
|---------|-------------|
| `gac alias list <user-email>` | List aliases for a user |
| `gac alias add <user-email> <alias-email>` | Add an alias to a user |
| `gac alias remove <user-email> <alias-email>` | Remove an alias from a user |

See: [Alias Management Guide](../guides/alias-management.md)

## Audit Commands

| Command | Description |
|---------|-------------|
| `gac audit export` | Export audit logs for a specific application |

### Audit Export Flags

| Flag | Description |
|------|-------------|
| `--app <application>` | Application type (required: admin, login, drive, calendar, groups, etc.) |
| `--start-time <time>` | Start time in RFC3339 format (default: 24 hours ago) |
| `--end-time <time>` | End time in RFC3339 format (default: now) |
| `--user <email>` | Filter by user email address |
| `--event-name <event>` | Filter by event name (can specify multiple) |
| `--actor-ip <ip>` | Filter by actor IP address |
| `-o, --output <format>` | Output format (json, csv) (default: json) |
| `-f, --output-file <path>` | Output file path (default: stdout) |
| `--max-results <number>` | Maximum number of results to return |

### Supported Application Types

- `admin` - Admin console activities
- `login` - User login/logout activities
- `drive` - Google Drive operations
- `calendar` - Calendar activities
- `groups` - Google Groups operations
- `mobile` - Mobile device activities
- `token` - OAuth token activities
- `groups_enterprise` - Groups for Enterprise
- `saml` - SAML authentication
- `chrome` - Chrome browser management
- `gcp` - Google Cloud Platform
- `chat` - Google Chat
- `meet` - Google Meet

See: [Audit Logs Guide](../guides/audit-logs.md)

## Transfer Commands

| Command | Description |
|---------|-------------|
| `gac transfer --from [email] --to [email]` | Transfer data ownership |

## Global Flags

All commands support the following global flags:

| Flag | Description |
|------|-------------|
| `--domain <domain>` | Google Workspace domain |
| `--client-secret <path>` | Path to OAuth2 client secret file |
| `--cache-file <path>` | Path to token cache file |
| `-h, --help` | Show help for command |

## Getting Detailed Help

For detailed help on any command, use the `--help` flag:

```bash
# Get help on a specific command
gac user create --help

# Get help on a command group
gac user --help

# Get help on subcommands
gac group-settings update --help
```

## Related Documentation

- [User Management Guide](../guides/user-management.md)
- [Group Settings Guide](../guides/group-settings.md)
- [Troubleshooting](troubleshooting.md)
- [Examples](../../examples/README.md)
