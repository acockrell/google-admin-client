package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	admin "google.golang.org/api/admin/directory/v1"

	"github.com/spf13/cobra"
)

// flags / parameters
var (
	fullOutput   bool
	csvOutput    bool
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// update-profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listUserCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "full export")
	listUserCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "csv export")
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

func listUserRunFunc(cmd *cobra.Command, args []string) {

	var email string

	if len(args) > 0 {
		email = args[0]
	}

	client, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// if email is supplied, display that user.  otherwise, display a list of
	// all users

	if email != "" {
		u, err := client.Users.Get(email).Do(Projection("FULL"))
		if err != nil {
			exitWithError(err.Error())
		}
		buf, _ := json.MarshalIndent(u, "", "  ")
		fmt.Printf("%s\n", buf)
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

		if fullOutput {
			buf, _ := json.MarshalIndent(u, "", "  ")
			fmt.Printf("%s\n", buf)
		} else if csvOutput {
			w := csv.NewWriter(os.Stdout)
			_ = w.Write([]string{"Name", "Email", "Admin", "OrgUnitPath"})
			for _, user := range u.Users {
				w.Write([]string{
					user.Name.FullName,
					user.PrimaryEmail,
					strconv.FormatBool(user.IsAdmin),
					user.OrgUnitPath,
				})
			}
			w.Flush()
		} else {
			for _, user := range u.Users {
				fmt.Printf("%s (%s) (Admin: %v) (OrgUnitPath: %s)\n", user.PrimaryEmail, user.Name.FullName, user.IsAdmin, user.OrgUnitPath)
			}
		}
	}
}
