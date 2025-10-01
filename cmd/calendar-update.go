package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// flags / parameters
var (
	eventID string
)

// listCalendarCmd represents the list-calendar command
var updateCalendarCmd = &cobra.Command{
	Use:   "update",
	Short: "update calendar event",
	Run:   updateCalendarRunFunc,
	Long: `
List calendar event(s).

Usage
-----

$ gac calendar update someRandomCalendarID@example.com -i randomEventID ...
`,
}

func init() {
	calendarCmd.AddCommand(updateCalendarCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// update-profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCalendarCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "full export")
	// listUserCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "csv export")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCalendarCmd.Flags().StringVarP(&numEvents, "calendar-id", "i", 10, "calendar ID")
	updateCalendarCmd.Flags().StringSliceVarP(&eventAttendees, "attendee", "a", eventAttendees, "event attendee (multiple ok)")
	updateCalendarCmd.Flags().StringVarP(&eventDescription, "description", "d", "", "event description")
	updateCalendarCmd.Flags().StringVarP(&eventLocation, "location", "l", "The Matrix", "event location")
	updateCalendarCmd.Flags().StringVarP(&eventSummary, "summary", "s", "", "event summary/title (required)")
	updateCalendarCmd.Flags().StringVarP(&eventStart, "begin", "b", "", "event start (RFC3339 format, required)")
	updateCalendarCmd.Flags().StringVarP(&eventEnd, "end", "e", "", "event end (RFC3339 format, required)")
	updateCalendarCmd.Flags().IntVarP(&eventRecurrenceCount, "count", "c", 1, "recurrence count")
	updateCalendarCmd.Flags().StringVarP(&eventRecurrenceFreq, "frequency", "f", "daily", "recurrence frequency (daily, weekly, monthly)")
	updateCalendarCmd.Flags().StringVarP(&eventID, "event-id", "i", "", "ID of event to update (required)")
}

func updateCalendarRunFunc(cmd *cobra.Command, args []string) {
	var calendarID string

	if len(args) > 0 {
		calendarID = args[0]
	} else {
		exitWithError("must provide calendar ID")
	}

	if eventID == "" {
		exitWithError("must provide event ID")
	}

	client, err := newCalendarClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	event, err := client.Events.Get(calendarID, eventID).Do()
	if err != nil {
		exitWithError(err.Error())
	}

	collectEventInfo(event)
	buf, _ := json.MarshalIndent(event, "", "  ")
	fmt.Printf("%s\n", buf)

	e, err := client.Events.Update(calendarID, eventID, event).Do()
	if err != nil {
		exitWithError(err.Error())
	}
	fmt.Printf("event URL: %s\n", e.HtmlLink)
}
