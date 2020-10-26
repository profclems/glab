package issues

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"

	"github.com/xanzy/go-gitlab"
)

func NewCmdIssues(f *cmdutils.Factory) *cobra.Command {
	var mrIssuesCmd = &cobra.Command{
		Use:     "issues [<id> | <branch>]",
		Short:   `Get issues related to a particular merge request.`,
		Long:    ``,
		Aliases: []string{"issue"},
		Args:    cobra.MaximumNArgs(1),
		Example: "$ glab mr issues 46",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			l := &gitlab.GetIssuesClosedOnMergeOptions{}

			mrIssues, err := api.GetMRLinkedIssues(apiClient, repo.FullName(), mr.IID, l)
			if err != nil {
				return err
			}

			title := utils.NewListTitle("issue")
			title.RepoName = repo.FullName()
			title.Page = 0
			title.ListActionType = "search"
			title.CurrentPageTotal = len(mrIssues)

			fmt.Fprintf(out, "%s\n%s\n", title.Describe(), issueutils.DisplayIssueList(mrIssues, repo.FullName()))
			return nil
		},
	}

	return mrIssuesCmd
}
