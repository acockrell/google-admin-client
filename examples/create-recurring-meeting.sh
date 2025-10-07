#!/bin/bash
#
# Create Recurring Meeting Script
#
# Create a recurring calendar event with customizable frequency.
#
# Usage:
#   ./create-recurring-meeting.sh <calendar-email> <title> <start> <end> [frequency] [count]

set -euo pipefail

# Check for required arguments
if [ $# -lt 4 ]; then
    echo "Usage: $0 <calendar-email> <meeting-title> <start-datetime> <end-datetime> [frequency] [count]"
    echo ""
    echo "Arguments:"
    echo "  calendar-email  - Email address of the calendar"
    echo "  meeting-title   - Title of the meeting"
    echo "  start-datetime  - Start time in RFC3339 format (e.g., 2025-10-15T09:00:00-04:00)"
    echo "  end-datetime    - End time in RFC3339 format (e.g., 2025-10-15T10:00:00-04:00)"
    echo "  frequency       - Recurrence: daily, weekly, or monthly (default: weekly)"
    echo "  count           - Number of occurrences (default: 52)"
    echo ""
    echo "Examples:"
    echo "  # Weekly team meeting for a year"
    echo "  $0 team@example.com 'Team Standup' \\"
    echo "     '2025-10-15T09:00:00-04:00' '2025-10-15T09:30:00-04:00' weekly 52"
    echo ""
    echo "  # Daily standup for 3 months"
    echo "  $0 team@example.com 'Daily Standup' \\"
    echo "     '2025-10-15T09:00:00-04:00' '2025-10-15T09:15:00-04:00' daily 90"
    echo ""
    echo "  # Monthly all-hands for a year"
    echo "  $0 company@example.com 'All Hands Meeting' \\"
    echo "     '2025-10-15T10:00:00-04:00' '2025-10-15T11:00:00-04:00' monthly 12"
    exit 1
fi

CALENDAR_EMAIL="$1"
MEETING_TITLE="$2"
START_DATE="$3"
END_DATE="$4"
FREQUENCY="${5:-weekly}"
COUNT="${6:-52}"

echo "========================================="
echo "Create Recurring Meeting"
echo "========================================="
echo "Calendar: $CALENDAR_EMAIL"
echo "Title: $MEETING_TITLE"
echo "Start: $START_DATE"
echo "End: $END_DATE"
echo "Frequency: $FREQUENCY"
echo "Occurrences: $COUNT"
echo "========================================="
echo ""

# Validate frequency
if [[ ! "$FREQUENCY" =~ ^(daily|weekly|monthly)$ ]]; then
    echo "✗ Error: Frequency must be daily, weekly, or monthly"
    exit 1
fi

# Confirm before creating
read -p "Create this recurring meeting? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled"
    exit 0
fi

# Create the recurring meeting
echo "Creating recurring meeting..."
if gac calendar create "$CALENDAR_EMAIL" \
    -s "$MEETING_TITLE" \
    -b "$START_DATE" \
    -e "$END_DATE" \
    -f "$FREQUENCY" \
    -c "$COUNT" \
    -l "Virtual - Video Conference"; then
    echo "✓ Meeting created successfully"
else
    echo "✗ Meeting creation failed"
    exit 1
fi

echo ""
echo "========================================="
echo "Meeting Created"
echo "========================================="
echo "A $FREQUENCY recurring meeting has been created"
echo "Total occurrences: $COUNT"
echo ""
echo "Next steps:"
echo "  1. Add attendees if needed (via Google Calendar UI)"
echo "  2. Add meeting link/location details"
echo "  3. Add meeting description/agenda"
echo "========================================="
