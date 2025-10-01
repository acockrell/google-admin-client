package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	calendar "google.golang.org/api/calendar/v3"
)

// flags / parameters
var (
	numEvents int64
	timeMin   string
	timeMax   string
)

// listCalendarCmd represents the list-calendar command
var listCalendarCmd = &cobra.Command{
	Use:   "list",
	Short: "list calendar events",
	Run:   listCalendarRunFunc,
	Long: `
List calendar event(s).

Usage
-----

$ gac calendar list someRandomCalendarID@example.com
$ gac calendar list someRandomCalendarID@example.com -n 3 --time-min 2011-06-03T10:00:00-04:00
`,
}

func init() {
	calendarCmd.AddCommand(listCalendarCmd)
	// Here you will define your flags and configuration settings.
	listCalendarCmd.Flags().Int64VarP(&numEvents, "num-events", "n", 10, "number of events")
	listCalendarCmd.Flags().StringVarP(&timeMin, "time-min", "", "", "number of events")
	listCalendarCmd.Flags().StringVarP(&timeMax, "time-max", "", "", "number of events")
}

func listCalendarRunFunc(cmd *cobra.Command, args []string) {
	var calendarID string

	if len(args) > 0 {
		calendarID = args[0]
	} else {
		exitWithError("must provide calendar ID")
	}

	client, err := newCalendarClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// list events from provided calendarID
	c := &calendar.Events{}
	if timeMin != "" && timeMax != "" {
		c, err = client.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(timeMin).TimeMax(timeMax).MaxResults(numEvents).OrderBy("startTime").Do()
	} else if timeMin != "" {
		c, err = client.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMin(timeMin).MaxResults(numEvents).OrderBy("startTime").Do()
	} else if timeMax != "" {
		c, err = client.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).TimeMax(timeMax).MaxResults(numEvents).OrderBy("startTime").Do()
	} else {
		c, err = client.Events.List(calendarID).ShowDeleted(false).SingleEvents(true).MaxResults(numEvents).OrderBy("startTime").Do()
	}
	if err != nil {
		exitWithError(err.Error())
	}
	buf, _ := json.MarshalIndent(c, "", "  ")
	fmt.Printf("%s\n", buf)
}
