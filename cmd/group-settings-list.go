package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	groupSettingsListFormat string
)

// groupSettingsListCmd represents the group-settings list command
var groupSettingsListCmd = &cobra.Command{
	Use:   "list [group-email]",
	Short: "View group settings",
	Long: `View settings for a Google Workspace group.

Usage
-----

$ gac group-settings list operations@example.com
$ gac group-settings list engineering
$ gac group-settings list sales@example.com --format json

Description
-----------

Displays the settings for a group, including:
- Who can join the group
- Who can view group messages
- Who can post messages
- Message moderation settings
- Email delivery preferences
- Archive settings
- Custom footer text

The --format flag controls output format:
  table - Human-readable table format (default)
  json  - JSON format for scripting

If you don't include the @domain part in the group email, the configured domain
will be automatically appended.
`,
	Args: cobra.ExactArgs(1),
	RunE: groupSettingsListRunFunc,
}

func init() {
	groupSettingsCmd.AddCommand(groupSettingsListCmd)
	groupSettingsListCmd.Flags().StringVarP(&groupSettingsListFormat, "format", "f", "table", "output format: table or json")
}

func groupSettingsListRunFunc(cmd *cobra.Command, args []string) error {
	groupEmail := args[0]
	if !strings.Contains(groupEmail, "@") {
		groupEmail = groupEmail + "@" + getDomain()
	}

	// Validate email
	if err := ValidateEmail(groupEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid group email: %v\n", err)
		return err
	}

	client, err := newGroupsSettingsClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	settings, err := client.Groups.Get(groupEmail).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting group settings for %s: %v\n", groupEmail, err)
		return err
	}

	if groupSettingsListFormat == "json" {
		buf, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
			return err
		}
		fmt.Printf("%s\n", buf)
	} else {
		// Table format
		displayGroupSettings(settings)
	}

	return nil
}

func displayGroupSettings(settings interface{}) {
	// Type assert to map for flexible access
	settingsMap, ok := settings.(map[string]interface{})
	if !ok {
		// If not a map, try to marshal and unmarshal to get a map
		bytes, _ := json.Marshal(settings)
		if err := json.Unmarshal(bytes, &settingsMap); err != nil {
			fmt.Fprintf(os.Stderr, "Error converting settings to map: %v\n", err)
			return
		}
	}

	fmt.Printf("Group Settings\n")
	fmt.Printf("==============\n\n")

	// Basic Information
	printSetting("Email", settingsMap["email"])
	printSetting("Name", settingsMap["name"])
	printSetting("Description", settingsMap["description"])
	fmt.Println()

	// Access Settings
	fmt.Printf("Access Settings:\n")
	fmt.Printf("----------------\n")
	printSetting("  Who Can Join", settingsMap["whoCanJoin"])
	printSetting("  Who Can View Group", settingsMap["whoCanViewGroup"])
	printSetting("  Who Can View Membership", settingsMap["whoCanViewMembership"])
	printSetting("  Allow External Members", settingsMap["allowExternalMembers"])
	fmt.Println()

	// Posting Settings
	fmt.Printf("Posting Settings:\n")
	fmt.Printf("-----------------\n")
	printSetting("  Who Can Post Message", settingsMap["whoCanPostMessage"])
	printSetting("  Allow Web Posting", settingsMap["allowWebPosting"])
	printSetting("  Message Moderation Level", settingsMap["messageModerationLevel"])
	printSetting("  Spam Moderation Level", settingsMap["spamModerationLevel"])
	fmt.Println()

	// Email Settings
	fmt.Printf("Email Settings:\n")
	fmt.Printf("---------------\n")
	printSetting("  Send Message Deny Notification", settingsMap["sendMessageDenyNotification"])
	printSetting("  Reply To", settingsMap["replyTo"])
	printSetting("  Custom Reply To", settingsMap["customReplyTo"])
	printSetting("  Include Custom Footer", settingsMap["includeCustomFooter"])
	if customFooter, ok := settingsMap["customFooterText"].(string); ok && customFooter != "" {
		printSetting("  Custom Footer Text", customFooter)
	}
	printSetting("  Include in Global Address List", settingsMap["includeInGlobalAddressList"])
	fmt.Println()

	// Moderation Settings
	fmt.Printf("Moderation Settings:\n")
	fmt.Printf("--------------------\n")
	printSetting("  Who Can Contact Owner", settingsMap["whoCanContactOwner"])
	printSetting("  Who Can Moderate Members", settingsMap["whoCanModerateMembers"])
	printSetting("  Who Can Moderate Content", settingsMap["whoCanModerateContent"])
	fmt.Println()

	// Archive Settings
	fmt.Printf("Archive Settings:\n")
	fmt.Printf("-----------------\n")
	printSetting("  Archive Only", settingsMap["archiveOnly"])
	printSetting("  Message Display Font", settingsMap["messageDisplayFont"])
	printSetting("  Show in Group Directory", settingsMap["showInGroupDirectory"])
	printSetting("  Max Message Bytes", settingsMap["maxMessageBytes"])
	printSetting("  Is Archived", settingsMap["isArchived"])
	fmt.Println()

	// Member Settings
	fmt.Printf("Member Settings:\n")
	fmt.Printf("----------------\n")
	printSetting("  Who Can Leave Group", settingsMap["whoCanLeaveGroup"])
	printSetting("  Who Can Add", settingsMap["whoCanAdd"])
	printSetting("  Who Can Invite", settingsMap["whoCanInvite"])
	printSetting("  Who Can Approve Members", settingsMap["whoCanApproveMembers"])
	printSetting("  Who Can Ban Users", settingsMap["whoCanBanUsers"])
	printSetting("  Allow Google Communication", settingsMap["allowGoogleCommunication"])
	printSetting("  Members Can Post As The Group", settingsMap["membersCanPostAsTheGroup"])
	fmt.Println()
}

func printSetting(label string, value interface{}) {
	if value == nil {
		return
	}

	// Convert value to string
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case bool:
		if v {
			strValue = "true"
		} else {
			strValue = "false"
		}
	case float64, int, int64:
		strValue = fmt.Sprintf("%v", v)
	default:
		strValue = fmt.Sprintf("%v", v)
	}

	if strValue != "" && strValue != "0" {
		fmt.Printf("%-35s: %s\n", label, strValue)
	}
}
