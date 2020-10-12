package issueutils

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/xanzy/go-gitlab"
)

func DisplayIssueList(issues []*gitlab.Issue, projectID string) string {
	table := tableprinter.NewTablePrinter()
	for _, issue := range issues {
		table.AddCell(IssueState(issue))
		table.AddCell(issue.Title)

		if len(issue.Labels) > 0 {
			table.AddCellf("(%s)", utils.Cyan(strings.Trim(strings.Join(issue.Labels, ", "), ",")))
		} else {
			table.AddCell("")
		}

		table.AddCell(utils.Gray(utils.TimeToPrettyTimeAgo(*issue.CreatedAt)))
		table.EndRow()
	}

	return table.Render()
}

func DisplayIssue(i *gitlab.Issue) string {
	duration := utils.TimeToPrettyTimeAgo(*i.CreatedAt)
	issueID := IssueState(i)

	return fmt.Sprintf("%s %s (%s)\n %s\n",
		issueID, i.Title, duration, i.WebURL)
}

func IssueState(i *gitlab.Issue) (issueID string) {
	if i.State == "opened" {
		issueID = utils.Green(fmt.Sprintf("#%d", i.IID))
	} else {
		issueID = utils.Red(fmt.Sprintf("#%d", i.IID))
	}
	return
}
