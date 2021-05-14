package action

import (
	"strings"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionGroups(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListGroupsOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if groups, err := api.ListGroups(client, opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(groups)*2)
			for _, group := range groups {
				vals = append(vals, group.Path, group.Description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}

func ActionSubgroups(cmd *cobra.Command, f *cmdutils.Factory, groupID string, opts *gitlab.ListSubgroupsOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if groups, err := api.ListSubgroups(client, groupID, opts); err != nil {
			if strings.Contains(err.Error(), "404 Group Not Found") {
				return carapace.ActionValues() // fail silently for repo completion
			}
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(groups)*2)
			for _, group := range groups {
				vals = append(vals, group.Path, group.Description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
