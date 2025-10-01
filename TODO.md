# TODO: Improvements for gac

## Critical Updates

### 1. Upgrade Go Version
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

### 2. Update Dependencies
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

### 3. Add Testing
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

## Security Improvements

### 4. Credential Management
- [ ] Document secure credential storage practices
- [ ] Add credential file permission checks (warn if world-readable)
- [ ] Consider system keychain integration
- [ ] Add environment variable support for credentials
- [ ] Document OAuth2 scope requirements

**Rationale:** Improve security posture for credential handling.

### 5. Input Validation
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

### 6. Fix gosec Security Issues
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

## Code Quality

### 7. Remove Hardcoded Domain
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

### 8. Error Handling
- [ ] Implement structured logging (zerolog or zap)
- [ ] Add context to error messages
- [ ] Use error wrapping with `fmt.Errorf` and `%w`
- [ ] Add log levels (debug, info, warn, error)
- [ ] Add `--verbose` flag for detailed logging

**Rationale:** Better debugging and error tracking.

### 9. CI/CD Modernization
- [ ] Create GitHub Actions workflows
  - [ ] Build and test on PR
  - [ ] Run golangci-lint
  - [ ] Security scanning (gosec, trivy)
  - [ ] Automated releases
- [ ] Set up goreleaser for multi-platform builds
- [ ] Add release automation
- [ ] Migrate from Jenkins or run both in parallel

**Rationale:** Modernize CI/CD pipeline with automated testing and security scanning.

### 10. Documentation
- [ ] Add installation instructions to README
- [ ] Document OAuth2 setup process
- [ ] Create CONTRIBUTING.md
- [ ] Document configuration file format (.google-admin.yaml)
- [ ] Add examples directory with sample configs
- [ ] Document all command flags and options
- [ ] Add troubleshooting section

**Rationale:** Improve onboarding and usability.

## Feature Enhancements

### 11. Output Formats
- [ ] Add `--format` flag (json, csv, yaml, table)
- [ ] Implement JSON output for all list commands
- [ ] Add CSV export for user/group lists
- [ ] Add `--quiet` flag for automation
- [ ] Add colored output support

**Rationale:** Better integration with automation scripts.

### 12. Batch Operations
- [ ] Support bulk user creation from CSV
- [ ] Support bulk user creation from YAML
- [ ] Add `--dry-run` flag for all commands
- [ ] Add progress bars for long operations
- [ ] Add rollback capability for batch operations

**Rationale:** Improve efficiency for large-scale operations.

### 13. Modern CLI Features
- [ ] Add shell completion (bash, zsh, fish)
- [ ] Add interactive prompts for destructive operations
- [ ] Add config validation command (`gac config validate`)
- [ ] Add version command with build info
- [ ] Add `--yes` flag to skip confirmations

**Rationale:** Improve user experience and safety.

### 14. Performance
- [ ] Add caching for group/user listings
- [ ] Implement concurrent API calls where safe
- [ ] Add request rate limiting
- [ ] Add retry logic with exponential backoff
- [ ] Add connection pooling

**Rationale:** Improve performance and handle API quotas gracefully.

## Nice to Have

### 15. Additional Features
- [ ] Add user suspension/unsuspension commands
- [ ] Add organizational unit management
- [ ] Add alias management for users
- [ ] Add calendar resource management
- [ ] Add group settings management
- [ ] Add audit log export

### 16. Developer Experience
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
