package action

import (
	"strconv"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func ActionIssues(cmd *cobra.Command, f *cmdutils.Factory, opts *gitlab.ListProjectIssuesOptions) carapace.Action {
	return ActionApiCallback(cmd, f, func(client *gitlab.Client, c carapace.Context) carapace.Action {
		opts.OrderBy = gitlab.String("updated_at")
		opts.PerPage = 100

		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if issues, err := api.ListIssues(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(issues)*2)
			for _, issue := range issues {
				vals = append(vals, strconv.Itoa(issue.IID), issue.Title)
			}
			return carapace.ActionValuesDescribed(vals...)
		}

	})
}
