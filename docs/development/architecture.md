# Architecture Documentation

## Overview

`gac` (Google Admin Client) is a command-line interface tool for managing Google Workspace administrative operations. It follows a modular architecture using the Cobra CLI framework and integrates with Google Admin SDK APIs.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        gac CLI                              │
├─────────────────────────────────────────────────────────────┤
│  main.go  →  cmd.Execute()                                  │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   Command Layer (cmd/)                       │
├──────────────┬──────────────┬──────────────┬────────────────┤
│   User Ops   │  Group Ops   │ Calendar Ops │  Transfer Ops  │
│              │              │              │                │
│ - create     │ - list       │ - create     │ - transfer     │
│ - list       │              │ - list       │                │
│ - update     │              │ - update     │                │
└──────┬───────┴──────┬───────┴──────┬───────┴────────┬───────┘
       │              │              │                │
       └──────────────┴──────────────┴────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Client Layer (cmd/client.go)                    │
├─────────────────────────────────────────────────────────────┤
│  - newAdminClient()        - newHTTPClient()                │
│  - newCalendarClient()     - OAuth2 Token Management        │
│  - newDataTransferClient() - Credential Caching             │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│              Google Admin SDK APIs                          │
├──────────────┬──────────────┬──────────────┬────────────────┤
│ Directory API│ Calendar API │  Groups API  │ DataTransfer   │
└──────────────┴──────────────┴──────────────┴────────────────┘
```

## Project Structure

```
gac/
├── main.go                     # Entry point
├── cmd/                        # Command implementations
│   ├── root.go                 # Root command and CLI setup
│   ├── client.go               # Google API client initialization
│   ├── helpers.go              # Shared utility functions
│   ├── init.go                 # OAuth2 initialization command
│   ├── user.go                 # User command group
│   ├── user-create.go          # User creation logic
│   ├── user-list.go            # User listing logic
│   ├── user-update.go          # User update logic
│   ├── group.go                # Group command group
│   ├── group-list.go           # Group listing logic
│   ├── calendar.go             # Calendar command group
│   ├── calendar-create.go      # Calendar event creation
│   ├── calendar-list.go        # Calendar event listing
│   ├── calendar-update.go      # Calendar event updates
│   └── transfer.go             # Data transfer operations
├── build/                      # Build artifacts (generated)
├── .devcontainer/              # VS Code dev container config
├── .github/workflows/          # CI/CD workflows
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── Makefile                    # Build automation
├── Dockerfile                  # Container image definition
├── README.md                   # User documentation
├── ARCHITECTURE.md             # This file
└── DEBUGGING.md                # Debugging guide
```

## Core Components

### 1. Main Entry Point (`main.go`)

- Minimal entry point that delegates to `cmd.Execute()`
- No business logic, just bootstraps the CLI

### 2. Command Layer (`cmd/`)

#### Root Command (`cmd/root.go`)

- Defines the root `gac` command
- Sets up global flags: `--config`, `--secret`, `--cache`
- Initializes Viper configuration management
- Handles configuration file loading from `$HOME/.google-admin.yaml`

#### Client Management (`cmd/client.go`)

**Responsibilities:**
- Creating authenticated HTTP clients
- Managing OAuth2 token lifecycle
- Caching credentials to `~/.credentials/gac.json`
- Providing factory methods for Google API service clients

**Key Functions:**
- `newHTTPClient()` - Creates OAuth2 authenticated HTTP client
- `newAdminClient()` - Returns Directory API service
- `newCalendarClient()` - Returns Calendar API service
- `newDataTransferClient()` - Returns Data Transfer API service
- `getTokenFromWeb()` - Handles OAuth2 authorization flow
- `tokenFromFile()` / `saveToken()` - Token persistence

**OAuth2 Scopes:**
```go
scopes = []string{
    admin.AdminDirectoryUserReadonlyScope,
    admin.AdminDirectoryUserScope,
    admin.AdminDirectoryGroupReadonlyScope,
    admin.AdminDirectoryGroupMemberReadonlyScope,
    admin.AdminDirectoryGroupMemberScope,
    calendar.CalendarScope,
    calendar.CalendarReadonlyScope,
    calendar.CalendarEventsScope,
    calendar.CalendarEventsReadonlyScope,
    datatransfer.AdminDatatransferScope,
}
```

#### Command Implementations

Each command follows a consistent pattern:

```go
var someCmd = &cobra.Command{
    Use:   "command-name",
    Short: "Brief description",
    Long:  "Detailed description with examples",
    Run:   commandRunFunc,
}

func init() {
    parentCmd.AddCommand(someCmd)
    someCmd.Flags().StringVar(&variable, "flag", "default", "description")
}

func commandRunFunc(cmd *cobra.Command, args []string) {
    // 1. Validate input
    // 2. Create API client
    // 3. Execute API calls
    // 4. Handle errors
    // 5. Output results
}
```

### 3. Helper Functions (`cmd/helpers.go`)

Provides shared utility functions:
- Error handling and exit functions
- Common parsing logic
- Validation functions

## Data Flow

### Example: Creating a User

```
1. User Input
   $ gac user create -g engineering newuser@example.com

2. CLI Parsing (Cobra)
   ├─ Parse command: "user create"
   ├─ Parse flags: -g engineering
   └─ Parse args: newuser@example.com

3. Command Execution (cmd/user-create.go)
   ├─ Validate email address
   ├─ Prompt for first/last name, personal email
   └─ Generate random password

4. Client Creation (cmd/client.go)
   ├─ Load client_secret.json
   ├─ Check for cached token
   ├─ Create OAuth2 HTTP client
   └─ Initialize Admin Directory API client

5. API Call (Google Admin SDK)
   ├─ Create User object
   ├─ Call client.Users.Insert()
   └─ Handle response/errors

6. Group Assignment (if -g flag provided)
   ├─ For each group:
   ├─ Append @domain.com (currently hardcoded)
   ├─ Call client.Members.Insert()
   └─ Handle errors

7. Output Results
   └─ Print user details and credentials
```

## Configuration Management

### Configuration File

Location: `$HOME/.google-admin.yaml`

Managed by Viper, supports:
- YAML configuration
- Environment variables (prefixed with `GOOGLE_ADMIN_`)
- Command-line flags (highest priority)

### Credentials

**Client Secret:**
- Location: `$HOME/.credentials/client_secret.json`
- Override: `--secret` flag
- Contains OAuth2 client credentials from Google Cloud Console

**Cached Token:**
- Location: `$HOME/.credentials/gac.json`
- Override: `--cache` flag
- Contains OAuth2 access/refresh tokens
- Automatically refreshed when expired

## Authentication Flow

```
1. First Run / Token Missing
   ├─ Load client_secret.json
   ├─ Generate OAuth2 authorization URL
   ├─ User opens URL in browser
   ├─ User grants permissions
   ├─ User copies authorization code
   ├─ Exchange code for access token
   └─ Save token to cache file

2. Subsequent Runs
   ├─ Load cached token
   ├─ Check if expired
   ├─ If expired: refresh using refresh token
   └─ Use access token for API calls
```

## Error Handling

**Current Approach:**
- Most errors result in `exitWithError()` which prints to stderr and exits
- API errors are caught and reported with context
- No retry logic (see TODO for improvements)

**Best Practices:**
- Always validate input before API calls
- Provide meaningful error messages
- Include affected resource (email, group, etc.) in errors

## API Integration

### Google Admin Directory API

**Users:**
- `client.Users.Insert()` - Create user
- `client.Users.Get()` - Get user details
- `client.Users.List()` - List all users
- `client.Users.Update()` - Update user attributes

**Groups:**
- `client.Groups.Get()` - Get group details
- `client.Groups.List()` - List all groups
- `client.Members.List()` - List group members
- `client.Members.Insert()` - Add user to group

### Google Calendar API

**Events:**
- `client.Events.Insert()` - Create event
- `client.Events.List()` - List events
- `client.Events.Update()` - Update event

### Data Transfer API

**Transfers:**
- `client.Transfers.Insert()` - Initiate data transfer

## Extension Points

### Adding a New Command

1. Create new file in `cmd/` (e.g., `cmd/newfeature.go`)
2. Define command structure:
   ```go
   var newFeatureCmd = &cobra.Command{
       Use:   "new-feature",
       Short: "Description",
       Run:   newFeatureRunFunc,
   }
   ```
3. Register in `init()`:
   ```go
   func init() {
       rootCmd.AddCommand(newFeatureCmd)
   }
   ```
4. Implement `newFeatureRunFunc()`
5. Update README with examples

### Adding a New API Client

1. Add to `cmd/client.go`:
   ```go
   func newSomeAPIClient() (*someapi.Service, error) {
       client, err := newHTTPClient()
       if err != nil {
           return nil, err
       }
       srv, err := someapi.NewService(context.Background(),
           option.WithHTTPClient(client))
       return srv, err
   }
   ```
2. Add required scopes to `scopes` variable
3. Import the API package in `go.mod`

## Design Decisions

### Why Cobra?

- Industry-standard CLI framework
- Automatic help generation
- Nested command support
- Flag parsing and validation

### Why Viper?

- Configuration file support
- Environment variable integration
- Works seamlessly with Cobra
- Multiple format support (YAML, JSON, etc.)

### Why Single Binary?

- Easy distribution
- No runtime dependencies
- Simple deployment

## Security Considerations

1. **Credential Storage:**
   - Tokens stored in `~/.credentials/` with 0700 permissions
   - Never commit credentials to version control
   - `.gitignore` excludes credential files

2. **OAuth2 Best Practices:**
   - Uses authorization code flow
   - Stores refresh tokens for long-term access
   - Scopes limited to required permissions

3. **API Communication:**
   - All API calls over HTTPS
   - Google handles TLS/SSL

## Performance Considerations

- **No Connection Pooling:** Each command creates new HTTP client
- **No Caching:** API responses not cached
- **Sequential Operations:** Bulk operations execute sequentially
- **Future Improvements:** See TODO.md for caching and concurrency plans

## Testing Strategy

Currently: No automated tests exist (see TODO.md)

Recommended approach:
- Unit tests for helper functions
- Integration tests with mocked Google APIs
- End-to-end tests against test Google Workspace domain

## Deployment

### Binary Distribution

Build for multiple platforms:
```bash
GOOS=linux GOARCH=amd64 go build -o gac-linux
GOOS=darwin GOARCH=amd64 go build -o gac-darwin
GOOS=windows GOARCH=amd64 go build -o gac-windows.exe
```

### Docker

```bash
docker build -t gac:latest .
docker run -v ~/.credentials:/root/.credentials gac user list
```

### CI/CD

GitHub Actions workflow (`.github/workflows/gosec.yml`):
- Runs on push and PRs
- Builds with Go 1.25
- Runs security scanning (gosec)
- Verifies dependencies

## Future Architecture Improvements

See TODO.md for planned enhancements:
- Add structured logging
- Implement retry logic with exponential backoff
- Add request rate limiting
- Support for multiple domains/organizations
- Configuration-based domain management
- Comprehensive test coverage
- Performance optimizations (caching, concurrency)
