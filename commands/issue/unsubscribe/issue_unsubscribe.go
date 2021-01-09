package unsubscribe

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdUnsubscribe(f *cmdutils.Factory) *cobra.Command {
	var issueUnsubscribeCmd = &cobra.Command{
		Use:     "unsubscribe <id>",
		Short:   `Unsubscribe to an issue`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Args:    cobra.ExactArgs(1),
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

			for _, issue := range issues {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintf(out, "- Unsubscribing from Issue #%d in %s\n", issue.IID, utils.Cyan(repo.FullName()))
				}

				issue, err := api.UnsubscribeFromIssue(apiClient, repo.FullName(), issue.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintln(out, utils.Red("✔"), "Unsubscribed")
				fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}

	return issueUnsubscribeCmd
}
