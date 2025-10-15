# gac - Google Admin Client

A powerful command-line tool for managing Google Workspace users, groups, calendars, and resources.

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Overview

**gac** (Google Admin Client) is a CLI tool for automating Google Workspace administrative tasks. Built with Go, it provides a simple interface for managing users, groups, calendars, and more through the Google Admin SDK APIs.

### Key Features

- **User Management** - Create, list, update, suspend users
- **Group Management** - Manage groups, memberships, and settings
- **Calendar Operations** - Create and manage calendar events
- **Calendar Resources** - Manage bookable resources (rooms, equipment)
- **Organizational Units** - Manage organizational structure
- **Alias Management** - Email aliases for users
- **Audit Log Export** - Export audit logs for compliance and analysis
- **Performance Caching** - Built-in caching for faster queries (30-90x speedup)
- **Shell Completion** - Bash, zsh, and fish completion support
- **Config Validation** - Validate configuration and credentials
- **Secure Authentication** - OAuth2 with automatic token refresh
- **Cross-Platform** - Linux and macOS (amd64 and arm64)

## Quick Start

### Installation

```bash
# macOS (Homebrew - coming soon)
# brew install acockrell/tap/gac

# Pre-built binaries
# Download from: https://github.com/acockrell/google-admin-client/releases

# macOS (Apple Silicon)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_darwin_arm64.tar.gz
tar xzf gac_darwin_arm64.tar.gz
sudo mv gac /usr/local/bin/

# Linux (amd64)
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_linux_amd64.tar.gz
tar xzf gac_linux_amd64.tar.gz
sudo mv gac /usr/local/bin/

# Build from source
git clone https://github.com/acockrell/google-admin-client.git
cd google-admin-client
make build
sudo mv build/gac /usr/local/bin/
```

### Authentication

1. **Set up OAuth2 credentials** in [Google Cloud Console](https://console.cloud.google.com)
   - Create a project
   - Enable Admin SDK API and Calendar API
   - Create OAuth2 credentials (Desktop app)
   - Download the JSON file

2. **Configure credentials**:
   ```bash
   mkdir -p ~/.credentials
   mv ~/Downloads/client_secret_*.json ~/.credentials/client_secret.json
   chmod 600 ~/.credentials/client_secret.json
   ```

3. **Authenticate**:
   ```bash
   gac user list
   # Follow the browser prompts to grant permissions
   ```

📖 **Detailed setup**: See [Authentication Guide](docs/authentication.md)

### Basic Usage

```bash
# List all users
gac user list

# Create a new user
gac user create john.doe@example.com

# Update user department
gac user update --dept Engineering john.doe@example.com

# View group settings
gac group-settings list team@example.com

# List calendar resources
gac cal-resource list --type room
```

## Configuration

Configure `gac` via config file, environment variables, or CLI flags:

### Config File (~/.google-admin.yaml)

```yaml
domain: example.com
client-secret: /path/to/client_secret.json
cache-file: /path/to/token.json
```

### Environment Variables

```bash
export GAC_DOMAIN=example.com
export GAC_CLIENT_SECRET=~/.credentials/client_secret.json
export GAC_CACHE_FILE=~/.credentials/gac.json
```

### CLI Flags

```bash
gac --domain example.com user list
```

📖 **Full configuration guide**: [docs/configuration.md](docs/configuration.md)

### Logging and Debugging

Control log output for troubleshooting and monitoring:

```bash
# Enable verbose/debug logging
gac --verbose user list
gac -v user list

# Set specific log level (debug, info, warn, error)
gac --log-level debug user list

# JSON log output for automation/parsing
gac --json-log --log-level debug user list > debug.log

# Combine with other flags
gac -v --domain example.com user suspend user@example.com
```

**Log Levels:**
- `debug` - Detailed API calls, requests, and responses
- `info` - General operational messages (default)
- `warn` - Warning messages (insecure permissions, deprecations)
- `error` - Error messages only

**Use Cases:**
- **Troubleshooting** - Use `-v` or `--log-level debug` to see API calls and diagnose issues
- **Automation** - Use `--json-log` for structured logs that can be parsed by log aggregators
- **Production** - Use `--log-level warn` or `--log-level error` to reduce noise

## Common Tasks

### User Management

```bash
# Create user with groups
gac user create \
  --first-name John \
  --last-name Doe \
  --groups engineering \
  --groups all-staff \
  john.doe@example.com

# Suspend user
gac user suspend user@example.com --reason "Left company"

# Unsuspend user
gac user unsuspend user@example.com
```

📖 **Full guide**: [User Management](docs/guides/user-management.md)

### Group Settings

```bash
# Configure moderated announcements group
gac group-settings update announcements@example.com \
  --who-can-post-message ALL_MANAGERS_CAN_POST \
  --message-moderation-level MODERATE_ALL_MESSAGES \
  --allow-external-members false

# Add custom footer
gac group-settings update support@example.com \
  --custom-footer-text "For help, contact support@example.com" \
  --include-custom-footer true
```

📖 **Full guide**: [Group Settings](docs/guides/group-settings.md)

### Calendar Resources

```bash
# Create conference room
gac cal-resource create conf-room-a \
  --name "Conference Room A" \
  --type room \
  --capacity 12 \
  --building-id main-building

# List all rooms
gac cal-resource list --type room
```

📖 **Full guide**: [Calendar Resources](docs/guides/calendar-resources.md)

### Organizational Units

```bash
# Create OU
gac ou create /Engineering/Backend

# List OUs
gac ou list

# Move user to OU
gac user update --ou /Engineering/Backend user@example.com
```

📖 **Full guide**: [Organizational Units](docs/guides/ou-management.md)

### Audit Log Export

```bash
# Export last 24h of admin console activities
gac audit export --app admin

# Export login activities for specific user
gac audit export --app login --user user@example.com

# Export drive activities to CSV
gac audit export --app drive \
  --start-time 2024-10-01T00:00:00Z \
  --end-time 2024-10-08T00:00:00Z \
  --output csv --output-file drive-audit.csv

# Filter by event types
gac audit export --app admin --event-name USER_CREATED
```

📖 **Full guide**: [Audit Logs](docs/guides/audit-logs.md)

### Performance and Caching

Speed up repeated queries with built-in caching:

```bash
# First call - fetches from API (~1200ms)
gac user list

# Subsequent calls - uses cache (~35ms) - 34x faster!
gac user list

# View cache statistics
gac cache status

# Clear cache when needed
gac cache clear users
gac cache clear groups
gac cache clear all

# Disable cache for fresh data
gac user list --no-cache

# Configure cache TTL
gac user list --cache-ttl 30m
```

**Benefits:**
- **30-90x faster** queries with caching enabled
- **80-90% reduction** in API quota usage
- Automatic cache expiration (default: 15 minutes)
- Configurable cache location and TTL

**Configuration** (`~/.google-admin.yaml`):
```yaml
cache:
  enabled: true
  ttl: 15m
  directory: ~/.cache/gac
```

📖 **Full guide**: [Caching](docs/guides/caching.md)

### CLI Utilities

```bash
# Show version information
gac version

# Show version number only
gac version --short

# Validate configuration
gac config validate

# Generate shell completion
gac completion bash > /etc/bash_completion.d/gac  # Linux
gac completion zsh > ~/.oh-my-zsh/completions/_gac  # zsh
gac completion fish > ~/.config/fish/completions/gac.fish  # fish

# Skip confirmations for automation (use with caution)
gac user suspend user@example.com --yes
gac ou delete /OldOU -y
```

📖 **Shell completion guide**: [Shell Completion](docs/guides/shell-completion.md)

## Documentation

### 📚 User Guides
- [User Management](docs/guides/user-management.md) - Create, update, suspend users
- [Group Management](docs/guides/group-management.md) - Manage groups and memberships
- [Group Settings](docs/guides/group-settings.md) - Configure group permissions and behavior
- [Calendar Operations](docs/guides/calendar-operations.md) - Manage calendar events
- [Calendar Resources](docs/guides/calendar-resources.md) - Manage rooms and equipment
- [Organizational Units](docs/guides/ou-management.md) - Manage organizational structure
- [Alias Management](docs/guides/alias-management.md) - Email aliases for users
- [Audit Logs](docs/guides/audit-logs.md) - Export audit logs for compliance and analysis
- [Shell Completion](docs/guides/shell-completion.md) - Set up tab completion for your shell

### 📖 Reference
- [Command Reference](docs/reference/commands.md) - Complete command list
- [Troubleshooting](docs/reference/troubleshooting.md) - Common issues and solutions

### 🔧 Configuration & Setup
- [Installation](docs/installation.md) - Detailed installation instructions
- [Authentication](docs/authentication.md) - OAuth2 setup and security
- [Configuration](docs/configuration.md) - Configuration options

### 💻 Development
- [Contributing](docs/development/contributing.md) - How to contribute
- [Architecture](docs/development/architecture.md) - Technical design
- [Debugging](docs/development/debugging.md) - Debug and profile gac
- [Releasing](docs/development/releasing.md) - Release process

### 📝 Examples
- [Examples Directory](examples/) - Runnable scripts for common scenarios
- [Examples Guide](examples/README.md) - Detailed scenario walkthroughs

## Command Reference

| Category | Commands |
|----------|----------|
| **Users** | `create`, `list`, `update`, `suspend`, `unsuspend` |
| **Groups** | `list` |
| **Group Settings** | `list`, `update` |
| **Calendar** | `create`, `list`, `update` |
| **Calendar Resources** | `list`, `create`, `update`, `delete` |
| **Organizational Units** | `list`, `create`, `update`, `delete` |
| **Aliases** | `list`, `add`, `remove` |
| **Audit** | `export` |
| **Data Transfer** | `transfer` |

📖 **Full command reference**: [docs/reference/commands.md](docs/reference/commands.md)

## Troubleshooting

### Common Issues

**Authentication errors**
```bash
# Delete cached token and re-authenticate
rm ~/.credentials/gac.json
gac user list
```

**Permission errors**
- Verify you have Google Workspace admin privileges
- Check OAuth scopes in [docs/authentication.md](docs/authentication.md)
- Ensure required APIs are enabled in Google Cloud Console

**Domain configuration**
```bash
# Set domain via environment variable
export GAC_DOMAIN=example.com
gac user list
```

📖 **Full troubleshooting guide**: [docs/reference/troubleshooting.md](docs/reference/troubleshooting.md)

## Contributing

We welcome contributions! See [CONTRIBUTING.md](docs/development/contributing.md) for:
- Reporting bugs
- Suggesting features
- Submitting pull requests
- Development setup

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Google Admin SDK](https://developers.google.com/admin-sdk) - Google Workspace APIs

## Links

- **Documentation**: [docs/](docs/)
- **Examples**: [examples/](examples/)
- **Issues**: [GitHub Issues](https://github.com/acockrell/google-admin-client/issues)
- **Releases**: [GitHub Releases](https://github.com/acockrell/google-admin-client/releases)
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)

---

Made with ❤️ for Google Workspace administrators
