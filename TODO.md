# TODO: Improvements for gac

## Critical Updates

<details>
<summary>🚧 1. Add Testing</summary>

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

## Feature Enhancements

<details>
<summary>📋 2. Output Formats</summary>

- [ ] Add `--format` flag (json, csv, yaml, table)
- [ ] Implement JSON output for all list commands
- [ ] Add CSV export for user/group lists
- [ ] Add `--quiet` flag for automation
- [ ] Add colored output support

**Rationale:** Better integration with automation scripts.

</details>

<details>
<summary>📋 3. Batch Operations</summary>

- [ ] Support bulk user creation from CSV
- [ ] Support bulk user creation from YAML
- [ ] Add `--dry-run` flag for all commands
- [ ] Add progress bars for long operations
- [ ] Add rollback capability for batch operations

**Rationale:** Improve efficiency for large-scale operations.

</details>

<details>
<summary>✅ 4. Modern CLI Features</summary>

- [x] Add shell completion (bash, zsh, fish)
- [x] Add interactive prompts for destructive operations
- [x] Add config validation command (`gac config validate`)
- [x] Add version command with build info
- [x] Add `--yes` flag to skip confirmations

**Rationale:** Improve user experience and safety.

**Status:** ✅ Complete

Implementation details:
- Created `cmd/version.go` with version, commit, date, and build info
- Updated `main.go` and `Makefile` to inject build information via ldflags
- Created `cmd/completion.go` with bash, zsh, and fish subcommands
- Created `cmd/prompts.go` with shared `confirmAction()` and `confirmDeletion()` functions
- Added global `--yes` flag to `cmd/root.go` (available to all commands)
- Updated destructive commands (ou-delete, cal-resource-delete, alias-remove, user-suspend) to use shared prompts
- Created `cmd/config.go` and `cmd/config-validate.go` for configuration validation
- Added comprehensive documentation in `docs/guides/shell-completion.md`
- Updated `README.md` and `docs/reference/commands.md` with new commands
- All features include unit tests with >80% coverage

</details>

<details>
<summary>📋 5. Performance</summary>

- [ ] Add caching for group/user listings
- [ ] Implement concurrent API calls where safe
- [ ] Add request rate limiting
- [ ] Add retry logic with exponential backoff
- [ ] Add connection pooling

**Rationale:** Improve performance and handle API quotas gracefully.

</details>

## Completed Items

For historical reference, the following items have been completed. See git history for implementation details:

1. ✅ Upgrade Go Version (Go 1.25)
2. ✅ Update Dependencies (google.golang.org/api v0.92.0 → v0.251.0)
3. ✅ Credential Management (file permissions, env vars, documentation)
4. ✅ Input Validation (email, phone, UUID, department, group names)
5. ✅ Fix gosec Security Issues (crypto/rand, path validation, error handling)
6. ✅ Remove Hardcoded Domain (configurable via flags/env/config)
7. ✅ Error Handling (zerolog structured logging)
8. ✅ CI/CD Modernization (GitHub Actions, GoReleaser, security scanning)
9. ✅ Documentation (README, CONTRIBUTING, guides, examples)
10. ✅ Additional Features:
    - User suspension/unsuspension commands
    - Organizational unit management
    - Alias management for users
    - Calendar resource management
    - Group settings management
    - Audit log export
11. ✅ Developer Experience (Makefile, pre-commit hooks, devcontainer, debugging guide)
