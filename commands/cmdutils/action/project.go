package action

import (
	"strings"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionGroupProjects(cmd *cobra.Command, f *cmdutils.Factory, groupID string, opts *gitlab.ListGroupProjectsOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.OrderBy = gitlab.String("updated_at")
		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if projects, err := api.ListGroupProjects(client, groupID, opts); err != nil {
			if strings.Contains(err.Error(), "404 Group Not Found") {
				return carapace.ActionValues() // fail silently for repo completion
			}
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(projects)*2)
			for _, project := range projects {
				vals = append(vals, project.Path, project.Description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}

func ActionUserProjects(cmd *cobra.Command, f *cmdutils.Factory, userID string, opts *gitlab.ListProjectsOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.OrderBy = gitlab.String("updated_at")
		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if projects, err := api.ListUserProjects(client, userID, opts); err != nil {
			if strings.Contains(err.Error(), "404 User Not Found") {
				return carapace.ActionValues() // fail silently for repo completion
			}
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(projects)*2)
			for _, project := range projects {
				vals = append(vals, project.Path, project.Description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
