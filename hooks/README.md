# Git Hooks

This directory contains Git hooks to help maintain code quality.

## Pre-commit Hook

The pre-commit hook runs the following checks before each commit:

1. **go fmt** - Format all Go code
2. **go vet** - Run static analysis
3. **golangci-lint** - Run comprehensive linting
4. **go mod tidy** - Ensure module dependencies are clean
5. **go build** - Verify the code compiles

## Installation

To install the pre-commit hook, run:

```bash
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

Or use this one-liner:

```bash
ln -sf ../../hooks/pre-commit .git/hooks/pre-commit
```

## Requirements

- Go 1.25 or later
- golangci-lint installed (`brew install golangci-lint` or see https://golangci-lint.run/usage/install/)

## Bypassing the Hook

If you need to bypass the hook in an emergency (not recommended):

```bash
git commit --no-verify
```
