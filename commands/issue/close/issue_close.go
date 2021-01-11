package close

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdClose(f *cmdutils.Factory) *cobra.Command {
	var issueCloseCmd = &cobra.Command{
		Use:   "close <id>",
		Short: `Close an issue`,
		Long:  ``,
		Example: heredoc.Doc(`
			$ glab issue close 123
			$ glab issue close https://gitlab.com/profclems/glab/-/issues/123
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			issues, repo, err := issueutils.IssuesFromArgs(apiClient, f.BaseRepo, args)
			if err != nil {
				return err
			}

			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("close")

			for _, issue := range issues {
				fmt.Fprintln(f.IO.StdOut, "- Closing Issue...")
				issue, err := api.UpdateIssue(apiClient, repo.FullName(), issue.IID, l)
				if err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "%s Closed Issue #%d\n", utils.RedCheck(), issue.IID)
				fmt.Fprintln(f.IO.StdOut, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}
	return issueCloseCmd
}
