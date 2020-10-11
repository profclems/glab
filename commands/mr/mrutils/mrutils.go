package mrutils

import (
	"fmt"

	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/xanzy/go-gitlab"
)

func DisplayMR(mr *gitlab.MergeRequest) string {
	mrID := MRState(mr)
	return fmt.Sprintf("%s %s (%s)\n %s\n",
		mrID, mr.Title, mr.SourceBranch, mr.WebURL)
}

func MRState(m *gitlab.MergeRequest) string {
	if m.State == "opened" {
		return utils.Green(fmt.Sprintf("!%d", m.IID))
	} else if m.State == "merged" {
		return utils.Blue(fmt.Sprintf("!%d", m.IID))
	} else {
		return utils.Red(fmt.Sprintf("!%d", m.IID))
	}
}

func DisplayAllMRs(mrs []*gitlab.MergeRequest, projectID string) string {
	title := utils.NewListTitle("Merge Requests")
	title.RepoName = projectID
	title.CurrentPageTotal = len(mrs)

	table := tableprinter.NewTablePrinter()
	for _, m := range mrs {
		table.AddCell(MRState(m))
		table.AddCell(m.Title)
		table.AddCell(utils.Cyan(fmt.Sprintf("(%s) ← (%s)", m.TargetBranch, m.SourceBranch)))
		table.EndRow()
	}

	return fmt.Sprintf("%s\n%s", title.Describe(), table.Render())
}
