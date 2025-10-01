package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	admin "google.golang.org/api/admin/directory/v1"
)

// staffOU lists Google Directory OUs which contain FTE vs
// contractors. This list is used to determine if a given
// G Suite user is an "internal" or "external" member when
// doing group audits.
//
// https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
var staffOU = map[string]struct{}{
	"/":                   {},
	"/Customer Education": {},
	"/Customer Success":   {},
	"/Engineering":        {},
	"/Marketing":          {},
	"/Product":            {},
	"/Sales":              {},
	"/Sales Engineering":  {},
}

var formerEmployeesOU = map[string]struct{}{
	"/Former employees": {},
}

type groupInfo struct {
	Name            string
	Description     string
	Email           string
	Owners          string
	InactiveMembers bool
	ExternalMembers bool
	FormerEmployees bool
}

var (
	getMembers   = false
	inactiveOnly = false
)

// listGroupCmd represents the group list command
var listGroupCmd = &cobra.Command{
	Use:   "list",
	Short: "list groups",
	Run:   listGroupRunFunc,
	Long: `
List user(s).

Usage
-----

$ gac group list
$ gac group list --contains-former-employees
$ gac group list operations-group@example.com
$ gac group list operations-group@example.com --get-members
`,
}

func init() {
	groupCmd.AddCommand(listGroupCmd)

	listGroupCmd.Flags().BoolVarP(&getMembers, "get-members", "m", getMembers, "lists the group members")
	listGroupCmd.Flags().BoolVarP(&inactiveOnly, "contains-former-employees", "i", inactiveOnly, "shows only groups with inactive members")
}

func displayGroupInfo(writer *csv.Writer, info groupInfo) {
	err := writer.Write([]string{
		info.Name,
		info.Description,
		info.Email,
		info.Owners,
		strconv.FormatBool(info.InactiveMembers),
		strconv.FormatBool(info.ExternalMembers),
		strconv.FormatBool(info.FormerEmployees),
	})
	if err != nil {
		exitWithError(err.Error())
	}
}

func getGroupInfo(wg *sync.WaitGroup, client *admin.Service, group *admin.Group, writer *csv.Writer) {
	var owners []string
	externalMembers := false
	inactiveMembers := false
	formerEmployees := false
	//
	r, err := client.Members.List(group.Id).Do()
	if err != nil {
		exitWithError(err.Error())
	}

	for _, m := range r.Members {
		configuredDomain := getDomain()
		if configuredDomain != "" && !strings.HasSuffix(m.Email, "@"+configuredDomain) {
			externalMembers = true
			continue
		}

		if m.Role != "MEMBER" {
			owners = append(owners, m.Email)
		}

		if m.Type == "USER" {
			if m.Status != "ACTIVE" {
				inactiveMembers = true
			}

			u, err := client.Users.Get(m.Email).Do()
			if err != nil {
				exitWithError(err.Error())
			}
			if _, ok := staffOU[u.OrgUnitPath]; !ok {
				externalMembers = true
			}

			if _, ok := formerEmployeesOU[u.OrgUnitPath]; ok {
				formerEmployees = true
			}
		}
	}

	gInfo := groupInfo{
		Name:            group.Name,
		Description:     group.Description,
		Email:           group.Email,
		Owners:          strings.Join(owners, ","),
		InactiveMembers: inactiveMembers,
		ExternalMembers: externalMembers,
		FormerEmployees: formerEmployees,
	}

	if inactiveOnly {
		if !formerEmployees {
			writer.Flush()
			wg.Done()
			return
		}
	}
	displayGroupInfo(writer, gInfo)

	writer.Flush()
	wg.Done()
}

func listGroupRunFunc(cmd *cobra.Command, args []string) {

	var group string

	if len(args) > 0 {
		group = args[0]
	}

	client, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// if a group is supplied, display that group.  otherwise, display a list of
	// all groups
	if group != "" {
		if getMembers {
			groupEmail := group
			if !strings.Contains(group, "@") {
				groupEmail = group + "@" + getDomain()
			}
			m, err := client.Members.List(groupEmail).Do()
			if err != nil {
				exitWithError(err.Error())
			}
			for _, i := range m.Members {
				mark := "\u2713"

				if i.Type == "GROUP" {
					fmt.Println(i.Email, "\u271B")
				}

				if i.Type == "USER" {
					u, err := client.Users.Get(i.Email).Do()

					if err != nil {
						exitWithError(err.Error())
					}
					if _, ok := formerEmployeesOU[u.OrgUnitPath]; ok {
						mark = "\u0078"
					}
					fmt.Println(i.Email, mark)
				}
			}
		} else {
			groupEmail := group
			if !strings.Contains(group, "@") {
				groupEmail = group + "@" + getDomain()
			}
			g, err := client.Groups.Get(groupEmail).Do()
			if err != nil {
				exitWithError(err.Error())
			}
			buf, _ := json.MarshalIndent(g, "", "  ")
			fmt.Printf("%s\n", buf)
		}
	} else {
		r, err := client.Groups.List().Customer("my_customer").Do()
		if err != nil {
			exitWithError(err.Error())
		}
		// header line
		w := csv.NewWriter(os.Stdout)
		err = w.Write([]string{"Name", "Description", "Email", "Owners", "Inactive Members", "External Members", "Former Employees"})
		if err != nil {
			exitWithError(err.Error())
		}

		wg := new(sync.WaitGroup)
		wg.Add(10)

		for idx, g := range r.Groups {
			if idx > 0 {
				if idx%10 == 0 {
					wg.Wait()
					wg.Add(10)
				}
			}
			go getGroupInfo(wg, client, g, w)
		}

		wg.Wait()
	}

}
