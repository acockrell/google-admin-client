package cmd

import (
	"fmt"

	admin "google.golang.org/api/admin/directory/v1"

	"github.com/spf13/cobra"
)

// flags / parameters
var (
	fullOutput   bool // Deprecated: use --format=json instead
	csvOutput    bool // Deprecated: use --format=csv instead
	disabledOnly = false
)

// listUserCmd represents the update-profile command
var listUserCmd = &cobra.Command{
	Use:   "list",
	Short: "list users",
	Run:   listUserRunFunc,
	Long: `
List user(s).

Usage
-----

$ gac user list
$ gac user list --disabled-only
$ gac user list username@example.com
`,
}

func init() {
	userCmd.AddCommand(listUserCmd)

	// Backward compatibility flags (deprecated)
	listUserCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "deprecated: use --format=json instead")
	listUserCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "deprecated: use --format=csv instead")
	if err := listUserCmd.Flags().MarkDeprecated("full", "use --format=json instead"); err != nil {
		Logger.Warn().Err(err).Msg("Failed to mark full flag as deprecated")
	}
	if err := listUserCmd.Flags().MarkDeprecated("csv", "use --format=csv instead"); err != nil {
		Logger.Warn().Err(err).Msg("Failed to mark csv flag as deprecated")
	}

	// Other flags
	listUserCmd.Flags().BoolVarP(&disabledOnly, "disabled-only", "d", disabledOnly, "lists only disabled accounts")
}

/*
NOTE: This is in place of using .Query("orgUnitPath=/path/foo") because the
existing of a space in the path name Regardless of encoding the path it
results in HTTP 400
*/
func filterDisabledOnly(u admin.Users) admin.Users {
	var formerEmployees admin.Users
	for _, user := range u.Users {
		if user.OrgUnitPath == "/Former employees" {
			formerEmployees.Users = append(formerEmployees.Users, user)
		}
	}
	return formerEmployees
}

// userListItem represents a simplified user for list output
type userListItem struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Admin       string `json:"admin"`
	OrgUnitPath string `json:"orgUnitPath"`
}

func listUserRunFunc(cmd *cobra.Command, args []string) {
	var email string

	if len(args) > 0 {
		email = args[0]
	}

	client, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// Handle backward compatibility with deprecated flags
	originalFormat := outputFormat
	if fullOutput {
		outputFormat = OutputFormatJSON
	} else if csvOutput {
		outputFormat = OutputFormatCSV
	}

	// if email is supplied, display that user. otherwise, display a list of all users
	if email != "" {
		u, err := client.Users.Get(email).Do(Projection("FULL"))
		if err != nil {
			exitWithError(err.Error())
		}

		// For single user, always show full details
		if err := FormatOutput(u, nil); err != nil {
			exitWithError(fmt.Sprintf("Failed to format output: %s", err))
		}
	} else {
		var u admin.Users
		var pageToken string
		for {
			res, err := client.Users.List().Customer("my_customer").PageToken(pageToken).Do()
			if err != nil {
				exitWithError(err.Error())
			}
			u.Users = append(u.Users, res.Users...)

			if res.NextPageToken == "" {
				break
			}
			pageToken = res.NextPageToken
		}

		if disabledOnly {
			u = filterDisabledOnly(u)
		}

		// Convert to simplified list items for CSV/table/plain formats
		headers := []string{"Name", "Email", "Admin", "OrgUnitPath"}
		var items []userListItem
		for _, user := range u.Users {
			item := userListItem{
				Name:        user.Name.FullName,
				Email:       user.PrimaryEmail,
				Admin:       fmt.Sprintf("%v", user.IsAdmin),
				OrgUnitPath: user.OrgUnitPath,
			}
			items = append(items, item)
		}

		// For JSON/YAML, output full user data
		var outputData interface{}
		if outputFormat == OutputFormatJSON || outputFormat == OutputFormatYAML {
			outputData = u.Users
		} else {
			outputData = items
		}

		if err := FormatOutput(outputData, headers); err != nil {
			exitWithError(fmt.Sprintf("Failed to format output: %s", err))
		}
	}

	// Restore original format
	outputFormat = originalFormat
}
