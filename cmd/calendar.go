package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	calendar "google.golang.org/api/calendar/v3"
)

// userCmd represents the user command
var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "List & modify calendars",
}

func init() {
	rootCmd.AddCommand(calendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func collectEventInfo(event *calendar.Event) {
	event.Summary = eventSummary

	if len(eventStart) > 10 {
		// create timed event
		event.Start = &calendar.EventDateTime{
			DateTime: eventStart,
		}
		event.End = &calendar.EventDateTime{
			DateTime: eventEnd,
		}
	} else if len(eventStart) > 0 {
		// create all day event
		event.Start = &calendar.EventDateTime{
			Date: eventStart,
		}
		event.End = &calendar.EventDateTime{
			Date: eventEnd,
		}
	} else if len(event.Start.DateTime) > 0 {
		// update timed event
		event.Start = &calendar.EventDateTime{
			DateTime: event.Start.DateTime,
		}
		event.End = &calendar.EventDateTime{
			DateTime: event.End.DateTime,
		}
	} else {
		// update all day event
		event.Start = &calendar.EventDateTime{
			DateTime: event.Start.Date,
		}
		event.End = &calendar.EventDateTime{
			DateTime: event.End.Date,
		}

	}

	if len(eventAttendees) > 0 {
		event.Attendees = []*calendar.EventAttendee{}
		for _, attendee := range eventAttendees {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{
				Email: attendee,
			})
		}
	}

	if eventLocation != "" {
		event.Location = eventLocation
	}

	if eventDescription != "" {
		event.Description = eventDescription
	}

	if eventRecurrenceCount != 1 {
		r := fmt.Sprintf("RRULE:FREQ=%s;COUNT=%d", strings.ToUpper(eventRecurrenceFreq), eventRecurrenceCount)
		event.Recurrence = []string{r}
	}
}
