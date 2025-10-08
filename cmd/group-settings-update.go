package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	groupssettings "google.golang.org/api/groupssettings/v1"
)

var (
	whoCanJoin                  string
	whoCanViewGroup             string
	whoCanViewMembership        string
	whoCanPostMessage           string
	whoCanContactOwner          string
	allowExternalMembers        string
	allowWebPosting             string
	archiveOnly                 string
	messageModerationLevel      string
	spamModerationLevel         string
	replyTo                     string
	customReplyTo               string
	includeCustomFooter         string
	customFooterText            string
	sendMessageDenyNotification string
	includeInGlobalAddressList  string
	showInGroupDirectory        string
	allowGoogleCommunication    string
	membersCanPostAsTheGroup    string
	whoCanLeaveGroup            string
	whoCanAdd                   string
	whoCanInvite                string
	whoCanApproveMembers        string
	whoCanModerateMembers       string
	whoCanModerateContent       string
	whoCanBanUsers              string
)

// groupSettingsUpdateCmd represents the group-settings update command
var groupSettingsUpdateCmd = &cobra.Command{
	Use:   "update [group-email]",
	Short: "Update group settings",
	Long: `Update settings for a Google Workspace group.

Usage
-----

$ gac group-settings update operations@example.com --who-can-join ALL_IN_DOMAIN_CAN_JOIN
$ gac group-settings update engineering --allow-external-members false
$ gac group-settings update sales@example.com --who-can-post-message ALL_MEMBERS_CAN_POST
$ gac group-settings update support@example.com --custom-footer-text "For help, contact support@example.com"

Description
-----------

Updates group settings. You can specify one or more flags to update specific settings.
Only the settings you specify will be updated; other settings will remain unchanged.

Common Settings:

Access Control:
  --who-can-join              Who can join the group
                              Values: CAN_REQUEST_TO_JOIN, ALL_IN_DOMAIN_CAN_JOIN,
                                      ANYONE_CAN_JOIN, INVITED_CAN_JOIN
  --who-can-view-group        Who can view group messages
                              Values: ANYONE_CAN_VIEW, ALL_IN_DOMAIN_CAN_VIEW,
                                      ALL_MEMBERS_CAN_VIEW, ALL_MANAGERS_CAN_VIEW
  --who-can-view-membership   Who can view the member list
  --allow-external-members    Allow external members (true/false)

Posting Permissions:
  --who-can-post-message      Who can post messages
                              Values: NONE_CAN_POST, ALL_MANAGERS_CAN_POST,
                                      ALL_MEMBERS_CAN_POST, ALL_IN_DOMAIN_CAN_POST,
                                      ANYONE_CAN_POST
  --allow-web-posting         Allow posting from web (true/false)
  --message-moderation-level  Message moderation level
                              Values: MODERATE_ALL_MESSAGES, MODERATE_NON_MEMBERS,
                                      MODERATE_NEW_MEMBERS, MODERATE_NONE

Email Settings:
  --reply-to                  Reply-to setting
                              Values: REPLY_TO_CUSTOM, REPLY_TO_SENDER,
                                      REPLY_TO_LIST, REPLY_TO_OWNER, REPLY_TO_IGNORE
  --custom-reply-to           Custom reply-to email address
  --custom-footer-text        Custom footer text for messages
  --include-custom-footer     Include custom footer (true/false)

Archive Settings:
  --archive-only              Make group archive-only (true/false)
  --show-in-group-directory   Show in group directory (true/false)

If you don't include the @domain part in the group email, the configured domain
will be automatically appended.
`,
	Args: cobra.ExactArgs(1),
	RunE: groupSettingsUpdateRunFunc,
}

func init() {
	groupSettingsCmd.AddCommand(groupSettingsUpdateCmd)

	// Access control flags
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanJoin, "who-can-join", "", "who can join the group")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanViewGroup, "who-can-view-group", "", "who can view group messages")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanViewMembership, "who-can-view-membership", "", "who can view membership")
	groupSettingsUpdateCmd.Flags().StringVar(&allowExternalMembers, "allow-external-members", "", "allow external members (true/false)")

	// Posting permission flags
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanPostMessage, "who-can-post-message", "", "who can post messages")
	groupSettingsUpdateCmd.Flags().StringVar(&allowWebPosting, "allow-web-posting", "", "allow web posting (true/false)")
	groupSettingsUpdateCmd.Flags().StringVar(&messageModerationLevel, "message-moderation-level", "", "message moderation level")
	groupSettingsUpdateCmd.Flags().StringVar(&spamModerationLevel, "spam-moderation-level", "", "spam moderation level")

	// Email settings flags
	groupSettingsUpdateCmd.Flags().StringVar(&replyTo, "reply-to", "", "reply-to setting")
	groupSettingsUpdateCmd.Flags().StringVar(&customReplyTo, "custom-reply-to", "", "custom reply-to email")
	groupSettingsUpdateCmd.Flags().StringVar(&customFooterText, "custom-footer-text", "", "custom footer text")
	groupSettingsUpdateCmd.Flags().StringVar(&includeCustomFooter, "include-custom-footer", "", "include custom footer (true/false)")
	groupSettingsUpdateCmd.Flags().StringVar(&sendMessageDenyNotification, "send-message-deny-notification", "", "send message deny notification (true/false)")
	groupSettingsUpdateCmd.Flags().StringVar(&includeInGlobalAddressList, "include-in-global-address-list", "", "include in global address list (true/false)")

	// Archive settings flags
	groupSettingsUpdateCmd.Flags().StringVar(&archiveOnly, "archive-only", "", "archive only (true/false)")
	groupSettingsUpdateCmd.Flags().StringVar(&showInGroupDirectory, "show-in-group-directory", "", "show in group directory (true/false)")

	// Member management flags
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanLeaveGroup, "who-can-leave-group", "", "who can leave the group")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanAdd, "who-can-add", "", "who can add members")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanInvite, "who-can-invite", "", "who can invite members")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanApproveMembers, "who-can-approve-members", "", "who can approve members")
	groupSettingsUpdateCmd.Flags().StringVar(&allowGoogleCommunication, "allow-google-communication", "", "allow Google communication (true/false)")
	groupSettingsUpdateCmd.Flags().StringVar(&membersCanPostAsTheGroup, "members-can-post-as-the-group", "", "members can post as the group (true/false)")

	// Moderation flags
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanContactOwner, "who-can-contact-owner", "", "who can contact owner")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanModerateMembers, "who-can-moderate-members", "", "who can moderate members")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanModerateContent, "who-can-moderate-content", "", "who can moderate content")
	groupSettingsUpdateCmd.Flags().StringVar(&whoCanBanUsers, "who-can-ban-users", "", "who can ban users")
}

func groupSettingsUpdateRunFunc(cmd *cobra.Command, args []string) error {
	groupEmail := args[0]
	if !strings.Contains(groupEmail, "@") {
		groupEmail = groupEmail + "@" + getDomain()
	}

	// Validate email
	if err := ValidateEmail(groupEmail); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid group email: %v\n", err)
		return err
	}

	// Validate custom reply-to email if provided
	if customReplyTo != "" {
		if err := ValidateEmail(customReplyTo); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid custom reply-to email: %v\n", err)
			return err
		}
	}

	client, err := newGroupsSettingsClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		return err
	}

	// Build the update object with only the fields that were specified
	groups := &groupssettings.Groups{
		ForceSendFields: []string{}, // Track which fields to send
	}

	// Track if any settings were provided
	hasUpdates := false

	// Access control settings
	if cmd.Flags().Changed("who-can-join") {
		groups.WhoCanJoin = whoCanJoin
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanJoin")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-view-group") {
		groups.WhoCanViewGroup = whoCanViewGroup
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanViewGroup")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-view-membership") {
		groups.WhoCanViewMembership = whoCanViewMembership
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanViewMembership")
		hasUpdates = true
	}
	if cmd.Flags().Changed("allow-external-members") {
		groups.AllowExternalMembers = allowExternalMembers
		groups.ForceSendFields = append(groups.ForceSendFields, "AllowExternalMembers")
		hasUpdates = true
	}

	// Posting permission settings
	if cmd.Flags().Changed("who-can-post-message") {
		groups.WhoCanPostMessage = whoCanPostMessage
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanPostMessage")
		hasUpdates = true
	}
	if cmd.Flags().Changed("allow-web-posting") {
		groups.AllowWebPosting = allowWebPosting
		groups.ForceSendFields = append(groups.ForceSendFields, "AllowWebPosting")
		hasUpdates = true
	}
	if cmd.Flags().Changed("message-moderation-level") {
		groups.MessageModerationLevel = messageModerationLevel
		groups.ForceSendFields = append(groups.ForceSendFields, "MessageModerationLevel")
		hasUpdates = true
	}
	if cmd.Flags().Changed("spam-moderation-level") {
		groups.SpamModerationLevel = spamModerationLevel
		groups.ForceSendFields = append(groups.ForceSendFields, "SpamModerationLevel")
		hasUpdates = true
	}

	// Email settings
	if cmd.Flags().Changed("reply-to") {
		groups.ReplyTo = replyTo
		groups.ForceSendFields = append(groups.ForceSendFields, "ReplyTo")
		hasUpdates = true
	}
	if cmd.Flags().Changed("custom-reply-to") {
		groups.CustomReplyTo = customReplyTo
		groups.ForceSendFields = append(groups.ForceSendFields, "CustomReplyTo")
		hasUpdates = true
	}
	if cmd.Flags().Changed("custom-footer-text") {
		groups.CustomFooterText = customFooterText
		groups.ForceSendFields = append(groups.ForceSendFields, "CustomFooterText")
		hasUpdates = true
	}
	if cmd.Flags().Changed("include-custom-footer") {
		groups.IncludeCustomFooter = includeCustomFooter
		groups.ForceSendFields = append(groups.ForceSendFields, "IncludeCustomFooter")
		hasUpdates = true
	}
	if cmd.Flags().Changed("send-message-deny-notification") {
		groups.SendMessageDenyNotification = sendMessageDenyNotification
		groups.ForceSendFields = append(groups.ForceSendFields, "SendMessageDenyNotification")
		hasUpdates = true
	}
	if cmd.Flags().Changed("include-in-global-address-list") {
		groups.IncludeInGlobalAddressList = includeInGlobalAddressList
		groups.ForceSendFields = append(groups.ForceSendFields, "IncludeInGlobalAddressList")
		hasUpdates = true
	}

	// Archive settings
	if cmd.Flags().Changed("archive-only") {
		groups.ArchiveOnly = archiveOnly
		groups.ForceSendFields = append(groups.ForceSendFields, "ArchiveOnly")
		hasUpdates = true
	}
	if cmd.Flags().Changed("show-in-group-directory") {
		groups.ShowInGroupDirectory = showInGroupDirectory
		groups.ForceSendFields = append(groups.ForceSendFields, "ShowInGroupDirectory")
		hasUpdates = true
	}

	// Member management settings
	if cmd.Flags().Changed("who-can-leave-group") {
		groups.WhoCanLeaveGroup = whoCanLeaveGroup
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanLeaveGroup")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-add") {
		groups.WhoCanAdd = whoCanAdd
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanAdd")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-invite") {
		groups.WhoCanInvite = whoCanInvite
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanInvite")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-approve-members") {
		groups.WhoCanApproveMembers = whoCanApproveMembers
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanApproveMembers")
		hasUpdates = true
	}
	if cmd.Flags().Changed("allow-google-communication") {
		groups.AllowGoogleCommunication = allowGoogleCommunication
		groups.ForceSendFields = append(groups.ForceSendFields, "AllowGoogleCommunication")
		hasUpdates = true
	}
	if cmd.Flags().Changed("members-can-post-as-the-group") {
		groups.MembersCanPostAsTheGroup = membersCanPostAsTheGroup
		groups.ForceSendFields = append(groups.ForceSendFields, "MembersCanPostAsTheGroup")
		hasUpdates = true
	}

	// Moderation settings
	if cmd.Flags().Changed("who-can-contact-owner") {
		groups.WhoCanContactOwner = whoCanContactOwner
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanContactOwner")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-moderate-members") {
		groups.WhoCanModerateMembers = whoCanModerateMembers
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanModerateMembers")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-moderate-content") {
		groups.WhoCanModerateContent = whoCanModerateContent
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanModerateContent")
		hasUpdates = true
	}
	if cmd.Flags().Changed("who-can-ban-users") {
		groups.WhoCanBanUsers = whoCanBanUsers
		groups.ForceSendFields = append(groups.ForceSendFields, "WhoCanBanUsers")
		hasUpdates = true
	}

	if !hasUpdates {
		fmt.Fprintf(os.Stderr, "Error: No settings specified to update. Use --help to see available flags.\n")
		return fmt.Errorf("no settings specified")
	}

	// Update the group settings
	result, err := client.Groups.Update(groupEmail, groups).Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating group settings for %s: %v\n", groupEmail, err)
		return err
	}

	fmt.Printf("Successfully updated settings for group: %s\n", result.Email)

	// Display updated settings
	fmt.Println("\nUpdated settings:")
	for _, field := range groups.ForceSendFields {
		fmt.Printf("  %s\n", field)
	}

	return nil
}
