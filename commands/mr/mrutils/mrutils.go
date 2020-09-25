package mrutils

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/profclems/glab/internal/utils"
	"github.com/xanzy/go-gitlab"
)

func DisplayMR(mr *gitlab.MergeRequest) string {
	var mrID string

	if mr.State == "opened" {
		mrID = utils.Green(fmt.Sprintf("#%d", mr.IID))
	} else {
		mrID = utils.Red(fmt.Sprintf("#%d", mr.IID))
	}

	return fmt.Sprintf("%s %s (%s)\n %s\n",
		mrID, mr.Title, mr.SourceBranch, mr.WebURL)
}

func DisplayAllMRs(m []*gitlab.MergeRequest, projectID string) *uitable.Table {
	return utils.DisplayList(utils.ListInfo{
		Name:    "Merge Requests",
		Columns: []string{"ID", "Title", "Branch"},
		Total:   len(m),
		GetCellValue: func(ri int, ci int) interface{} {
			mr := m[ri]
			switch ci {
			case 0:
				if mr.State == "opened" {
					return utils.Green(fmt.Sprintf("#%d", mr.IID))
				} else {
					return utils.Red(fmt.Sprintf("#%d", mr.IID))
				}
			case 1:
				return mr.Title
			case 2:
				return utils.Cyan(fmt.Sprintf("(%s) ← (%s)", mr.TargetBranch, mr.SourceBranch))
			default:
				return ""
			}
		},
	}, projectID)
}