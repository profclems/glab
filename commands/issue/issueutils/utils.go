package issueutils

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/utils"

	"github.com/gosuri/uitable"
	"github.com/xanzy/go-gitlab"
)

func DisplayAllIssues(m []*gitlab.Issue, projectID string) *uitable.Table {
	return utils.DisplayList(utils.ListInfo{
		Name:    "issues",
		Columns: []string{"IssueID", "Title", "Labels", "CreatedAt"},
		Total:   len(m),
		GetCellValue: func(ri int, ci int) interface{} {
			issue := m[ri]
			switch ci {
			case 0:
				if issue.State == "opened" {
					return utils.Green(fmt.Sprintf("#%d", issue.IID))
				} else {
					return utils.Red(fmt.Sprintf("#%d", issue.IID))
				}
			case 1:
				return issue.Title
			case 2:
				if len(issue.Labels) > 0 {
					return fmt.Sprintf("(%s)", utils.Cyan(strings.Trim(strings.Join(issue.Labels, ", "), ",")))
				}
				return ""
			case 3:
				return utils.Gray(utils.TimeToPrettyTimeAgo(*issue.CreatedAt))
			default:
				return ""
			}
		},
	}, projectID)
}

func DisplayIssue(i *gitlab.Issue) string {
	duration := utils.TimeToPrettyTimeAgo(*i.CreatedAt)
	var issueID string
	if i.State == "opened" {
		issueID = utils.Green(fmt.Sprintf("#%d", i.IID))
	} else {
		issueID = utils.Red(fmt.Sprintf("#%d", i.IID))
	}

	return fmt.Sprintf("%s %s (%s)\n %s\n",
		issueID, i.Title, duration, i.WebURL)
}
