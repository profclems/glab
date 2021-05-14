package action

import (
	"fmt"
	"strconv"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionPipelines(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListProjectPipelinesOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.OrderBy = gitlab.String("updated_at")
		opts.PerPage = 100

		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if pipelines, err := api.ListProjectPipelines(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(pipelines)*2)
			for _, pipeline := range pipelines {
				description := pipeline.Ref
				if pipeline.CreatedAt != nil {
					description = fmt.Sprintf("%v (%v)", pipeline.Ref, utils.TimeToPrettyTimeAgo(*pipeline.CreatedAt))
				}
				vals = append(vals, strconv.Itoa(pipeline.ID), description)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
