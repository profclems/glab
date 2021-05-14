package action

import (
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionUsers(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListUsersOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if users, err := api.ListUsers(client, opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(users)*2)
			for _, user := range users {
				vals = append(vals, user.Username, user.Name)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}

func ActionProjectMembers(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListProjectMembersOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.Query = &c.CallbackValue
		opts.PerPage = 100

		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if members, err := api.ListProjectMembers(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(members)*2)
			for _, member := range members {
				vals = append(vals, member.Username, member.Name)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
