# TODO: Improvements for gac

## Critical Updates

<details>
<summary>✅ 1. Upgrade Go Version</summary>

- [x] Update go.mod to Go 1.25
- [x] Update Dockerfile to use golang:1.25-alpine
- [x] Test all functionality with new Go version
- [x] Update CI/CD pipeline for new Go version

**Status:** ✅ Complete
- go.mod updated to 1.25
- Dockerfile updated to golang:1.25.1-alpine (official Docker Hub repository)
- GitHub Actions workflow (.github/workflows/gosec.yml) updated to Go 1.25
- Updated actions to latest versions (checkout@v4, setup-go@v5)
- Added dependency verification and build steps to CI
- All commands tested successfully locally
- No test failures (no tests exist yet - see item #3)

</details>

<details>
<summary>✅ 2. Update Dependencies</summary>

- [x] Run `go get -u ./...` to update all dependencies
- [x] Update Google Cloud libraries (v0.103 → v0.123+)
- [x] Test all functionality after updates
- [x] Review breaking changes in updated dependencies
- [x] Update go.sum with new checksums

**Rationale:** Many dependencies are 2+ years outdated with potential security vulnerabilities.

**Status:** ✅ Complete
Major dependency updates:
- **google.golang.org/api**: v0.92.0 → v0.251.0 (major version jump!)
- **github.com/spf13/cobra**: v1.5.0 → v1.10.1
- **github.com/spf13/viper**: v1.12.0 → v1.21.0
- **golang.org/x/oauth2**: v0.0.0-20220808 → v0.31.0
- **google.golang.org/grpc**: v1.48.0 → v1.75.1
- **google.golang.org/protobuf**: v1.28.1 → v1.36.9
- **cloud.google.com/go/compute**: v1.8.0 → v1.38.0
- All golang.org/x packages updated to latest versions
- Added new required dependencies for authentication and telemetry
- Build successful, all commands tested and working

</details>

<details>
<summary>🚧 3. Add Testing</summary>

- [x] Add unit tests for helper functions (parsePhone, parseAddress, etc.)
- [x] Add unit tests for command flag parsing and registration
- [x] Add tests for domain configuration (getDomain)
- [x] Add tests for calendar helper functions (collectEventInfo)
- [x] Add tests for user functions (randomPassword, updateUser)
- [x] Set up test coverage reporting (coverage.out, coverage.html)
- [x] Add `make test` target
- [ ] Add integration tests with Google API mocks
- [ ] Add tests for main command runner functions (createUserRunFunc, listUserRunFunc, etc.)
- [ ] Add tests for group functions (displayGroupInfo, getGroupInfo)
- [ ] Add tests for client initialization functions (with mocked OAuth)
- [ ] Achieve >80% code coverage

**Rationale:** No test files exist currently. Testing is critical for reliability and maintenance.

**Status:** 🚧 In Progress (23.6% coverage achieved)
Test files created:
- `cmd/user_test.go`: Tests for randomPassword, updateUser
- `cmd/root_test.go`: Tests for getDomain and root command flags
- `cmd/group_test.go`: Tests for group email construction
- `cmd/commands_test.go`: Tests for command registration and flags
- `cmd/user_update_test.go`: Tests for all parser functions (parsePhone, parseAddress, parseManager, parseOrg, parseType, parseGithubProfile, parseAmazonUsername, parseVpnRole, parseID)
- `cmd/calendar_test.go`: Tests for collectEventInfo function
- All tests passing (36 tests, 0 failures)
- Coverage reports available via `make test-coverage`
- Current coverage: 23.6% (started at 0%)

**Next Steps to Reach >80% Coverage:**
1. Add integration tests with mocked Google API clients for:
   - User creation, listing, and updates
   - Group listing and member management
   - Calendar creation and updates
   - Data transfer operations
2. Add tests for command runner functions (currently 0% coverage)
3. Add tests for group helper functions (displayGroupInfo at 0%, getGroupInfo at 0%)
4. Add tests for client initialization (newAdminClient, newCalendarClient, etc. currently at 0%)
5. Consider testing error handling paths and edge cases

</details>

## Security Improvements

<details>
<summary>✅ 4. Credential Management</summary>

- [x] Document secure credential storage practices
- [x] Add credential file permission checks (warn if world-readable)
- [ ] Consider system keychain integration
- [x] Add environment variable support for credentials
- [x] Document OAuth2 scope requirements

**Rationale:** Improve security posture for credential handling.

**Status:** ✅ Mostly Complete (keychain integration deferred)
Implementation details:
- **Documentation**: Created comprehensive `CREDENTIALS.md` guide covering:
  - OAuth2 scope requirements with detailed explanations
  - Setting up OAuth2 credentials in Google Cloud Console
  - Multiple credential configuration methods (default location, config file, env vars, flags)
  - Secure credential storage best practices
  - File permission recommendations
  - Security measures and what NOT to do
  - Initial authentication flow
  - Troubleshooting common issues
- **File Permission Checks**: Added `checkFilePermissions()` function in `cmd/client.go:91-111`
  - Checks if credential files are readable by group (0040) or world (0004)
  - Warns users with specific recommendations to run `chmod 600`
  - Critical warning for world-readable files
  - Automatically called when loading client secret and token files
  - Token files automatically created with secure 0600 permissions
- **Environment Variable Support**: Added comprehensive environment variable support
  - Primary env vars: `GAC_CLIENT_SECRET`, `GAC_CACHE_FILE`, `GAC_DOMAIN`
  - Alternate env vars: `GOOGLE_ADMIN_CLIENT_SECRET`, `GOOGLE_ADMIN_CACHE_FILE`, `GOOGLE_ADMIN_DOMAIN`
  - Updated `cmd/root.go:64-72` with viper environment variable bindings
  - Updated `cmd/client.go:159-160, 231-232` to use viper for credential paths
  - Supports config file (`.google-admin.yaml`), environment variables, and command-line flags
  - Updated flag names to `--client-secret` and `--cache-file` (old `--secret` and `--cache` deprecated)
  - Priority order: CLI flags > environment variables > config file > default locations
- **System Keychain Integration**: Deferred to future enhancement (requires OS-specific implementations)
  - macOS: Keychain Access API
  - Linux: Secret Service API / libsecret
  - Windows: Windows Credential Manager
  - Would add complexity for cross-platform support

</details>

<details>
<summary>✅ 5. Input Validation</summary>

- [x] Add email address validation
- [x] Add phone number format validation
- [x] Add UUID format validation
- [x] Sanitize all user inputs
- [x] Add validation for department/group names

**Rationale:** Prevent invalid data and potential injection issues.

**Status:** ✅ Complete
Implementation details:
- Created `cmd/validation.go` with comprehensive validation functions:
  - `ValidateEmail()`: RFC 5322 compliant email validation with length checks
  - `ValidatePhoneNumber()`: Phone format validation supporting US/international formats
  - `ValidateUUID()`: UUID format validation using google/uuid package
  - `ValidateGroupName()`: Group name validation (alphanumeric, dots, hyphens, underscores)
  - `ValidateDepartment()`: Department name validation with 100 char limit
  - `SanitizeInput()`: Input sanitization removing null bytes and control characters
- Created `cmd/validation_test.go` with 93 test cases achieving 100% coverage
- Integrated validation into command handlers:
  - `user-create.go`: Email and group name validation
  - `user-update.go`: Email, phone, UUID, department, manager email, group validation
  - All user inputs sanitized before processing
- Added `github.com/google/uuid v1.6.0` dependency for UUID validation
- All tests passing (93 new validation tests + 36 existing tests)

</details>

<details>
<summary>✅ 6. Fix gosec Security Issues</summary>

- [x] **HIGH**: Fix weak RNG in password generation (`cmd/user.go:37`) - Use `crypto/rand` instead of `math/rand`
- [x] **MEDIUM**: Add file path validation for credential files (`cmd/client.go:92, 161, 182`) - Prevent directory traversal
- [x] **LOW**: Handle errors from `csv.Writer.Write()` (`cmd/user-list.go:120-125`)
- [x] **LOW**: Handle errors from `viper.BindPFlag()` (`cmd/root.go:40`)

**Rationale:** Fix security vulnerabilities identified by gosec scanner.

**Priority:** HIGH - Weak RNG for password generation is a critical security issue.

**Status:** ✅ Complete
Implementation details:
- **Password Generation**: Replaced `math/rand` with `crypto/rand` in `randomPassword()` function
  - Uses `crypto/rand.Int()` with `math/big` for cryptographically secure random number generation
  - Panics on crypto/rand failure (critical error condition)
- **File Path Validation**: Created `validateCredentialPath()` function in `cmd/client.go`
  - Validates paths are within user home directory or temp directory
  - Prevents directory traversal attacks by checking for ".." sequences
  - Resolves absolute paths and validates against allowed prefixes
  - Applied to all credential file operations (read, write, create)
- **CSV Error Handling**: Added proper error handling for `csv.Writer.Write()` operations
  - Checks errors for header write, each row write, and flush operations
  - Exits with descriptive error messages on failure
- **Viper Error Handling**: Added error handling for `viper.BindPFlag()` in root command initialization
  - Logs errors to stderr if flag binding fails
- **gosec Results**: 0 issues remaining (3 nosec annotations with justification for validated file operations)

</details>

## Code Quality

<details>
<summary>✅ 7. Remove Hardcoded Domain</summary>

- [x] Add `domain` field to configuration file (.google-admin.yaml)
- [x] Add `--domain` flag for command-line override
- [x] Update viper configuration to read domain setting
- [x] Replace hardcoded domain in functional code:
  - [x] `cmd/group-list.go:110` - Remove email filtering or make configurable
  - [x] `cmd/group-list.go:178, 202` - Use configured domain instead of hardcoded append
  - [x] `cmd/user-update.go:236` - Use configured domain for group insertion
  - [x] `cmd/user-create.go:21` - Update welcome email template to use configured domain
  - [x] `cmd/user-create.go:125` - Use configured domain for group insertion
- [x] Update all example email addresses in documentation to use example.com
- [x] Support default domain fallback for backward compatibility

**Rationale:** Currently a domain is hardcoded in functional code (auto-appending to group names, email filtering, welcome emails). This makes the tool non-reusable for other organizations. Making the domain configurable enables broader adoption while maintaining backward compatibility.

**Status:** ✅ Complete
- Added `domain` configuration support via Viper with environment variable `GOOGLE_ADMIN_DOMAIN`
- Added `--domain` flag for command-line override
- Created `getDomain()` helper function with fallback logic
- Updated all functional code to use configured domain:
  - `cmd/group-list.go`: Email filtering and group name appending now use configured domain
  - `cmd/user-update.go`: Group insertion uses configured domain
  - `cmd/user-create.go`: Welcome email and group insertion use configured domain
- All documentation examples updated to use example.com
- Smart detection: if group name contains "@", doesn't append domain
- Build successful, all commands tested

</details>

<details>
<summary>📋 8. Error Handling</summary>

- [ ] Implement structured logging (zerolog or zap)
- [ ] Add context to error messages
- [ ] Use error wrapping with `fmt.Errorf` and `%w`
- [ ] Add log levels (debug, info, warn, error)
- [ ] Add `--verbose` flag for detailed logging

**Rationale:** Better debugging and error tracking.

</details>

<details>
<summary>✅ 9. CI/CD Modernization</summary>

- [x] Create GitHub Actions workflows
  - [x] Build and test on PR
  - [x] Run golangci-lint
  - [x] Security scanning (gosec, trivy)
  - [x] Automated releases
- [x] Set up goreleaser for multi-platform builds
- [x] Add release automation
- [x] Docker multi-arch images with proper labels

**Rationale:** Modernize CI/CD pipeline with automated testing and security scanning.

**Status:** ✅ Complete
Implementation details:
- **CI Workflow** (`.github/workflows/ci.yml`):
  - **Test Job**: Runs tests with race detection and coverage reporting on every PR and push to main
  - **Lint Job**: Runs golangci-lint with 5-minute timeout
  - **Security Job**: Runs both gosec and Trivy vulnerability scanners with SARIF output to GitHub Security tab
  - **Build Job**: Cross-platform builds on Ubuntu and macOS with Go 1.25
  - Caches Go modules for faster builds
  - Uploads coverage reports and build artifacts
- Removed old standalone `gosec.yml` workflow (now integrated into ci.yml)
- **Release Automation** (`.github/workflows/release.yml`):
  - Triggered automatically on version tags (e.g., `v0.1.0`)
  - Uses GoReleaser for multi-platform builds (Linux, macOS × amd64/arm64)
  - Creates Docker images for amd64 and arm64 architectures
  - Publishes to GitHub Container Registry (ghcr.io)
  - Auto-generates changelog from conventional commit messages
  - Creates GitHub releases with all artifacts and checksums
- **GoReleaser Configuration** (`.goreleaser.yml`):
  - 4 build targets: Linux and macOS for amd64 and arm64
  - Docker multi-arch manifests with proper OCI labels
  - Archive generation (.tar.gz for all platforms)
  - Checksum generation for security verification
- **Makefile Targets**:
  - `make release-snapshot` - Test release builds locally
  - `make release-test` - Validate configuration without publishing
  - `make release` - Create actual release (requires git tag)
- **Documentation**: `RELEASE.md` with comprehensive release process guide
- **PR #7**: https://github.com/acockrell/google-admin-client/pull/7

</details>

<details>
<summary>✅ 10. Documentation</summary>

- [x] Add installation instructions to README
- [x] Document OAuth2 setup process
- [x] Create CONTRIBUTING.md
- [x] Document configuration file format (.google-admin.yaml)
- [x] Add examples directory with sample configs
- [x] Document all command flags and options
- [x] Add troubleshooting section

**Rationale:** Improve onboarding and usability.

**Status:** ✅ Complete
Implementation details:
- **README.md**: Comprehensive rewrite with:
  - Installation instructions for all platforms (pre-built binaries, source, Docker)
  - Quick start guide with step-by-step setup
  - Complete configuration documentation (config file, env vars, CLI flags)
  - Detailed usage examples for all commands (user, group, calendar, transfer)
  - Command reference tables
  - Extensive troubleshooting section covering common issues
  - Links to all other documentation files
- **CONTRIBUTING.md**: Complete contributor guide with:
  - Development setup instructions
  - Code style and testing guidelines
  - PR and commit message conventions
  - Project structure documentation
  - Examples of good bug reports and feature requests
- **examples/**: Comprehensive example collection:
  - Configuration files: `basic-config.yaml`, `production-config.yaml`, `development-config.yaml`
  - Scripts: `onboarding-example.sh`, `offboarding-example.sh`, `bulk-export.sh`, `group-audit.sh`, `create-recurring-meeting.sh`, `batch-update-users.sh`
  - Sample CSV: `users-to-update.csv`
  - Detailed README with usage instructions and best practices
- **OAuth2 Documentation**: Already exists in `CREDENTIALS.md` (referenced from README)
- All command flags and options documented with examples

</details>

## Feature Enhancements

<details>
<summary>📋 11. Output Formats</summary>

- [ ] Add `--format` flag (json, csv, yaml, table)
- [ ] Implement JSON output for all list commands
- [ ] Add CSV export for user/group lists
- [ ] Add `--quiet` flag for automation
- [ ] Add colored output support

**Rationale:** Better integration with automation scripts.

</details>

<details>
<summary>📋 12. Batch Operations</summary>

- [ ] Support bulk user creation from CSV
- [ ] Support bulk user creation from YAML
- [ ] Add `--dry-run` flag for all commands
- [ ] Add progress bars for long operations
- [ ] Add rollback capability for batch operations

**Rationale:** Improve efficiency for large-scale operations.

</details>

<details>
<summary>📋 13. Modern CLI Features</summary>

- [ ] Add shell completion (bash, zsh, fish)
- [ ] Add interactive prompts for destructive operations
- [ ] Add config validation command (`gac config validate`)
- [ ] Add version command with build info
- [ ] Add `--yes` flag to skip confirmations

**Rationale:** Improve user experience and safety.

</details>

<details>
<summary>📋 14. Performance</summary>

- [ ] Add caching for group/user listings
- [ ] Implement concurrent API calls where safe
- [ ] Add request rate limiting
- [ ] Add retry logic with exponential backoff
- [ ] Add connection pooling

**Rationale:** Improve performance and handle API quotas gracefully.

</details>

## Nice to Have

<details>
<summary>🚧 15. Additional Features</summary>

- [x] Add user suspension/unsuspension commands
- [x] Add organizational unit management
- [x] Add alias management for users
- [ ] Add calendar resource management
- [ ] Add group settings management
- [ ] Add audit log export

**Status:** 🚧 In Progress (3/6 features complete)

**Organizational Unit Management** - ✅ Complete
Implementation details:
- **`gac ou list`**: List all organizational units or specific OU with children
  - Supports `--type` flag (all/children) for filtering
  - Hierarchical display with indentation based on depth
  - Shows path, description, parent, ID, and inheritance settings
- **`gac ou create`**: Create new organizational units
  - Validates OU path format (must start with /)
  - Auto-detects parent from path or accepts `--parent` flag
  - Supports `--description` and `--block-inheritance` flags
  - Provides clear examples for top-level and nested OUs
- **`gac ou update`**: Update existing organizational units
  - Update name, description, parent (move OU), or block inheritance
  - Supports partial updates (only specified fields changed)
  - Warns about impacts of moving OUs
- **`gac ou delete`**: Delete organizational units
  - Confirmation prompt (skippable with `--force`)
  - OU must be empty (no users or sub-OUs)
  - Provides helpful error messages for common failures
- **Tests**: Added comprehensive test suite in `cmd/ou_test.go`
  - Tests for command existence and structure
  - Validates all flags are present
  - All 140 tests passing
- **Files Created**:
  - `cmd/ou.go` - Root OU command
  - `cmd/ou-list.go` - List organizational units
  - `cmd/ou-create.go` - Create organizational units
  - `cmd/ou-update.go` - Update organizational units
  - `cmd/ou-delete.go` - Delete organizational units
  - `cmd/ou_test.go` - Test suite

**User Alias Management** - ✅ Complete
Implementation details:
- **`gac alias list`**: List all email aliases for a user
  - Displays primary user email, all aliases, and total count
  - Handles Google API's []interface{} response with type assertions
  - Email validation for user addresses
- **`gac alias add`**: Add email alias to a user
  - Validates both user email and alias email formats
  - Checks alias is in managed domains
  - Provides clear error messages for conflicts
  - Shows requirements and common use cases
- **`gac alias remove`**: Remove email alias from a user
  - Confirmation prompt (skippable with `--force` flag)
  - Email validation for both addresses
  - Warning about mail delivery impact
  - Helpful error messages for common failures
- **Tests**: Added comprehensive test suite in `cmd/alias_test.go`
  - Tests for command existence and structure
  - Validates all flags are present
  - All 146 tests passing
- **Documentation**:
  - README.md updated with alias management section and examples
  - examples/alias-management.sh demo script
  - examples/README.md updated with alias example (Section 8)
- **Files Created**:
  - `cmd/alias.go` - Root alias command
  - `cmd/alias-list.go` - List user aliases
  - `cmd/alias-add.go` - Add user alias
  - `cmd/alias-remove.go` - Remove user alias
  - `cmd/alias_test.go` - Test suite
  - `examples/alias-management.sh` - Demo script
- **Common Use Cases**:
  - Department addresses (support@, sales@, info@)
  - Role-based addresses (admin@, webmaster@)
  - Alternative name formats (first.last@, firstlast@)
  - Legacy addresses when renaming users

**User Suspension/Unsuspension** - ✅ Complete
Implementation details:
- **`gac user suspend`**: Suspend user accounts
  - Optional suspension reason tracking (`--reason` flag)
  - Confirmation prompt (skippable with `--force` flag)
  - Email validation for user addresses
  - Prevents users from signing in and accessing services
  - Preserves all account data for restoration
  - ForceSendFields to ensure boolean values are sent to API
- **`gac user unsuspend`**: Unsuspend (restore) user accounts
  - Restores full account access
  - Confirmation prompt (skippable with `--force` flag)
  - Clear communication of restored capabilities
  - Email validation for user addresses
- **Tests**: Added comprehensive test suite in `cmd/user_suspend_test.go`
  - Tests for command existence and structure
  - Validates all flags are present
  - Tests for command registration
  - All 152 tests passing
- **Documentation**:
  - README.md updated with suspension commands, examples, and use cases
  - examples/user-suspension.sh demo script
  - examples/README.md updated with suspension example (Section 9)
- **Files Created**:
  - `cmd/user-suspend.go` - Suspend user accounts
  - `cmd/user-unsuspend.go` - Unsuspend user accounts
  - `cmd/user_suspend_test.go` - Test suite
  - `examples/user-suspension.sh` - Demo script
- **Use Cases**:
  - Employee termination or departure
  - Policy violations or security incidents
  - Account compromise or suspicious activity
  - Extended leave or sabbatical
- **Workflows Documented**:
  - Employee departure workflow
  - Security incident workflow
  - Extended leave workflow
  - Best practices for suspension management

</details>

<details>
<summary>✅ 16. Developer Experience</summary>

- [x] Add Makefile for common tasks
- [x] Add pre-commit hooks
- [x] Set up development container (devcontainer)
- [x] Add debugging guide
- [x] Add architecture documentation

**Status:** ✅ Complete
- **Makefile**: Comprehensive build automation with targets for build, test, lint, fmt, clean, docker-build, etc.
- **Pre-commit hooks**: Both `.pre-commit-config.yaml` for pre-commit framework and `.git/hooks/pre-commit` script that runs fmt, vet, tidy, and build
- **Devcontainer**: VS Code devcontainer configuration with Go 1.25, extensions, and credential mounting
- **DEBUGGING.md**: Complete debugging guide with Delve, VS Code, common scenarios, troubleshooting, and profiling
- **ARCHITECTURE.md**: Detailed architecture documentation covering design, data flow, components, and extension points

</details>
