package cmd

import (
	"testing"

	calendar "google.golang.org/api/calendar/v3"
)

func TestCollectEventInfo(t *testing.T) {
	// Save original values
	origSummary := eventSummary
	origStart := eventStart
	origEnd := eventEnd
	origAttendees := eventAttendees
	origLocation := eventLocation
	origDescription := eventDescription
	origRecurrenceCount := eventRecurrenceCount
	origRecurrenceFreq := eventRecurrenceFreq

	defer func() {
		eventSummary = origSummary
		eventStart = origStart
		eventEnd = origEnd
		eventAttendees = origAttendees
		eventLocation = origLocation
		eventDescription = origDescription
		eventRecurrenceCount = origRecurrenceCount
		eventRecurrenceFreq = origRecurrenceFreq
	}()

	tests := []struct {
		name            string
		summary         string
		start           string
		end             string
		attendees       []string
		location        string
		description     string
		recurrenceCount int
		recurrenceFreq  string
	}{
		{
			name:        "timed event",
			summary:     "Meeting",
			start:       "2024-01-15T10:00:00-05:00",
			end:         "2024-01-15T11:00:00-05:00",
			attendees:   []string{"attendee@example.com"},
			location:    "Conference Room",
			description: "Important meeting",
		},
		{
			name:     "all day event",
			summary:  "Holiday",
			start:    "2024-01-15",
			end:      "2024-01-16",
			location: "Office",
		},
		{
			name:            "recurring event",
			summary:         "Weekly Standup",
			start:           "2024-01-15T09:00:00-05:00",
			end:             "2024-01-15T09:30:00-05:00",
			recurrenceCount: 10,
			recurrenceFreq:  "weekly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventSummary = tt.summary
			eventStart = tt.start
			eventEnd = tt.end
			eventAttendees = tt.attendees
			eventLocation = tt.location
			eventDescription = tt.description
			eventRecurrenceCount = tt.recurrenceCount
			eventRecurrenceFreq = tt.recurrenceFreq

			event := &calendar.Event{}
			collectEventInfo(event)

			if event.Summary != tt.summary {
				t.Errorf("Summary = %q, want %q", event.Summary, tt.summary)
			}

			if len(tt.start) > 10 {
				if event.Start.DateTime != tt.start {
					t.Errorf("Start.DateTime = %q, want %q", event.Start.DateTime, tt.start)
				}
			} else if len(tt.start) > 0 {
				if event.Start.Date != tt.start {
					t.Errorf("Start.Date = %q, want %q", event.Start.Date, tt.start)
				}
			}

			if tt.location != "" && event.Location != tt.location {
				t.Errorf("Location = %q, want %q", event.Location, tt.location)
			}

			if tt.description != "" && event.Description != tt.description {
				t.Errorf("Description = %q, want %q", event.Description, tt.description)
			}

			if len(tt.attendees) > 0 {
				if len(event.Attendees) != len(tt.attendees) {
					t.Errorf("Attendees count = %d, want %d", len(event.Attendees), len(tt.attendees))
				}
			}

			if tt.recurrenceCount > 1 && len(event.Recurrence) == 0 {
				t.Error("Expected recurrence rules but got none")
			}
		})
	}
}

func TestCollectEventInfoUpdateExisting(t *testing.T) {
	// Save original values
	origSummary := eventSummary
	origStart := eventStart
	origEnd := eventEnd

	defer func() {
		eventSummary = origSummary
		eventStart = origStart
		eventEnd = origEnd
	}()

	t.Run("update existing timed event", func(t *testing.T) {
		eventSummary = "Updated Meeting"
		eventStart = ""
		eventEnd = ""

		event := &calendar.Event{
			Start: &calendar.EventDateTime{
				DateTime: "2024-01-15T10:00:00-05:00",
			},
			End: &calendar.EventDateTime{
				DateTime: "2024-01-15T11:00:00-05:00",
			},
		}

		collectEventInfo(event)

		if event.Summary != "Updated Meeting" {
			t.Errorf("Summary = %q, want %q", event.Summary, "Updated Meeting")
		}
		if event.Start.DateTime != "2024-01-15T10:00:00-05:00" {
			t.Errorf("Start.DateTime preserved = %q", event.Start.DateTime)
		}
	})

	t.Run("update existing all day event", func(t *testing.T) {
		eventSummary = "Updated Holiday"
		eventStart = ""
		eventEnd = ""

		event := &calendar.Event{
			Start: &calendar.EventDateTime{
				Date: "2024-01-15",
			},
			End: &calendar.EventDateTime{
				Date: "2024-01-16",
			},
		}

		collectEventInfo(event)

		if event.Summary != "Updated Holiday" {
			t.Errorf("Summary = %q, want %q", event.Summary, "Updated Holiday")
		}
		if event.Start.DateTime != "2024-01-15" {
			t.Errorf("Start.DateTime = %q", event.Start.DateTime)
		}
	})
}
