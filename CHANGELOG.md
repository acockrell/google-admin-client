# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive documentation reorganization
  - Created `docs/` directory with organized structure
  - Added user guides for all major features
  - Added reference documentation (commands, troubleshooting)
  - Moved development docs to `docs/development/`
  - Slimmed down README.md from 966 to 290 lines
- Documentation hub at `docs/README.md`
- Detailed guides:
  - User Management
  - Group Settings Management
  - Command Reference
  - Troubleshooting Guide

## [0.3.0] - 2025-10-07

### Added
- Group settings management feature (#17)
  - `gac group-settings list` - View group settings in table or JSON format
  - `gac group-settings update` - Update group settings with 26+ configurable options
  - Comprehensive access control, posting permissions, moderation settings
  - Email customization (custom footers, reply-to)
  - Archive settings for read-only groups
  - Member management and invitation workflows
- Pre-commit hook improvements
  - Added golangci-lint to pre-commit checks
  - Created `hooks/` directory with versioned pre-commit hook
  - Added `hooks/README.md` with installation instructions

### Fixed
- JSON unmarshal error handling in group settings display
- Linter errors (errcheck) in group settings code

## [0.2.0] - 2025-10-06

### Added
- Calendar resource management (#16, #15)
  - `gac cal-resource list` - List calendar resources with type filtering
  - `gac cal-resource create` - Create conference rooms, equipment, etc.
  - `gac cal-resource update` - Update resource properties
  - `gac cal-resource delete` - Delete calendar resources
  - Support for rooms, equipment, and other resource types
  - Capacity management and feature tracking
- User suspension/unsuspension (#14)
  - `gac user suspend` - Suspend user accounts with optional reason
  - `gac user unsuspend` - Restore suspended accounts
  - Confirmation prompts with `--force` bypass
  - Suspension reason tracking
- User alias management (#12)
  - `gac alias list` - List email aliases for users
  - `gac alias add` - Add email aliases
  - `gac alias remove` - Remove aliases with confirmation
  - Email validation for aliases
- Organizational unit management (#11)
  - `gac ou list` - List organizational units
  - `gac ou create` - Create new OUs
  - `gac ou update` - Update OU properties
  - `gac ou delete` - Delete empty OUs
  - Hierarchical display with inheritance settings

### Changed
- Applied `go fmt` to all calendar resource files
- Improved error handling in OU commands
- Enhanced documentation for all new features

## [0.1.0] - Earlier

### Added
- Initial release
- User management commands
  - Create, list, update users
  - Comprehensive profile field support
- Group management
  - List groups and members
  - Group membership management
- Calendar operations
  - Create, list, update calendar events
  - Recurring event support
- Data transfer functionality
- OAuth2 authentication with token caching
- Configuration via file, environment variables, and CLI flags
- Input validation (emails, phone numbers, UUIDs)
- Cross-platform support (Linux, macOS - amd64, arm64)
- Security features
  - File permission checks for credentials
  - Path validation to prevent traversal attacks
  - Cryptographically secure password generation
- Comprehensive documentation
  - README with usage examples
  - OAuth2 setup guide (CREDENTIALS.md)
  - Architecture documentation
  - Contributing guidelines
  - Debugging guide
- CI/CD with GitHub Actions
  - Automated testing and linting
  - Security scanning (gosec, trivy)
  - Multi-platform builds
  - Automated releases with GoReleaser

## Release Notes

### How to Upgrade

```bash
# Download latest release
curl -LO https://github.com/acockrell/google-admin-client/releases/latest/download/gac_darwin_arm64.tar.gz
tar xzf gac_darwin_arm64.tar.gz
sudo mv gac /usr/local/bin/

# Or build from source
git pull
make build
sudo mv build/gac /usr/local/bin/
```

### Breaking Changes

None yet. This project follows semantic versioning.

### Deprecations

- Old flag names `--secret` and `--cache` are deprecated in favor of `--client-secret` and `--cache-file`
  - Old flags still work for backward compatibility but will be removed in v1.0.0

## Links

- [GitHub Releases](https://github.com/acockrell/google-admin-client/releases)
- [Documentation](docs/)
- [Contributing Guide](docs/development/contributing.md)
