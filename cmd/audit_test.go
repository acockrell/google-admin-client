package cmd

import (
	"testing"
	"time"
)

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name          string
		startTimeFlag string
		endTimeFlag   string
		wantErr       bool
		description   string
	}{
		{
			name:          "default time range (empty flags)",
			startTimeFlag: "",
			endTimeFlag:   "",
			wantErr:       false,
			description:   "should default to last 24 hours",
		},
		{
			name:          "valid custom time range",
			startTimeFlag: "2024-10-01T00:00:00Z",
			endTimeFlag:   "2024-10-08T00:00:00Z",
			wantErr:       false,
			description:   "should accept valid RFC3339 times",
		},
		{
			name:          "only start time provided",
			startTimeFlag: "2024-10-01T00:00:00Z",
			endTimeFlag:   "",
			wantErr:       false,
			description:   "should default end time to now",
		},
		{
			name:          "only end time provided",
			startTimeFlag: "",
			endTimeFlag:   "2024-10-08T00:00:00Z",
			wantErr:       false,
			description:   "should default start time to 24h ago",
		},
		{
			name:          "invalid start time format",
			startTimeFlag: "2024-10-01",
			endTimeFlag:   "2024-10-08T00:00:00Z",
			wantErr:       true,
			description:   "should reject non-RFC3339 format",
		},
		{
			name:          "invalid end time format",
			startTimeFlag: "2024-10-01T00:00:00Z",
			endTimeFlag:   "2024-10-08",
			wantErr:       true,
			description:   "should reject non-RFC3339 format",
		},
		{
			name:          "end time before start time",
			startTimeFlag: "2024-10-08T00:00:00Z",
			endTimeFlag:   "2024-10-01T00:00:00Z",
			wantErr:       true,
			description:   "should reject inverted time range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set package-level variables
			startTime = tt.startTimeFlag
			endTime = tt.endTimeFlag

			start, end, err := parseTimeRange()

			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimeRange() error = %v, wantErr %v (%s)", err, tt.wantErr, tt.description)
				return
			}

			if !tt.wantErr {
				// Validate that times are in RFC3339 format
				_, err := time.Parse(time.RFC3339, start)
				if err != nil {
					t.Errorf("parseTimeRange() returned invalid start time format: %v", err)
				}

				_, err = time.Parse(time.RFC3339, end)
				if err != nil {
					t.Errorf("parseTimeRange() returned invalid end time format: %v", err)
				}

				// Validate that start is before end
				startParsed, _ := time.Parse(time.RFC3339, start)
				endParsed, _ := time.Parse(time.RFC3339, end)
				if endParsed.Before(startParsed) {
					t.Errorf("parseTimeRange() returned end before start")
				}
			}

			// Reset package-level variables
			startTime = ""
			endTime = ""
		})
	}
}

func TestValidAppTypes(t *testing.T) {
	validApps := []string{
		"admin", "login", "drive", "calendar", "groups",
		"mobile", "token", "groups_enterprise", "saml",
		"chrome", "gcp", "chat", "meet",
	}

	for _, app := range validApps {
		if !validAppTypes[app] {
			t.Errorf("validAppTypes missing expected app type: %s", app)
		}
	}
}

func TestAuditCommandRegistration(t *testing.T) {
	// Check that audit command is registered
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "audit" {
			found = true

			// Check for export subcommand
			exportFound := false
			for _, subcmd := range cmd.Commands() {
				if subcmd.Name() == "export" {
					exportFound = true

					// Check required flags
					appFlag := subcmd.Flag("app")
					if appFlag == nil {
						t.Error("audit export command missing --app flag")
					}

					// Check optional flags
					flags := []string{"start-time", "end-time", "user", "event-name", "actor-ip", "output", "output-file", "max-results"}
					for _, flagName := range flags {
						if subcmd.Flag(flagName) == nil {
							t.Errorf("audit export command missing --%s flag", flagName)
						}
					}
					break
				}
			}

			if !exportFound {
				t.Error("audit command missing export subcommand")
			}
			break
		}
	}

	if !found {
		t.Error("audit command not registered with root command")
	}
}
