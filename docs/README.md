# Documentation

Welcome to the **gac** (Google Admin Client) documentation! This guide will help you get started and make the most of this powerful CLI tool for managing Google Workspace.

## 📚 Table of Contents

### Getting Started
- [Installation](installation.md) - Installing gac on your system
- [Configuration](configuration.md) - Setting up credentials and config files
- [Authentication](authentication.md) - OAuth2 setup and credential management

### User Guides
- [User Management](guides/user-management.md) - Create, update, list, suspend users
- [Group Management](guides/group-management.md) - Manage groups and memberships
- [Group Settings](guides/group-settings.md) - Configure group access, posting, moderation
- [Organizational Units](guides/ou-management.md) - Manage organizational structure
- [Alias Management](guides/alias-management.md) - Email aliases for users
- [Calendar Operations](guides/calendar-operations.md) - Create and manage calendar events
- [Calendar Resources](guides/calendar-resources.md) - Manage rooms and equipment

### Reference
- [Command Reference](reference/commands.md) - Complete command documentation
- [Troubleshooting](reference/troubleshooting.md) - Common issues and solutions

### Development
- [Architecture](development/architecture.md) - System design and components
- [Contributing](development/contributing.md) - How to contribute to gac
- [Debugging](development/debugging.md) - Debug and profile gac
- [Releasing](development/releasing.md) - Release process and versioning

## 🚀 Quick Links

### First Time Users
1. **[Install gac](installation.md)** - Get the CLI installed
2. **[Set up OAuth](authentication.md)** - Configure Google Cloud credentials
3. **[Run your first command](../README.md#quick-start)** - List users in your domain

### Common Tasks
- **[Create a user](guides/user-management.md#create-a-user)** - Onboard new employees
- **[Manage group settings](guides/group-settings.md)** - Configure group permissions
- **[List calendar resources](guides/calendar-resources.md#list-resources)** - Find available rooms

### Developers
- **[Architecture overview](development/architecture.md)** - Understand the codebase
- **[Contributing guide](development/contributing.md)** - Submit your first PR
- **[Running tests](development/contributing.md#testing)** - Ensure code quality

## 📖 Documentation Organization

```
docs/
├── README.md                    # You are here
├── installation.md              # Installation instructions
├── configuration.md             # Configuration guide
├── authentication.md            # OAuth and credentials
│
├── guides/                      # User guides by feature
│   ├── user-management.md
│   ├── group-management.md
│   ├── group-settings.md
│   ├── ou-management.md
│   ├── alias-management.md
│   ├── calendar-operations.md
│   └── calendar-resources.md
│
├── reference/                   # Reference documentation
│   ├── commands.md
│   └── troubleshooting.md
│
└── development/                 # Developer documentation
    ├── architecture.md
    ├── contributing.md
    ├── debugging.md
    └── releasing.md
```

## 🔍 Finding What You Need

### By Task
- **Setting up for the first time?** → [Installation](installation.md) & [Authentication](authentication.md)
- **Managing users?** → [User Management Guide](guides/user-management.md)
- **Configuring groups?** → [Group Settings Guide](guides/group-settings.md)
- **Something not working?** → [Troubleshooting](reference/troubleshooting.md)
- **Want to contribute?** → [Contributing Guide](development/contributing.md)

### By Role
- **Administrators** - Start with [User Management](guides/user-management.md) and [Group Management](guides/group-management.md)
- **Developers** - Check out [Architecture](development/architecture.md) and [Contributing](development/contributing.md)
- **New Users** - Begin with [Installation](installation.md) and [Quick Start](../README.md#quick-start)

## 💡 Additional Resources

- **[Examples](../examples/README.md)** - Real-world scenarios with runnable scripts
- **[Command Reference](reference/commands.md)** - Complete list of all commands
- **[Changelog](../CHANGELOG.md)** - What's new in each version
- **[GitHub Issues](https://github.com/acockrell/google-admin-client/issues)** - Report bugs or request features

## 🤝 Getting Help

- **Documentation issues?** Open an issue on [GitHub](https://github.com/acockrell/google-admin-client/issues)
- **Found a bug?** See [Troubleshooting](reference/troubleshooting.md) first, then file a bug report
- **Want to contribute?** Read the [Contributing Guide](development/contributing.md)

---

**Need to see the code?** Check out the [main README](../README.md) for project overview and quick start.
