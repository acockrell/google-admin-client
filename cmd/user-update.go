package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"

	"github.com/spf13/cobra"
)

// flags / parameters
var (
	address      string
	dept         string
	employeeID   string
	employeeType string
	forceUpdate  = false
	managerEmail string
	ou           string
	phone        string
	removeUser   = false
	title        string
	clearPII     = false
	// defined in group-list
	// groups  []string
)

// updateUserCmd represents the update-profile command
var updateUserCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the specified user",
	Run:   updateUserRunFunc,
	Long: `
Update the speficied user.

Usage
-----

	$ gac user update --address "Columbus, OH" jdoe@example.com
	$ gac user update --dept Engineering jdoe@example.com
	$ gac user update --group dev jdoe@example.com
	$ gac user update -g dev -g info jdoe@example.com
	$ gac user update --id $(uuidgen) jdoe@example.com
	$ gac user update --id $(uuidgen) --force jdoe@example.com
	$ gac user update --manager manager@example.com jdoe@example.com
	$ gac user update --ou /some/path jdoe@example.com
	$ gac user update --phone mobile:703-555-5555 jdoe@example.com
	$ gac user update --phone 'mobile:703-555-5555; work:301-684-8080,555' jdoe@example.com
	$ gac user update --remove jdoe@example.com
	$ gac user update --title "Sales Engineer" jdoe@example.com
	$ gac user update --clear-pii jdoe@example.com

`,
}

type orgArgs struct {
	Dept  string
	Title string
}

// Projection and associated type/method implements CallOption
// interface to enable seeing custom user attributes.
//
// https://developers.google.com/admin-sdk/directory/reference/rest/v1/users/get#Projection
// https://github.com/googleapis/google-api-go-client/blob/v0.14.0/googleapi/googleapi.go#L378
func Projection(p string) googleapi.CallOption {
	return projection(p)
}

type projection string

func (p projection) Get() (string, string) {
	return "projection", string(p)
}

func init() {
	userCmd.AddCommand(updateUserCmd)
	// Here you will define your flags and configuration settings.

	// ensure groups is empty for update operations
	groups = []string{}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// update-profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	updateUserCmd.Flags().StringVarP(&address, "address", "a", address, "address")
	updateUserCmd.Flags().StringVarP(&dept, "dept", "d", "", "department")
	updateUserCmd.Flags().StringVarP(&employeeID, "id", "i", "", "employee UUID")
	updateUserCmd.Flags().StringVarP(&employeeType, "type", "e", "", "staff or contractor")
	updateUserCmd.Flags().BoolVarP(&forceUpdate, "force", "f", forceUpdate, "overwrite existing values (e.g. employee ID)")
	updateUserCmd.Flags().StringSliceVarP(&groups, "group", "g", groups, "groups")
	updateUserCmd.Flags().StringVarP(&ou, "ou", "o", "", "org unit path")
	updateUserCmd.Flags().StringVarP(&managerEmail, "manager", "m", "", "manager's email")
	updateUserCmd.Flags().StringVarP(&phone, "phone", "p", "", "phone")
	updateUserCmd.Flags().BoolVarP(&removeUser, "remove", "r", removeUser, "disable user account")
	updateUserCmd.Flags().StringVarP(&title, "title", "t", "", "title")
	updateUserCmd.Flags().BoolVarP(&clearPII, "clear-pii", "", clearPII, "clear personal information")
}

func updateUserRunFunc(cmd *cobra.Command, args []string) {
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

	// if parameters aren't supplied, create a user based on stdin
	user := new(admin.User)
	if address == "" && dept == "" && employeeID == "" && employeeType == "" && ou == "" && managerEmail == "" && phone == "" && title == "" && len(groups) == 0 {
		j, _ := io.ReadAll(os.Stdin)
		err := json.Unmarshal(j, &user)
		if err != nil {
			exitWithError(err.Error())
		}
	} else {
		if removeUser {
			// signout sessions on all devices, reset cookies
			client.Users.SignOut(email)

			// asp/token edits == scope error
			// Error 403: Request had insufficient authentication scopes.
			//
			// remove any application specific passwords
			// as, err := client.Asps.List(email).Do()
			// if err != nil {
			// 	exitWithError(err.Error())
			// }
			// for _, a := range as.Items {
			// 	client.Asps.Delete(email, a.CodeId)
			// }
			//
			// remove any tokens issues to 3rd party apps
			// ts, err := client.Tokens.List(email).Do()
			// if err != nil {
			// 	exitWithError(err.Error())
			// }
			// for _, t := range ts.Items {
			// 	client.Tokens.Delete(email, t.ClientId)
			// }

			clearUserPII(user)
			disableGsuiteUser(user)

			gs, _ := client.Groups.List().UserKey(email).Do()
			for _, g := range gs.Groups {
				err := client.Members.Delete(g.Email, email).Do()
				if err != nil {
					exitWithError(fmt.Sprintf("Unable to remove %s from group %s: %s", email, g.Email, err))
				}
			}
		} else if clearPII {
			// If you just want to Clear PII without disabling the user.  Useful for testing.
			clearUserPII(user)
		} else {
			if address != "" {
				user.Addresses = parseAddress(SanitizeInput(address))
			}
			if dept != "" || title != "" {
				// Validate department if provided
				if dept != "" {
					if err := ValidateDepartment(dept); err != nil {
						exitWithError(fmt.Sprintf("invalid department: %s", err))
					}
				}
				user.Organizations = parseOrg(&orgArgs{Dept: SanitizeInput(dept), Title: SanitizeInput(title)})
			}
			if employeeID != "" {
				// Validate employee ID as UUID
				if err := ValidateUUID(employeeID); err != nil {
					exitWithError(fmt.Sprintf("invalid employee ID: %s", err))
				}

				u, err := client.Users.Get(email).Do()
				if err != nil {
					exitWithError(err.Error())
				}

				if u.ExternalIds == nil || forceUpdate {
					user.ExternalIds = parseID(employeeID)
				} else {
					fmt.Println("Skipping update of existing Employee ID, use --force.")
				}
			}
			if employeeType != "" {
				user.CustomSchemas = parseType(employeeType)
				// lists user with custom attributes
				// e.g. to reverse engineer wtf custom attributes look like :-/
				// u, _ := client.Users.Get(email).Do(Projection("FULL"))
				// fmt.Printf("DEBUG: %+v", string(u.CustomSchemas["Employee_Type"]))
			}
			if ou != "" {
				user.OrgUnitPath = ou
			}
			if managerEmail != "" {
				managerEmail = SanitizeInput(managerEmail)
				// Validate manager email
				if err := ValidateEmail(managerEmail); err != nil {
					exitWithError(fmt.Sprintf("invalid manager email: %s", err))
				}
				user.Relations = parseManager(managerEmail)
			}
			if phone != "" {
				// Validate phone numbers (handles multiple phones separated by semicolon)
				phones := strings.Split(phone, ";")
				for _, p := range phones {
					if err := ValidatePhoneNumber(strings.TrimSpace(p)); err != nil {
						exitWithError(fmt.Sprintf("invalid phone number: %s", err))
					}
				}
				user.Phones = parsePhone(phone)
			}
		}
	}

	_, err = client.Users.Update(email, user).Do()
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
		_, err = client.Members.Insert(groupEmail, &admin.Member{Email: email}).Do()
		if err != nil {
			exitWithError(fmt.Sprintf("Unable to add %s to group %s: %s", email, g, err))
		}
	}
}

// ...
func disableGsuiteUser(u *admin.User) {
	u.ChangePasswordAtNextLogin = false
	u.IncludeInGlobalAddressList = false
	u.OrgUnitPath = "/Former employees"
	u.Password = randomPassword(12)

	// this stops mail delivery...
	//u.Suspended = true
}

func clearUserPII(u *admin.User) {
	u.ForceSendFields = []string{
		"RecoveryEmail",
		"RecoveryPhone",
	}

	// These are JSON arrays.  Sending Null fields will flush them out.
	// No need to set them as blank strings.
	u.NullFields = []string{
		"Addresses",
		"Emails",
	}

	u.RecoveryEmail = ""
	u.RecoveryPhone = ""
}

// parse a phone string like "mobile:<number>" or "mobile:<number>;work:<number>"
func parsePhone(str string) (phones []admin.UserPhone) {
	p := strings.Split(str, ";")
	for _, s := range p {
		split := strings.SplitN(strings.Trim(s, " "), ":", 2)
		if len(split) > 0 {
			phones = append(phones, admin.UserPhone{
				Type:  split[0],
				Value: split[1],
			})
		}
	}
	return phones
}

// parse manager email
func parseManager(managerEmail string) (relations []admin.UserRelation) {
	r := admin.UserRelation{
		Type:  "manager",
		Value: managerEmail,
	}
	return []admin.UserRelation{r}
}

// parse an address string
func parseAddress(address string) []admin.UserAddress {
	a := admin.UserAddress{Formatted: address}
	return []admin.UserAddress{a}
}

// parse user organization
func parseOrg(args *orgArgs) []admin.UserOrganization {
	o := admin.UserOrganization{Primary: true}
	if args.Dept != "" {
		o.Department = args.Dept
	}
	if args.Title != "" {
		o.Title = args.Title
	}
	return []admin.UserOrganization{o}
}

// parse employee type custom attribute
func parseType(employeeType string) map[string]googleapi.RawMessage {
	schema := make(map[string]googleapi.RawMessage)

	// U-G-L-Y but works for now...
	var t = []byte(`{"Staff":[{"type":"work","value":"Yes"}],"Contractor":[{"type":"work","value":""}]}`)
	if employeeType == "contractor" {
		t = []byte(`{"Staff":[{"type":"work","value":""}],"Contractor":[{"type":"work","value":"Yes"}]}`)
	}
	schema["Employee_Type"] = t

	return schema
}

// parse employee ID
func parseID(employeeID string) []admin.UserExternalId {
	i := admin.UserExternalId{
		Type:  "organization",
		Value: employeeID,
	}
	return []admin.UserExternalId{i}
}
