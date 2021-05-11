package action

import (
	"strconv"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/xanzy/go-gitlab"
)

func ActionIssues(f *cmdutils.Factory, opts *gitlab.ListProjectIssuesOptions) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		client, err := f.HttpClient()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
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
