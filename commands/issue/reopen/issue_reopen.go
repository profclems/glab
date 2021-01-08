package reopen

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var issueReopenCmd = &cobra.Command{
		Use:     "reopen <id>",
		Short:   `Reopen a closed issue`,
		Long:    ``,
		Aliases: []string{"open"},
		Example: heredoc.Doc(`
			$ glab issue reopen 123
			$ glab issue open 123
			$ glab issue reopen https://gitlab.com/profclems/glab/-/issues/123
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := f.IO.StdOut

			apiClient, err := f.HttpClient()

			if err != nil {
				return err
			}

			issues, repo, err := issueutils.IssuesFromArgs(apiClient, f.BaseRepo, args)
			if err != nil {
				return err
			}

			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("reopen")

			for _, issue := range issues {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintln(out, "- Reopening Issue...")
				}

				issue, err := api.UpdateIssue(apiClient, repo.FullName(), issue.IID, l)
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s Reopened Issue #%d\n", utils.GreenCheck(), issue.IID)
				fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}

	return issueReopenCmd
}
