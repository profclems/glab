package action

import (
	"time"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/cache"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionLabels(cmd *cobra.Command, f *cmdutils.Factory) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if labels, err := api.ListLabels(client, repo.FullName(), &gitlab.ListLabelsOptions{}); err != nil {
				return carapace.ActionMessage(err.Error())
			} else {
				vals := make([]string, 0, len(labels)*2)
				for _, label := range labels {
					vals = append(vals, label.Name, label.Description)
				}
				return carapace.ActionValuesDescribed(vals...)
			}
		}).Cache(1*time.Hour, cache.String(repo.FullName()))
	})
}
