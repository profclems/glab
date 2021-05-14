package action

import (
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionTags(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListTagsOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		opts.Search = &c.CallbackValue
		opts.PerPage = 100

		if tags, err := api.ListTags(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(tags)*2)
			for _, tag := range tags {
				vals = append(vals, tag.Name, tag.Commit.Title)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
