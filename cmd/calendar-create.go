package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	calendar "google.golang.org/api/calendar/v3"
)

// flags / parameters
var (
	eventAttendees       = []string{}
	eventDescription     string
	eventLocation        string
	eventSummary         string
	eventStart           string
	eventEnd             string
	eventRecurrenceCount int
	eventRecurrenceFreq  string
)

// listCalendarCmd represents the list-calendar command
var createCalendarCmd = &cobra.Command{
	Use:   "create",
	Short: "create calendar event",
	Run:   createCalendarRunFunc,
	Long: `
Create calendar event.

Usage
-----

$ gac calendar create someRandomCalendarID@example.com -s "all day event" -b 1970-01-01 -e 1970-01-02
$ gac calendar create someRandomCalendarID@example.com -s "short event" -d "with a long description" -b 1970-01-01T00:00:00-04:00 -e 1970-01-01T01:00:00-04:00
`,
}

func init() {
	calendarCmd.AddCommand(createCalendarCmd)

	// Here you will define your flags and configuration settings.
	createCalendarCmd.Flags().StringSliceVarP(&eventAttendees, "attendee", "a", eventAttendees, "event attendee (multiple ok)")
	createCalendarCmd.Flags().StringVarP(&eventDescription, "description", "d", "", "event description")
	createCalendarCmd.Flags().StringVarP(&eventLocation, "location", "l", "The Matrix", "event location")
	createCalendarCmd.Flags().StringVarP(&eventSummary, "summary", "s", "", "event summary/title (required)")
	createCalendarCmd.Flags().StringVarP(&eventStart, "begin", "b", "", "event start (RFC3339 format, required)")
	createCalendarCmd.Flags().StringVarP(&eventEnd, "end", "e", "", "event end (RFC3339 format, required)")
	createCalendarCmd.Flags().IntVarP(&eventRecurrenceCount, "count", "c", 1, "recurrence count")
	createCalendarCmd.Flags().StringVarP(&eventRecurrenceFreq, "frequency", "f", "daily", "recurrence frequency (daily, weekly, monthly)")
}

func createCalendarRunFunc(cmd *cobra.Command, args []string) {
	var calendarID string

	if len(args) > 0 {
		calendarID = args[0]
	} else {
		exitWithError("must provide calendar ID")
	}

	if eventSummary == "" || eventStart == "" || eventEnd == "" {
		exitWithError("--summary, --begin and --end required.")
	}

	client, err := newCalendarClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	event := calendar.Event{}
	collectEventInfo(&event)

	buf, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("%s\n", buf)

	e, err := client.Events.Insert(calendarID, &event).Do()
	if err != nil {
		exitWithError(err.Error())
	}
	fmt.Printf("event URL: %s\n", e.HtmlLink)
}
