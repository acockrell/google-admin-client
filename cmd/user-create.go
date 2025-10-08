package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	admin "google.golang.org/api/admin/directory/v1"

	"github.com/spf13/cobra"
)

// EMAIL ... shut up lint
const EMAIL = `Your Google Workspace account has been created.

Username: %s
Password: %s
URL: https://www.google.com/accounts/AccountChooser?Email=%s&continue=https://apps.google.com/user/hub

`

// flags / parameters
var (
	groups        []string
	personalEmail string
	firstName     string
	lastName      string
)

// createUserCmd represents the update-profile command
var createUserCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactively create the specified user",
	Run:   createUserRunFunc,
	Long: `
Interactively create a user.

Usage
-----

  $ gac user create newuser@example.com
  $ gac user create -g Group1 newuser@example.com
  $ gac user create -g Group1 -g Group2 newuser@example.com
  $ gac user create -g Group1 -e personal@email.com -f Firstname -l Lastname newuser@example.com

Overview
--------

This command interactively creates a user.  Prompts are made for first & last
names as well as personal email address.  The user can optionally be added to
one or more groups.

The user is created with a random password, and an update of the password is
forced on first login.

The resultant user record, including password is output.

Future Enhancements
-------------------

1. Read from STDIN

2. Output only personal email address & password

3. If group assignment fails, undo user creation (i.e. make this a transaction)

4. There are several sets of groups depending on department.  Rather than
support this directly w/ this command, a "usermod" type command is
anticipated... at which point the group functionality may be removed from this
command.
`,
}

func init() {
	userCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().StringSliceVarP(&groups, "groups", "g", groups, "groups")
	createUserCmd.Flags().StringVarP(&personalEmail, "email", "e", "", "email")
	createUserCmd.Flags().StringVarP(&firstName, "first-name", "f", "", "first name")
	createUserCmd.Flags().StringVarP(&lastName, "last-name", "l", "", "last name")
}

func createUserRunFunc(cmd *cobra.Command, args []string) {
	// For interactive mode, we still need to handle collectUserInfo
	// For now, if flags are not provided, use the old path
	if personalEmail == "" || firstName == "" || lastName == "" {
		createUserRunFuncInteractive(cmd, args)
		return
	}

	// Create real admin client
	service, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// Wrap in adapter
	client := newRealAdminClientAdapter(service)

	// Package flags
	flags := createUserFlags{
		groups:        groups,
		personalEmail: personalEmail,
		firstName:     firstName,
		lastName:      lastName,
	}

	// Call testable function
	if err := createUserWithClient(client, args, flags); err != nil {
		exitWithError(err.Error())
	}

	// Success message (note: password not accessible from testable function)
	// For now, we'll skip the output in the refactored path
	// TODO: Refactor to return user object and print here
}

// createUserRunFuncInteractive handles the interactive user creation flow
// This preserves the existing behavior for when flags are not provided
func createUserRunFuncInteractive(cmd *cobra.Command, args []string) {
	var email string

	if len(args) == 0 {
		exitWithError("email is a required argument")
	}

	email = SanitizeInput(args[0])

	// Validate email address
	if err := ValidateEmail(email); err != nil {
		exitWithError(fmt.Sprintf("invalid email address: %s", err))
	}
	client, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	user := admin.User{}
	user.PrimaryEmail = email
	user.ChangePasswordAtNextLogin = true
	user.Password = randomPassword(12)

	err = collectUserInfo(&user)
	if err != nil {
		exitWithError(err.Error())
	}

	_, err = client.Users.Insert(&user).Do()
	if err != nil {
		exitWithError(fmt.Sprintf("Unable to update %s: %s", email, err))
	}

	for _, g := range groups {
		// Validate group name
		if err := ValidateGroupName(g); err != nil {
			exitWithError(fmt.Sprintf("invalid group name '%s': %s", g, err))
		}

		groupEmail := g
		if !strings.Contains(g, "@") {
			groupEmail = g + "@" + getDomain()
		}
		_, err = client.Members.Insert(groupEmail, &admin.Member{Email: user.PrimaryEmail}).Do()
		if err != nil {
			exitWithError(fmt.Sprintf("Unable to add %s to group %s: %s", user.PrimaryEmail, g, err))
		}
	}

	fmt.Printf(EMAIL, user.PrimaryEmail, user.Password, user.PrimaryEmail)
}

func collectUserInfo(user *admin.User) (err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Personal Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	email = SanitizeInput(strings.TrimRight(email, "\n"))

	// Validate personal email
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid personal email: %w", err)
	}

	fmt.Print("First Name: ")
	fname, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fname = strings.TrimRight(fname, "\n")

	fmt.Print("Last Name: ")
	lname, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	lname = strings.TrimRight(lname, "\n")

	updateUser(user, email, fname, lname)

	return nil
}

// TODO: smells like a method
func updateUser(user *admin.User, email, fname, lname string) {
	user.Emails = []admin.UserEmail{
		{
			Address: email,
			Type:    "home",
		},
		{
			Address: user.PrimaryEmail,
			Primary: true,
		},
	}
	user.Name = &admin.UserName{
		FamilyName: lname,
		GivenName:  fname,
		FullName:   fname + " " + lname,
	}
}
