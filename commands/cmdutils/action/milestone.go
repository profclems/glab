package action

import (
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionMilestones(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListMilestonesOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if milestones, err := api.ListMilestones(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(milestones)*2)
			for _, milestone := range milestones {
				vals = append(vals, milestone.Title, milestone.Description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
