package cmd

import (
	"fmt"
	"strings"

	admin "google.golang.org/api/admin/directory/v1"
)

// createUserFlags holds the flags for user creation
type createUserFlags struct {
	groups        []string
	personalEmail string
	firstName     string
	lastName      string
}

// createUserWithClient is the testable version of createUserRunFunc
// It accepts a client interface and returns errors instead of calling os.Exit
func createUserWithClient(client adminClientInterface, args []string, flags createUserFlags) error {
	if len(args) == 0 {
		return fmt.Errorf("email is a required argument")
	}

	email := SanitizeInput(args[0])

	// Validate email address
	if err := ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	user := admin.User{}
	user.PrimaryEmail = email
	user.ChangePasswordAtNextLogin = true
	user.Password = randomPassword(12)

	// Skip interactive collection if flags are provided
	if flags.personalEmail == "" || flags.firstName == "" || flags.lastName == "" {
		return fmt.Errorf("all user details must be provided via flags (use -e, -f, -l)")
	}

	updateUser(&user, flags.personalEmail, flags.firstName, flags.lastName)

	_, err := client.InsertUser(&user)
	if err != nil {
		return fmt.Errorf("unable to create user %s: %w", email, err)
	}

	// Add user to groups
	for _, g := range flags.groups {
		// Validate group name
		if err := ValidateGroupName(g); err != nil {
			return fmt.Errorf("invalid group name '%s': %w", g, err)
		}

		groupEmail := g
		if !strings.Contains(g, "@") {
			groupEmail = g + "@" + getDomain()
		}

		_, err = client.InsertMember(groupEmail, &admin.Member{Email: user.PrimaryEmail})
		if err != nil {
			return fmt.Errorf("unable to add %s to group %s: %w", user.PrimaryEmail, g, err)
		}
	}

	return nil
}
