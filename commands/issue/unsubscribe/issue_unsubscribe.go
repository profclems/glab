package unsubscribe

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/rsteube/carapace"
	"github.com/xanzy/go-gitlab"

	"github.com/spf13/cobra"
)

func NewCmdUnsubscribe(f *cmdutils.Factory) *cobra.Command {
	var issueUnsubscribeCmd = &cobra.Command{
		Use:     "unsubscribe <id>",
		Short:   `Unsubscribe to an issue`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Example: heredoc.Doc(`
			$ glab issue unsubscribe 123
			$ glab issue unsub 123
			$ glab issue unsubscribe https://gitlab.com/profclems/glab/-/issues/123
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := f.IO.Color()
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			issues, repo, err := issueutils.IssuesFromArgs(apiClient, f.BaseRepo, args)
			if err != nil {
				return err
			}

			for _, issue := range issues {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintf(f.IO.StdErr, "- Unsubscribing from Issue #%d in %s\n", issue.IID, c.Cyan(repo.FullName()))
				}

				issue, err := api.UnsubscribeFromIssue(apiClient, repo.FullName(), issue.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintln(f.IO.StdErr, c.RedCheck(), "Unsubscribed")
				fmt.Fprintln(f.IO.StdOut, issueutils.DisplayIssue(c, issue))
			}
			return nil
		},
	}

	carapace.Gen(issueUnsubscribeCmd).PositionalCompletion(
		action.ActionIssues(issueUnsubscribeCmd, f, &gitlab.ListProjectIssuesOptions{State: gitlab.String("opened")}),
	)

	return issueUnsubscribeCmd
}
