package action

import (
	"strconv"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/rsteube/carapace"
	"github.com/xanzy/go-gitlab"
)

func ActionMergeRequests(f *cmdutils.Factory, opts *gitlab.ListProjectMergeRequestsOptions) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		client, err := f.HttpClient()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		repo, err := f.BaseRepo()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		if mergeRequests, err := api.ListMRs(client, repo.FullName(), opts); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			vals := make([]string, 0, len(mergeRequests)*2)
			for _, mr := range mergeRequests {
				vals = append(vals, strconv.Itoa(mr.IID), mr.Title)
			}
			return carapace.ActionValuesDescribed(vals...)
		}
	})
}
