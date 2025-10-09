package cmd

import (
	"fmt"
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

// displayGroupInfo is no longer needed with unified formatter

func getGroupInfo(wg *sync.WaitGroup, client *admin.Service, group *admin.Group, results chan<- groupInfo) {
	defer wg.Done()

	var owners []string
	externalMembers := false
	formerEmployees := false

	r, err := client.Members.List(group.Id).Do()
	if err != nil {
		Logger.Error().Err(err).Str("group", group.Email).Msg("Failed to list group members")
		return
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
			u, err := client.Users.Get(m.Email).Do()
			if err != nil {
				Logger.Error().Err(err).Str("user", m.Email).Msg("Failed to get user details")
				continue
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
		InactiveMembers: inactiveOnly,
		ExternalMembers: externalMembers,
		FormerEmployees: formerEmployees,
	}

	// Skip if filtering for former employees only and this group doesn't have any
	if inactiveOnly && !formerEmployees {
		return
	}

	results <- gInfo
}

// groupMember represents a group member for list output
type groupMember struct {
	Email  string `json:"email"`
	Type   string `json:"type"`
	Status string `json:"status"`
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

	// if a group is supplied, display that group. otherwise, display a list of all groups
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

			// Build list of members with status
			var members []groupMember
			for _, i := range m.Members {
				status := "active"
				if i.Type == "GROUP" {
					status = "group"
				} else if i.Type == "USER" {
					u, err := client.Users.Get(i.Email).Do()
					if err != nil {
						Logger.Error().Err(err).Str("user", i.Email).Msg("Failed to get user details")
						continue
					}
					if _, ok := formerEmployeesOU[u.OrgUnitPath]; ok {
						status = "former"
					}
				}
				members = append(members, groupMember{
					Email:  i.Email,
					Type:   i.Type,
					Status: status,
				})
			}

			headers := []string{"Email", "Type", "Status"}
			if err := FormatOutput(members, headers); err != nil {
				exitWithError(fmt.Sprintf("Failed to format output: %s", err))
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

			if err := FormatOutput(g, nil); err != nil {
				exitWithError(fmt.Sprintf("Failed to format output: %s", err))
			}
		}
	} else {
		r, err := client.Groups.List().Customer("my_customer").Do()
		if err != nil {
			exitWithError(err.Error())
		}

		// Collect group info concurrently
		results := make(chan groupInfo, len(r.Groups))
		wg := new(sync.WaitGroup)

		for _, g := range r.Groups {
			wg.Add(1)
			go getGroupInfo(wg, client, g, results)
		}

		// Close results channel when all goroutines are done
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect all results
		var groupInfos []groupInfo
		for info := range results {
			groupInfos = append(groupInfos, info)
		}

		// Output using unified formatter
		headers := []string{"Name", "Description", "Email", "Owners", "Inactive Members", "External Members", "Former Employees"}
		if err := FormatOutput(groupInfos, headers); err != nil {
			exitWithError(fmt.Sprintf("Failed to format output: %s", err))
		}
	}
}
