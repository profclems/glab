package subscribe

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
)

func NewCmdSubscribe(f *cmdutils.Factory) *cobra.Command {
	var issueSubscribeCmd = &cobra.Command{
		Use:     "subscribe <id>",
		Short:   `Subscribe to an issue`,
		Long:    ``,
		Aliases: []string{"sub"},
		Example: heredoc.Doc(`
			$ glab issue subscribe 123
			$ glab issue sub 123
			$ glab issue subscribe https://gitlab.com/profclems/glab/-/issues/123
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
				if f.IO.IsErrTTY && f.IO.IsaTTY {
					fmt.Fprintf(f.IO.StdErr, "- Subscribing to Issue #%d in %s\n", issue.IID, c.Cyan(repo.FullName()))
				}

				issue, err := api.SubscribeToIssue(apiClient, repo.FullName(), issue.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintln(f.IO.StdErr, c.GreenCheck(), "Subscribed")
				fmt.Fprintln(f.IO.StdOut, issueutils.DisplayIssue(c, issue))
			}
			return nil
		},
	}

	return issueSubscribeCmd
}
