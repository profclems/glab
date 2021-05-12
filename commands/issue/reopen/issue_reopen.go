package reopen

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/rsteube/carapace"

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
			c := f.IO.Color()

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

				fmt.Fprintf(out, "%s Reopened Issue #%d\n", c.GreenCheck(), issue.IID)
				fmt.Fprintln(out, issueutils.DisplayIssue(c, issue))
			}
			return nil
		},
	}

	carapace.Gen(issueReopenCmd).PositionalCompletion(
		action.ActionIssues(f, &gitlab.ListProjectIssuesOptions{State: gitlab.String("closed")}),
	)

	return issueReopenCmd
}
