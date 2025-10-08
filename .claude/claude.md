# Google Admin Client (gac) - Claude Code Context

## Project Overview

**gac** is a command-line tool for managing Google Workspace administrative operations (users, groups, calendars, organizational units, etc.). Built with Go 1.25, it provides a simple interface for automating Google Workspace admin tasks through the Google Admin SDK APIs.

**Tech Stack:**
- Go 1.25 (recently upgraded)
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Zerolog](https://github.com/rs/zerolog) - Structured logging
- Google Admin SDK (Directory, Calendar, Groups Settings, Data Transfer APIs)
- OAuth2 for authentication

> **📘 Go-Specific Coding Standards**
> For comprehensive Go coding rules and best practices, see [claude-golang.md](./claude-golang.md).
> These rules cover error handling, concurrency, testing patterns, security, and more—and take precedence for all Go code in this project.

## Quick Architecture

```
main.go → cmd/root.go → cmd/{command}.go → Google Admin SDK APIs
                       ↓
                 cmd/client.go (OAuth2, API clients)
                 cmd/logger.go (zerolog)
                 cmd/helpers.go (utilities)
```

**Key Points:**
- **No `pkg/` directory** - All code lives in `cmd/`
- Each command is a separate file: `cmd/{resource}-{action}.go`
- Commands registered via Cobra in each file's `init()` function
- Centralized client initialization in `cmd/client.go`
- Structured logging via zerolog in `cmd/logger.go`

## Code Organization

### Directory Structure
```
cmd/
├── root.go                    # Root command, global flags, config loading
├── client.go                  # OAuth2 + Google API client initialization
├── logger.go                  # Zerolog logger setup and configuration
├── helpers.go                 # Shared utility functions
├── validation.go              # Input validation functions
├── {resource}.go              # Command group (e.g., user.go, group.go)
├── {resource}-{action}.go     # Specific commands (e.g., user-create.go)
└── {resource}_test.go         # Tests for resource commands
```

### Command Pattern

All commands follow this structure:

```go
package cmd

import (
    "github.com/spf13/cobra"
    admin "google.golang.org/api/admin/directory/v1"
)

// Flags/parameters
var (
    someFlag string
    anotherFlag []string
)

// Command definition
var someActionCmd = &cobra.Command{
    Use:   "action",
    Short: "Brief description",
    Long:  `Detailed description with examples`,
    Run:   someActionRunFunc,
}

// Registration in init()
func init() {
    parentCmd.AddCommand(someActionCmd)
    someActionCmd.Flags().StringVar(&someFlag, "flag-name", "", "description")
}

// Run function
func someActionRunFunc(cmd *cobra.Command, args []string) {
    // 1. Validate input
    // 2. Initialize API client
    // 3. Execute API operations
    // 4. Log and handle errors
    // 5. Output results
}
```

## Coding Patterns

### 1. Logging (Zerolog)

**Always use zerolog, never `fmt.Println()` or `log.Println()`:**

```go
// Import the global Logger from logger.go
Logger.Info().Str("user", email).Msg("Creating user")
Logger.Error().Err(err).Msg("Failed to create user")
Logger.Debug().Interface("response", resp).Msg("API response")
Logger.Fatal().Err(err).Msg("Critical error") // exits with status 1
```

**Log Levels:**
- `debug` - Detailed API calls, requests, responses (enabled with `-v` or `--log-level debug`)
- `info` - General operational messages (default)
- `warn` - Warnings (insecure permissions, deprecations)
- `error` - Errors
- `fatal` - Critical errors that require exit

**Structured Fields:**
- Use `.Str()`, `.Int()`, `.Bool()`, `.Err()`, etc. for typed fields
- Use `.Interface()` for complex objects
- Always end with `.Msg()` or `.Msgf()`

### 2. Error Handling

**Pattern:**
```go
if err != nil {
    Logger.Error().Err(err).Str("user", email).Msg("Failed to create user")
    os.Exit(1)
}
```

**For API errors with additional context:**
```go
if err != nil {
    Logger.Error().
        Err(err).
        Str("user", email).
        Str("group", group).
        Msg("Failed to add user to group")
    os.Exit(1)
}
```

### 3. API Client Initialization

**Always use factory functions from `cmd/client.go`:**

```go
// Directory API (users, groups, OUs)
client, err := newAdminClient()
if err != nil {
    Logger.Fatal().Err(err).Msg("Failed to initialize admin client")
}

// Calendar API
calClient, err := newCalendarClient()
if err != nil {
    Logger.Fatal().Err(err).Msg("Failed to initialize calendar client")
}

// Groups Settings API
groupsClient, err := newGroupsSettingsClient()
if err != nil {
    Logger.Fatal().Err(err).Msg("Failed to initialize groups settings client")
}
```

### 4. Configuration & Domain Management

**Domain configuration hierarchy** (highest to lowest priority):
1. CLI flag: `--domain example.com`
2. Environment variable: `GAC_DOMAIN` or `GOOGLE_ADMIN_DOMAIN`
3. Config file: `~/.google-admin.yaml` (`domain: example.com`)

**Get domain in commands:**
```go
domain := getDomain() // from root.go
if domain == "" {
    Logger.Error().Msg("Domain not configured")
    os.Exit(1)
}
```

### 5. Testing Conventions

**Current State:** 23.6% coverage, working toward >80%

**Test file naming:** `{resource}_test.go` (e.g., `user_test.go`, `group_test.go`)

**Test patterns:**
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "value",
            expected: "expected",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr = %v, got %v", tt.wantErr, err)
            }
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

**Testing priorities:**
- Helper functions (parsers, validators)
- Command flag parsing and registration
- Business logic functions
- Integration tests with mocked Google APIs (future)

## Development Workflow

### Make Targets

```bash
make build              # Build binary to build/gac
make test               # Run all tests
make test-coverage      # Generate coverage report (coverage.html)
make clean              # Remove build artifacts
make lint               # Run golangci-lint
make release            # Build for all platforms (goreleaser)
```

### Git Workflow

**IMPORTANT: The main branch is protected. All changes MUST be made via Pull Requests.**

**Required workflow for ALL changes:**

1. **Create feature branch FIRST (before any changes)**
   ```bash
   git checkout -b feat/feature-name     # For new features
   git checkout -b fix/bug-name          # For bug fixes
   git checkout -b chore/task-name       # For maintenance tasks
   ```

2. **Make your code changes**
   - Write/edit files
   - Add tests for new functionality

3. **Run quality checks BEFORE staging changes**
   ```bash
   make check    # MUST pass (includes: fmt, vet, lint, security scan, and tests)
   ```

4. **Stage and commit changes**
   ```bash
   git add <files>
   git commit -m "type: descriptive message"
   ```

5. **Push branch to remote**
   ```bash
   git push -u origin <branch-name>
   ```

6. **Create Pull Request**
   ```bash
   gh pr create --title "..." --body "..."
   ```

**Key Rules:**
- ❌ NEVER commit directly to main
- ❌ NEVER push to main
- ✅ ALWAYS create a feature branch first
- ✅ ALWAYS run `make check` before `git add` (includes all tests)
- ✅ All changes merge to main via approved PRs only

**Branch Naming Conventions:**
- `feat/` - New features or enhancements
- `fix/` - Bug fixes
- `chore/` - Maintenance, documentation, dependencies
- `test/` - Test additions or improvements

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestFunctionName ./cmd
```

### Pre-commit Hooks

The project uses pre-commit hooks (`.pre-commit-config.yaml`):
- Go formatting
- Go linting
- Trailing whitespace removal
- YAML/JSON validation

### Debugging

```bash
# Enable verbose logging
gac --verbose user list
gac -v user list

# Set specific log level
gac --log-level debug user list

# JSON log output (for parsing)
gac --json-log user list
```

## Google API Integration

### OAuth2 Flow

1. **First run:** User authorizes via browser, code is exchanged for token
2. **Subsequent runs:** Token loaded from cache (`~/.credentials/gac.json`)
3. **Token refresh:** Automatic refresh when expired

**Credential files:**
- `~/.credentials/client_secret.json` - OAuth2 client credentials from Google Cloud Console
- `~/.credentials/gac.json` - Cached access/refresh tokens

### Required OAuth2 Scopes

```go
admin.AdminDirectoryUserScope                    // Create/update users
admin.AdminDirectoryUserReadonlyScope            // Read users
admin.AdminDirectoryGroupScope                   // Manage groups
admin.AdminDirectoryGroupReadonlyScope           // Read groups
admin.AdminDirectoryGroupMemberScope             // Manage group members
admin.AdminDirectoryOrgunitScope                 // Manage organizational units
groupssettings.AppsGroupsSettingsScope           // Manage group settings
calendar.CalendarScope                           // Full calendar access
calendar.CalendarEventsScope                     // Manage calendar events
datatransfer.AdminDatatransferScope              // Data transfer operations
```

### Google Admin SDK APIs Used

- **Directory API** (`admin/directory/v1`): Users, Groups, OUs, Aliases
- **Groups Settings API** (`groupssettings/v1`): Group permissions and settings
- **Calendar API** (`calendar/v3`): Events and calendar resources
- **Data Transfer API** (`datatransfer/v1`): Data ownership transfers

## Current Development Priorities

### From TODO.md (as of Oct 2024)

**Completed:**
✅ Upgraded to Go 1.25
✅ Updated all dependencies (google.golang.org/api v0.92.0 → v0.251.0)
✅ Implemented structured logging with zerolog
✅ Added error handling improvements

**In Progress:**
🚧 **Testing** - Current: 23.6% coverage, Goal: >80%
  - Helper functions tested ✅
  - Command flag parsing tested ✅
  - Need: Integration tests with mocked APIs
  - Need: Command runner function tests
  - Need: Client initialization tests

**Upcoming:**
- Rate limiting for API calls
- Retry logic with exponential backoff
- Better error messages with recovery suggestions
- System keychain integration for credentials

## Key Conventions

### Flag Naming
- Use kebab-case: `--first-name`, `--group-email`
- Short flags for common options: `-v` (verbose), `-g` (groups), `-e` (email)
- Boolean flags: No value needed (e.g., `--verbose`)

### Variable Naming
- Go standard: camelCase
- Exported: PascalCase
- Package-level: camelCase with descriptive names

### Command Naming
- Pattern: `{resource} {action}` (e.g., `user create`, `group-settings update`)
- Resources: `user`, `group`, `group-settings`, `calendar`, `cal-resource`, `ou`, `alias`, `transfer`
- Actions: `create`, `list`, `update`, `delete`, `suspend`, `unsuspend`, `add`, `remove`

### File Naming
- Command files: `{resource}-{action}.go` (e.g., `user-create.go`, `group-list.go`)
- Test files: `{resource}_test.go` (e.g., `user_test.go`)
- Core files: `{purpose}.go` (e.g., `client.go`, `logger.go`, `helpers.go`)

## Common Tasks

### Adding a New Command

1. **Create command file:** `cmd/{resource}-{action}.go`
2. **Define command struct:**
   ```go
   var newCmd = &cobra.Command{
       Use:   "action",
       Short: "Brief description",
       Long:  `Detailed description`,
       Run:   actionRunFunc,
   }
   ```
3. **Register in `init()`:**
   ```go
   func init() {
       resourceCmd.AddCommand(newCmd)
       newCmd.Flags().StringVar(&flag, "flag-name", "", "description")
   }
   ```
4. **Implement run function:**
   ```go
   func actionRunFunc(cmd *cobra.Command, args []string) {
       // Implementation
   }
   ```
5. **Add tests:** Create/update `cmd/{resource}_test.go`
6. **Update docs:** README.md, relevant docs/ files

### Adding a New Google API Integration

1. **Add import:** Import the API package (e.g., `"google.golang.org/api/newapi/v1"`)
2. **Add factory function in `cmd/client.go`:**
   ```go
   func newNewAPIClient() (*newapi.Service, error) {
       client, err := newHTTPClient()
       if err != nil {
           return nil, err
       }
       return newapi.NewService(context.Background(), option.WithHTTPClient(client))
   }
   ```
3. **Add required scopes** to `scopes` variable in `client.go`
4. **Update README** with new features

## Documentation References

- **README.md** - User-facing documentation, installation, usage examples
- **ARCHITECTURE.md** - Detailed technical architecture, design decisions
- **CONTRIBUTING.md** - Development setup, coding guidelines, PR process
- **DEBUGGING.md** - Debugging techniques, profiling
- **TODO.md** - Current priorities, roadmap, completed items
- **docs/** - Comprehensive guides for each feature
- **examples/** - Runnable example scripts

## Tips for Working with This Project

1. **Always use zerolog for logging** - No fmt.Println() in production code
2. **Test your changes** - Add tests for new functions (working toward >80% coverage)
3. **Follow the command pattern** - Look at existing commands (e.g., `user-create.go`) as templates
4. **Check TODO.md** - See current priorities and avoid duplicate work
5. **Use make targets** - `make test`, `make build`, `make test-coverage`
6. **Run pre-commit hooks** - Ensure code formatting and linting
7. **Update docs** - Keep README and relevant docs/ files in sync with code changes
8. **Be mindful of API rate limits** - Google Admin APIs have quotas
9. **Test with a test domain** - Don't test against production Google Workspace

## Common Pitfalls to Avoid

- ❌ Don't create a `pkg/` directory (all code in `cmd/`)
- ❌ Don't use `fmt.Println()` or `log.Println()` (use zerolog)
- ❌ Don't hardcode domains (use `getDomain()` function)
- ❌ Don't skip tests (coverage goal is >80%)
- ❌ Don't handle errors without logging context
- ❌ Don't create API clients directly (use factory functions in `client.go`)
- ❌ Don't commit credentials or tokens (`.gitignore` protects you, but be careful)
