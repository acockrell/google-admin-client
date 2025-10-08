package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected zerolog.Level
	}{
		{"debug lowercase", "debug", zerolog.DebugLevel},
		{"debug uppercase", "DEBUG", zerolog.DebugLevel},
		{"info lowercase", "info", zerolog.InfoLevel},
		{"info uppercase", "INFO", zerolog.InfoLevel},
		{"warn lowercase", "warn", zerolog.WarnLevel},
		{"warn uppercase", "WARN", zerolog.WarnLevel},
		{"warning", "warning", zerolog.WarnLevel},
		{"error lowercase", "error", zerolog.ErrorLevel},
		{"error uppercase", "ERROR", zerolog.ErrorLevel},
		{"fatal", "fatal", zerolog.FatalLevel},
		{"panic", "panic", zerolog.PanicLevel},
		{"disabled", "disabled", zerolog.Disabled},
		{"invalid defaults to info", "invalid", zerolog.InfoLevel},
		{"empty defaults to info", "", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		logLevel   string
		jsonFormat bool
		wantLevel  zerolog.Level
	}{
		{
			name:       "verbose enables debug",
			verbose:    true,
			logLevel:   "info",
			jsonFormat: false,
			wantLevel:  zerolog.DebugLevel,
		},
		{
			name:       "info level without verbose",
			verbose:    false,
			logLevel:   "info",
			jsonFormat: false,
			wantLevel:  zerolog.InfoLevel,
		},
		{
			name:       "error level",
			verbose:    false,
			logLevel:   "error",
			jsonFormat: false,
			wantLevel:  zerolog.ErrorLevel,
		},
		{
			name:       "warn level",
			verbose:    false,
			logLevel:   "warn",
			jsonFormat: false,
			wantLevel:  zerolog.WarnLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogger(tt.verbose, tt.logLevel, tt.jsonFormat)
			currentLevel := zerolog.GlobalLevel()
			if currentLevel != tt.wantLevel {
				t.Errorf("InitLogger() set level to %v, want %v", currentLevel, tt.wantLevel)
			}
		})
	}
}

func TestLogAPICall(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to debug level to capture debug logs
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	LogAPICall("test-service", "test-method", map[string]interface{}{
		"key": "value",
	})

	output := buf.String()
	if output == "" {
		t.Error("LogAPICall() produced no output")
	}
	// Basic check that the log contains expected fields
	if !bytes.Contains(buf.Bytes(), []byte("test-service")) {
		t.Error("Log output should contain service name")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test-method")) {
		t.Error("Log output should contain method name")
	}
}

func TestLogAPIResponse(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to debug level to capture debug logs
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	LogAPIResponse("test-service", "test-method", 200, 100*time.Millisecond)

	output := buf.String()
	if output == "" {
		t.Error("LogAPIResponse() produced no output")
	}
	// Basic check that the log contains expected fields
	if !bytes.Contains(buf.Bytes(), []byte("test-service")) {
		t.Error("Log output should contain service name")
	}
	if !bytes.Contains(buf.Bytes(), []byte("200")) {
		t.Error("Log output should contain status code")
	}
}

func TestLogError(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to error level
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	testErr := &testError{msg: "test error"}
	LogError(testErr, "test error message", map[string]interface{}{
		"field1": "value1",
	})

	output := buf.String()
	if output == "" {
		t.Error("LogError() produced no output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test error message")) {
		t.Error("Log output should contain error message")
	}
}

func TestLogWarn(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to warn level
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	LogWarn("test warning", map[string]interface{}{
		"field1": "value1",
	})

	output := buf.String()
	if output == "" {
		t.Error("LogWarn() produced no output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test warning")) {
		t.Error("Log output should contain warning message")
	}
}

func TestLogInfo(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to info level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	LogInfo("test info", map[string]interface{}{
		"field1": "value1",
	})

	output := buf.String()
	if output == "" {
		t.Error("LogInfo() produced no output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test info")) {
		t.Error("Log output should contain info message")
	}
}

func TestLogDebug(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	testLogger := zerolog.New(&buf).With().Timestamp().Logger()
	Logger = testLogger

	// Set to debug level
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	LogDebug("test debug", map[string]interface{}{
		"field1": "value1",
	})

	output := buf.String()
	if output == "" {
		t.Error("LogDebug() produced no output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test debug")) {
		t.Error("Log output should contain debug message")
	}
}

// testError is a simple error type for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
