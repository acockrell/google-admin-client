package cmd

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger is the global logger instance
	Logger zerolog.Logger
)

func init() {
	// Initialize with a default logger (console output, info level)
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}).With().Timestamp().Logger()

	// Set default level to Info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// InitLogger initializes the logger based on configuration
func InitLogger(verbose bool, logLevel string, jsonFormat bool) {
	// Determine log level
	var level zerolog.Level
	if verbose {
		level = zerolog.DebugLevel
	} else {
		level = parseLogLevel(logLevel)
	}

	zerolog.SetGlobalLevel(level)

	// Determine output format
	var output io.Writer = os.Stderr
	if jsonFormat {
		// JSON output for scripting/automation
		Logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		// Human-friendly console output
		Logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}).With().Timestamp().Logger()
	}

	// Update global logger
	log.Logger = Logger
}

// parseLogLevel converts string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

// LogAPICall logs an API call with debug level
func LogAPICall(service, method string, params map[string]interface{}) {
	Logger.Debug().
		Str("service", service).
		Str("method", method).
		Interface("params", params).
		Msg("API call")
}

// LogAPIResponse logs an API response with debug level
func LogAPIResponse(service, method string, statusCode int, duration time.Duration) {
	Logger.Debug().
		Str("service", service).
		Str("method", method).
		Int("status_code", statusCode).
		Dur("duration", duration).
		Msg("API response")
}

// LogError logs an error with context
func LogError(err error, msg string, fields map[string]interface{}) {
	event := Logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// LogWarn logs a warning with context
func LogWarn(msg string, fields map[string]interface{}) {
	event := Logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// LogInfo logs an info message with context
func LogInfo(msg string, fields map[string]interface{}) {
	event := Logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// LogDebug logs a debug message with context
func LogDebug(msg string, fields map[string]interface{}) {
	event := Logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
