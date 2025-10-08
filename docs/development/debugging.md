# Debugging Guide

This guide provides information on how to debug the `gac` (Google Admin Client) application.

## Prerequisites

- Go 1.25+ installed
- Valid Google Workspace OAuth2 credentials
- Delve debugger installed (optional): `go install github.com/go-delve/delve/cmd/dlv@latest`

## Building for Debugging

Build the binary with debug symbols:

```bash
go build -gcflags="all=-N -l" -o gac
```

Or use the Makefile:

```bash
make build
```

## Debugging Methods

### 1. Using Print Statements

Add debug print statements to track execution:

```go
import "fmt"

func someFunction() {
    fmt.Printf("DEBUG: variable value = %+v\n", variable)
}
```

### 2. Using Delve (Command Line)

Start debugging with Delve:

```bash
# Debug the main package
dlv debug

# Debug with arguments
dlv debug -- user list

# Debug a specific test
dlv test ./cmd -- -test.run TestUserCreate
```

Common Delve commands:
- `break (b)` - Set breakpoint
- `continue (c)` - Continue execution
- `next (n)` - Step over
- `step (s)` - Step into
- `print (p)` - Print variable
- `list (l)` - Show source code
- `quit (q)` - Exit debugger

### 3. Using VS Code

Add this to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug gac",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": ["--help"]
    },
    {
      "name": "Debug user list",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": ["user", "list"]
    },
    {
      "name": "Attach to Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickProcess}"
    }
  ]
}
```

### 4. Using GoLand / IntelliJ IDEA

1. Set breakpoints by clicking in the gutter
2. Right-click on `main.go` → Debug 'go build'
3. Configure run configurations with arguments

## Common Debugging Scenarios

### OAuth2 Authentication Issues

Enable verbose OAuth2 debugging:

```go
// In cmd/client.go, add before creating HTTP client:
import "log"
import "net/http/httputil"

// Add request/response logging
client.Transport = &loggingTransport{http.DefaultTransport}

type loggingTransport struct {
    rt http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    dump, _ := httputil.DumpRequest(req, true)
    log.Printf("REQUEST:\n%s", dump)
    resp, err := t.rt.RoundTrip(req)
    if resp != nil {
        dump, _ := httputil.DumpResponse(resp, true)
        log.Printf("RESPONSE:\n%s", dump)
    }
    return resp, err
}
```

### API Call Issues

Check the API response:

```bash
# Set verbose mode (if implemented)
./gac --verbose user list

# Check credential cache
cat ~/.credentials/gac.json | jq .

# Verify client secret
cat ~/.credentials/client_secret.json | jq .
```

### Configuration Issues

Debug configuration loading:

```go
// In cmd/root.go, add to initConfig():
import "github.com/spf13/viper"

fmt.Printf("Config file used: %s\n", viper.ConfigFileUsed())
fmt.Printf("All settings: %+v\n", viper.AllSettings())
```

### Google API Rate Limiting

Add retry logic with exponential backoff:

```go
import "time"
import "google.golang.org/api/googleapi"

func retryableCall(fn func() error) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }

        if apiErr, ok := err.(*googleapi.Error); ok {
            if apiErr.Code == 429 || apiErr.Code >= 500 {
                backoff := time.Duration(1<<uint(i)) * time.Second
                fmt.Printf("Retrying after %v (attempt %d/%d)\n", backoff, i+1, maxRetries)
                time.Sleep(backoff)
                continue
            }
        }
        return err
    }
    return fmt.Errorf("max retries exceeded")
}
```

## Debugging Tests

Run tests with verbose output:

```bash
# Run all tests verbosely
go test -v ./...

# Run specific test
go test -v ./cmd -run TestUserCreate

# Run with race detector
go test -race ./...

# Generate coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Environment Variables

Set these for debugging:

```bash
# Enable Go debugging
export GODEBUG=gctrace=1

# HTTP debugging
export GODEBUG=http2debug=1

# Disable HTTP/2
export GODEBUG=http2client=0
```

## Troubleshooting

### "Unable to create client" Error

1. Check if `client_secret.json` exists in `~/.credentials/`
2. Verify the JSON format is valid
3. Ensure OAuth2 scopes are correct

### "Token expired" Error

1. Delete cached token: `rm ~/.credentials/gac.json`
2. Re-authenticate: `./gac init -f`

### "Permission denied" Errors

1. Verify the Google Workspace admin account has necessary permissions
2. Check OAuth2 scopes include required permissions
3. Review Google Admin Console → Security → API Controls

### Build Errors

```bash
# Clean and rebuild
make clean
make build

# Update dependencies
go mod tidy
go mod download
```

## Logging

Add structured logging (future enhancement):

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("User operation", "email", email, "action", "create")
```

## Profiling

Profile CPU usage:

```bash
go build -o gac
./gac --cpuprofile=cpu.prof user list
go tool pprof cpu.prof
```

Profile memory usage:

```bash
go build -o gac
./gac --memprofile=mem.prof user list
go tool pprof mem.prof
```

## Resources

- [Delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation)
- [VS Code Go Debugging](https://github.com/golang/vscode-go/wiki/debugging)
- [Google APIs Go Client](https://pkg.go.dev/google.golang.org/api)
- [Cobra Documentation](https://github.com/spf13/cobra)
