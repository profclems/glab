package subscribe

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
)

func NewCmdSubscribe(f *cmdutils.Factory) *cobra.Command {
	var issueSubscribeCmd = &cobra.Command{
		Use:     "subscribe <id>",
		Short:   `Subscribe to an issue`,
		Long:    ``,
		Aliases: []string{"sub"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := f.IO.StdOut
			var err error

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
					fmt.Fprintf(out, "- Subscribing to Issue #%d in %s\n", issue.IID, utils.Cyan(repo.FullName()))
				}

				issue, err := api.SubscribeToIssue(apiClient, repo.FullName(), issue.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintln(out, utils.GreenCheck(), "Subscribed")
				fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}

	return issueSubscribeCmd
}
