package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	reports "google.golang.org/api/admin/reports/v1"
)

// flags / parameters for audit export
var (
	appType           string
	startTime         string
	endTime           string
	userEmail         string
	eventNames        []string
	auditOutputFormat string
	outputFile        string
	maxResults        int64
	actorIPAddr       string
)

// validAppTypes lists all supported application types for audit logs
var validAppTypes = map[string]bool{
	"admin":             true,
	"login":             true,
	"drive":             true,
	"calendar":          true,
	"groups":            true,
	"mobile":            true,
	"token":             true,
	"groups_enterprise": true,
	"saml":              true,
	"chrome":            true,
	"gcp":               true,
	"chat":              true,
	"meet":              true,
}

// auditExportCmd represents the audit export command
var auditExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export audit logs for a specific application",
	Long: `Export Google Workspace audit logs for a specific application.

The export command retrieves audit logs from the Google Admin Reports API and outputs them
in JSON or CSV format. You can filter by time range, user, event type, and more.

Application Types:
  admin              Admin console activities (user creation, settings changes, etc.)
  login              User login and logout activities
  drive              Google Drive file operations (view, share, download, etc.)
  calendar           Calendar event activities
  groups             Google Groups activities
  mobile             Mobile device activities
  token              OAuth token activities
  groups_enterprise  Groups for Enterprise activities
  saml               SAML authentication activities
  chrome             Chrome browser activities
  gcp                Google Cloud Platform activities
  chat               Google Chat activities
  meet               Google Meet activities

Time Range:
  By default, exports the last 24 hours of audit logs.
  Use --start-time and --end-time to specify a custom range.
  Dates must be in RFC3339 format (e.g., 2024-10-08T00:00:00Z).

Output Formats:
  json               Full event details in JSON format (default)
  csv                Tabular format with key fields

Examples:
  # Export last 24h of admin console activities
  gac audit export --app admin

  # Export login activities for specific user
  gac audit export --app login --user user@example.com

  # Export drive activities with custom time range
  gac audit export --app drive \
    --start-time 2024-10-01T00:00:00Z \
    --end-time 2024-10-08T00:00:00Z

  # Export to CSV file
  gac audit export --app admin --output csv --output-file admin-audit.csv

  # Filter by specific event types
  gac audit export --app admin --event-name USER_CREATED --event-name GROUP_CREATED

  # Filter by IP address
  gac audit export --app login --actor-ip 192.168.1.100
`,
	Run: auditExportRunFunc,
}

func init() {
	auditCmd.AddCommand(auditExportCmd)

	// Required flags
	auditExportCmd.Flags().StringVar(&appType, "app", "", "application type (required: admin, login, drive, calendar, groups, mobile, token, etc.)")
	if err := auditExportCmd.MarkFlagRequired("app"); err != nil {
		Logger.Error().Err(err).Msg("Failed to mark app flag as required")
	}

	// Time range flags
	auditExportCmd.Flags().StringVar(&startTime, "start-time", "", "start time in RFC3339 format (default: 24 hours ago)")
	auditExportCmd.Flags().StringVar(&endTime, "end-time", "", "end time in RFC3339 format (default: now)")

	// Filter flags
	auditExportCmd.Flags().StringVar(&userEmail, "user", "", "filter by user email address")
	auditExportCmd.Flags().StringSliceVar(&eventNames, "event-name", []string{}, "filter by event name (can specify multiple)")
	auditExportCmd.Flags().StringVar(&actorIPAddr, "actor-ip", "", "filter by actor IP address")

	// Output flags
	auditExportCmd.Flags().StringVarP(&auditOutputFormat, "output", "o", "json", "output format (json, csv)")
	auditExportCmd.Flags().StringVarP(&outputFile, "output-file", "f", "", "output file path (default: stdout)")
	auditExportCmd.Flags().Int64Var(&maxResults, "max-results", 0, "maximum number of results to return (default: all)")
}

// parseTimeRange parses and validates start and end times
func parseTimeRange() (string, string, error) {
	var start, end time.Time
	var err error

	// Parse end time first to determine proper defaults
	if endTime == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return "", "", fmt.Errorf("invalid end-time format (expected RFC3339): %w", err)
		}
	}

	// Parse start time (default: 24 hours before end time)
	if startTime == "" {
		start = end.Add(-24 * time.Hour)
	} else {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			return "", "", fmt.Errorf("invalid start-time format (expected RFC3339): %w", err)
		}
	}

	// Validate time range
	if end.Before(start) {
		return "", "", fmt.Errorf("end-time must be after start-time")
	}

	return start.Format(time.RFC3339), end.Format(time.RFC3339), nil
}

// writeJSONOutput writes activities to JSON format
func writeJSONOutput(activities []*reports.Activity, filename string) error {
	var output *os.File
	var err error

	if filename == "" {
		output = os.Stdout
	} else {
		// Validate output file path to prevent directory traversal
		if err := validateCredentialPath(filename); err != nil {
			return fmt.Errorf("invalid output file path: %w", err)
		}

		// #nosec G304 - Path is validated by validateCredentialPath() above
		output, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if cerr := output.Close(); cerr != nil && err == nil {
				err = fmt.Errorf("failed to close output file: %w", cerr)
			}
		}()
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(activities); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// writeCSVOutput writes activities to CSV format
func writeCSVOutput(activities []*reports.Activity, filename string) error {
	var output *os.File
	var err error

	if filename == "" {
		output = os.Stdout
	} else {
		// Validate output file path to prevent directory traversal
		if err := validateCredentialPath(filename); err != nil {
			return fmt.Errorf("invalid output file path: %w", err)
		}

		// #nosec G304 - Path is validated by validateCredentialPath() above
		output, err = os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if cerr := output.Close(); cerr != nil && err == nil {
				err = fmt.Errorf("failed to close output file: %w", cerr)
			}
		}()
	}

	writer := csv.NewWriter(output)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Timestamp", "Actor", "Event", "IP Address", "Application"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write activity rows
	for _, activity := range activities {
		timestamp := ""
		if activity.Id != nil && activity.Id.Time != "" {
			t, err := time.Parse(time.RFC3339, activity.Id.Time)
			if err == nil {
				timestamp = t.Format("2006-01-02 15:04:05 MST")
			}
		}

		actor := ""
		if activity.Actor != nil && activity.Actor.Email != "" {
			actor = activity.Actor.Email
		}

		eventName := ""
		if len(activity.Events) > 0 && activity.Events[0].Name != "" {
			eventName = activity.Events[0].Name
		}

		ipAddress := ""
		if activity.IpAddress != "" {
			ipAddress = activity.IpAddress
		}

		appName := ""
		if activity.Id != nil && activity.Id.ApplicationName != "" {
			appName = activity.Id.ApplicationName
		}

		row := []string{timestamp, actor, eventName, ipAddress, appName}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

func auditExportRunFunc(cmd *cobra.Command, args []string) {
	// Validate application type
	appType = strings.ToLower(appType)
	if !validAppTypes[appType] {
		Logger.Error().
			Str("app", appType).
			Msg("Invalid application type")
		Logger.Info().Msg("Valid application types: admin, login, drive, calendar, groups, mobile, token, groups_enterprise, saml, chrome, gcp, chat, meet")
		os.Exit(1)
	}

	// Validate output format
	auditOutputFormat = strings.ToLower(auditOutputFormat)
	if auditOutputFormat != "json" && auditOutputFormat != "csv" {
		Logger.Error().
			Str("format", auditOutputFormat).
			Msg("Invalid output format (must be json or csv)")
		os.Exit(1)
	}

	// Parse and validate time range
	startTimeStr, endTimeStr, err := parseTimeRange()
	if err != nil {
		Logger.Error().Err(err).Msg("Failed to parse time range")
		os.Exit(1)
	}

	Logger.Info().
		Str("app", appType).
		Str("start_time", startTimeStr).
		Str("end_time", endTimeStr).
		Msg("Exporting audit logs")

	// Initialize Reports API client
	client, err := newReportsClient()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Failed to initialize reports client")
	}

	// Build the API call
	// For Reports API, userKey can be "all" or a specific user email
	userKey := "all"
	if userEmail != "" {
		userKey = userEmail
		Logger.Debug().Str("user", userEmail).Msg("Filtering by user")
	}

	call := client.Activities.List(userKey, appType).
		StartTime(startTimeStr).
		EndTime(endTimeStr)

	// Apply additional filters

	if actorIPAddr != "" {
		call = call.ActorIpAddress(actorIPAddr)
		Logger.Debug().Str("ip", actorIPAddr).Msg("Filtering by IP address")
	}

	if len(eventNames) > 0 {
		call = call.EventName(strings.Join(eventNames, ","))
		Logger.Debug().Strs("events", eventNames).Msg("Filtering by event names")
	}

	if maxResults > 0 {
		call = call.MaxResults(maxResults)
		Logger.Debug().Int64("max_results", maxResults).Msg("Limiting results")
	}

	// Fetch activities with pagination
	var activities []*reports.Activity
	pageToken := ""
	pageCount := 0

	for {
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		Logger.Debug().
			Int("page", pageCount+1).
			Msg("Fetching audit log page")

		resp, err := call.Do()
		if err != nil {
			Logger.Error().
				Err(err).
				Int("page", pageCount+1).
				Msg("Failed to fetch activities")
			os.Exit(1)
		}

		if resp.Items != nil {
			activities = append(activities, resp.Items...)
			Logger.Debug().
				Int("count", len(resp.Items)).
				Int("total", len(activities)).
				Msg("Retrieved activities")
		}

		pageCount++

		// Check if we've reached max results
		if maxResults > 0 && int64(len(activities)) >= maxResults {
			activities = activities[:maxResults]
			break
		}

		// Check for next page
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	Logger.Info().
		Int("total", len(activities)).
		Int("pages", pageCount).
		Msg("Retrieved audit logs")

	// Output results
	if len(activities) == 0 {
		Logger.Warn().Msg("No audit logs found for the specified criteria")
		return
	}

	var outputErr error
	if auditOutputFormat == "csv" {
		outputErr = writeCSVOutput(activities, outputFile)
	} else {
		outputErr = writeJSONOutput(activities, outputFile)
	}

	if outputErr != nil {
		Logger.Error().Err(outputErr).Msg("Failed to write output")
		os.Exit(1)
	}

	if outputFile != "" {
		Logger.Info().
			Str("file", outputFile).
			Str("format", auditOutputFormat).
			Msg("Audit logs exported successfully")
	}
}
