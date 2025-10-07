# Contributing to Google Admin Client (gac)

Thank you for your interest in contributing to the Google Admin Client! We welcome contributions of all kinds: bug reports, feature requests, documentation improvements, and code contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Features](#suggesting-features)
  - [Submitting Pull Requests](#submitting-pull-requests)
- [Development Guidelines](#development-guidelines)
  - [Code Style](#code-style)
  - [Testing](#testing)
  - [Documentation](#documentation)
  - [Commit Messages](#commit-messages)
- [Project Structure](#project-structure)
- [Testing](#testing-1)
- [Release Process](#release-process)

## Code of Conduct

This project follows a Code of Conduct to ensure a welcoming environment for all contributors. By participating, you are expected to:

- Be respectful and inclusive
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/google-admin-client.git
   cd google-admin-client
   ```
3. **Add the upstream repository** as a remote:
   ```bash
   git remote add upstream https://github.com/acockrell/google-admin-client.git
   ```
4. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- **Go 1.25 or later** - [Install Go](https://golang.org/doc/install)
- **Make** - For build automation
- **Git** - For version control
- **Google Cloud account** - For testing with Google Workspace APIs

### Local Development Environment

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Build the project**:
   ```bash
   make build
   ```

3. **Run tests**:
   ```bash
   make test
   ```

4. **Run linters**:
   ```bash
   make lint
   ```

### Development Container (Optional)

We provide a VS Code devcontainer for consistent development environments:

1. Install [Docker](https://www.docker.com/get-started) and [VS Code](https://code.visualstudio.com/)
2. Install the [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension
3. Open the project in VS Code
4. Click "Reopen in Container" when prompted

The devcontainer includes:
- Go 1.25
- All development tools (golangci-lint, delve, etc.)
- Recommended VS Code extensions
- Credential file mounting support

### Available Make Targets

```bash
make help          # Show all available targets
make build         # Build the binary
make test          # Run all tests
make test-coverage # Run tests with coverage report
make lint          # Run golangci-lint
make fmt           # Format code with go fmt
make vet           # Run go vet
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make release-test  # Test release configuration
```

## How to Contribute

### Reporting Bugs

Before creating a bug report:
1. **Search existing issues** to avoid duplicates
2. **Update to the latest version** to see if the issue persists
3. **Gather information** about your environment

When creating a bug report, include:
- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs. actual behavior
- **Environment details**:
  - OS and version
  - Go version
  - `gac` version (`gac version`)
  - Relevant configuration (redact sensitive info)
- **Error messages and logs**
- **Screenshots** if applicable

**Example bug report**:
```markdown
## Bug Description
User creation fails when adding user to multiple groups

## Steps to Reproduce
1. Run: `gac user create -g group1 -g group2 newuser@example.com`
2. Observe error: "failed to add user to group2"

## Expected Behavior
User should be created and added to both groups

## Actual Behavior
User is created but only added to group1

## Environment
- OS: macOS 14.5
- Go: 1.25.1
- gac: v0.2.0
```

### Suggesting Features

Feature requests are welcome! Before submitting:
1. **Check existing feature requests** in GitHub Issues
2. **Consider if it fits** the project scope
3. **Think about backward compatibility**

When suggesting a feature:
- **Describe the problem** you're trying to solve
- **Propose a solution** with examples
- **Explain why it's useful** to the broader community
- **Consider implementation details** if possible

**Example feature request**:
```markdown
## Feature Request: Bulk User Import from CSV

### Problem
Creating multiple users one-by-one is time-consuming for onboarding

### Proposed Solution
Add a `gac user import` command that accepts a CSV file:
```csv
email,firstName,lastName,department,groups
jdoe@example.com,John,Doe,Engineering,"dev,all-staff"
```

### Benefits
- Faster bulk user creation
- Easier migration from other systems
- Reduced human error

### Implementation Considerations
- Support dry-run mode
- Validate all entries before creating
- Provide detailed error reporting
```

### Submitting Pull Requests

1. **Create an issue first** for significant changes to discuss the approach
2. **Keep PRs focused** - one feature or bug fix per PR
3. **Write tests** for new functionality
4. **Update documentation** as needed
5. **Follow code style** guidelines
6. **Ensure CI passes** before requesting review

**PR Checklist**:
- [ ] Tests pass locally (`make test`)
- [ ] Linters pass (`make lint`)
- [ ] Code is formatted (`make fmt`)
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages follow conventions
- [ ] PR description explains the change

**Example PR description**:
```markdown
## Summary
Adds input validation for phone numbers to prevent invalid formats

## Changes
- Created `ValidatePhoneNumber()` function in `cmd/validation.go`
- Added comprehensive test cases in `cmd/validation_test.go`
- Integrated validation into `user create` and `user update` commands
- Updated README with phone number format requirements

## Testing
- Added 15 test cases covering valid and invalid formats
- Tested manually with various phone number inputs
- All existing tests still pass

## Related Issues
Fixes #42
```

## Development Guidelines

### Code Style

We follow standard Go conventions:

1. **Use `gofmt`** for formatting:
   ```bash
   make fmt
   ```

2. **Follow Go Code Review Comments**: https://go.dev/wiki/CodeReviewComments

3. **Run linters**:
   ```bash
   make lint
   ```
   We use `golangci-lint` with the configuration in `.golangci.yml`

4. **Keep functions focused** - functions should do one thing well

5. **Use meaningful names** - prioritize clarity over brevity

6. **Add comments** for:
   - Exported functions and types (required)
   - Complex logic
   - Non-obvious decisions

7. **Error handling**:
   - Always check errors
   - Provide context with error wrapping: `fmt.Errorf("operation failed: %w", err)`
   - Use meaningful error messages

### Testing

We aim for high test coverage (>80%) to ensure reliability:

1. **Write tests for new code**:
   ```go
   func TestValidateEmail(t *testing.T) {
       tests := []struct {
           name    string
           email   string
           wantErr bool
       }{
           {"valid email", "user@example.com", false},
           {"invalid format", "not-an-email", true},
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               err := ValidateEmail(tt.email)
               if (err != nil) != tt.wantErr {
                   t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
               }
           })
       }
   }
   ```

2. **Run tests before submitting**:
   ```bash
   make test
   ```

3. **Check coverage**:
   ```bash
   make test-coverage
   ```

4. **Test types**:
   - **Unit tests**: Test individual functions in isolation
   - **Integration tests**: Test components working together (use mocks for external APIs)
   - **Edge cases**: Test boundary conditions and error paths

5. **Use table-driven tests** for testing multiple inputs

6. **Mock external dependencies** (Google APIs) for reliable tests

### Documentation

Good documentation is crucial:

1. **Code comments**:
   - All exported functions, types, and constants must have doc comments
   - Start with the name of the thing being documented
   - Use complete sentences

   ```go
   // ValidateEmail checks if the provided email address is valid according to RFC 5322.
   // It returns an error if the email is invalid or exceeds length limits.
   func ValidateEmail(email string) error {
       // ...
   }
   ```

2. **README updates**:
   - Add new features to the usage examples
   - Update command reference for new commands/flags
   - Add troubleshooting entries for common issues

3. **Inline comments**:
   - Explain "why" not "what" (the code shows what)
   - Document non-obvious decisions
   - Mark TODOs clearly: `// TODO: Add retry logic`

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

Format: `<type>(<scope>): <description>`

**Types**:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, missing semi-colons, etc.)
- `refactor:` - Code refactoring without changing functionality
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks (dependencies, build config, etc.)
- `perf:` - Performance improvements
- `ci:` - CI/CD changes

**Examples**:
```bash
feat(user): add bulk user import from CSV
fix(auth): resolve token refresh race condition
docs(readme): add installation instructions for Windows
test(validation): add test cases for phone number validation
refactor(client): simplify OAuth2 client initialization
chore(deps): update google.golang.org/api to v0.251.0
```

**Commit message guidelines**:
- Use imperative mood: "add feature" not "added feature"
- Keep subject line under 72 characters
- Add body for complex changes explaining what and why
- Reference issues: `Fixes #123` or `Closes #456`

## Project Structure

```
google-admin-client/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and configuration
│   ├── user.go            # User command group
│   ├── user-create.go     # User creation command
│   ├── user-list.go       # User listing command
│   ├── user-update.go     # User update command
│   ├── group.go           # Group command group
│   ├── group-list.go      # Group listing command
│   ├── calendar.go        # Calendar command group
│   ├── transfer.go        # Transfer command
│   ├── client.go          # Google API client setup
│   ├── validation.go      # Input validation functions
│   └── *_test.go          # Test files
├── .github/
│   └── workflows/         # GitHub Actions workflows
├── examples/              # Example configurations
├── main.go               # Application entry point
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── Makefile              # Build automation
├── Dockerfile            # Docker image definition
├── .goreleaser.yml       # Release configuration
├── README.md             # User documentation
├── CONTRIBUTING.md       # This file
├── ARCHITECTURE.md       # Technical architecture docs
├── DEBUGGING.md          # Debugging guide
├── CREDENTIALS.md        # OAuth2 setup guide
└── TODO.md               # Project roadmap
```

### Adding a New Command

1. Create the command file in `cmd/`:
   ```go
   // cmd/myfeature-action.go
   package cmd

   import (
       "github.com/spf13/cobra"
   )

   var myfeatureActionCmd = &cobra.Command{
       Use:   "action [args]",
       Short: "Brief description",
       Long:  `Detailed description with examples`,
       RunE:  myfeatureActionRunFunc,
   }

   func init() {
       myfeatureCmd.AddCommand(myfeatureActionCmd)
       // Add flags here
   }

   func myfeatureActionRunFunc(cmd *cobra.Command, args []string) error {
       // Implementation
       return nil
   }
   ```

2. Write tests in `cmd/myfeature-action_test.go`
3. Update README with usage examples
4. Add integration with existing client setup

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
go test -race ./...

# Run specific test
go test -v -run TestValidateEmail ./cmd

# Run tests in a specific package
go test ./cmd
```

### Test Organization

- Place test files alongside the code they test
- Name test files with `_test.go` suffix
- Use table-driven tests for multiple test cases
- Group related tests with subtests

### Mocking External Dependencies

For Google API clients, use interfaces and provide mock implementations:

```go
type UserService interface {
    Insert(*admin.User) (*admin.User, error)
    Get(userKey string) (*admin.User, error)
}

type mockUserService struct {
    insertFunc func(*admin.User) (*admin.User, error)
    getFunc    func(string) (*admin.User, error)
}

func (m *mockUserService) Insert(u *admin.User) (*admin.User, error) {
    return m.insertFunc(u)
}
```

## Release Process

Releases are automated via GitHub Actions and GoReleaser. See [RELEASE.md](RELEASE.md) for details.

### Version Numbering

We use [Semantic Versioning](https://semver.org/):
- **MAJOR** version for incompatible API changes
- **MINOR** version for backward-compatible functionality
- **PATCH** version for backward-compatible bug fixes

### Creating a Release

1. Update version in relevant files
2. Update CHANGELOG (if maintained)
3. Create and push a version tag:
   ```bash
   git tag -a v0.3.0 -m "Release v0.3.0"
   git push origin v0.3.0
   ```
4. GitHub Actions will automatically build and publish the release

## Getting Help

- **Documentation**: Check [README.md](README.md), [ARCHITECTURE.md](ARCHITECTURE.md), and [DEBUGGING.md](DEBUGGING.md)
- **Issues**: Search existing [GitHub Issues](https://github.com/acockrell/google-admin-client/issues)
- **Discussions**: Start a [GitHub Discussion](https://github.com/acockrell/google-admin-client/discussions) for questions

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License).

## Acknowledgments

Thank you for contributing to making Google Admin Client better for everyone!
