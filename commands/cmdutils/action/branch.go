package action

import (
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionBranches(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListBranchesOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if branches, err := api.ListBranches(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(branches)*2)
			for _, branch := range branches {
				vals = append(vals, branch.Name, branch.Commit.Title)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
