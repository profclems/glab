package unsubscribe

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
)

func NewCmdUnsubscribe(f *cmdutils.Factory) *cobra.Command {
	var mrUnsubscribeCmd = &cobra.Command{
		Use:     "unsubscribe [<id> | <branch>]",
		Short:   `Unsubscribe from merge requests`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Example: heredoc.Doc(`
			$ glab mr unsubscribe 123
			$ glab mr unsub 123
			$ glab mr unsubscribe branch
			$ glab mr unsubscribe 123 branch  # unsubscribe from multiple MRs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args)
			if err != nil {
				return err
			}

			for _, mr := range mrs {
				if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
					Unsubscribed: true,
				}); err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "- Unsubscribing from Merge Request !%d\n", mr.IID)

				mr, err = api.UnsubscribeFromMR(apiClient, repo.FullName(), mr.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "%s You have successfully unsubscribed from merge request !%d\n", c.RedCheck(), mr.IID)
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr))
			}

			return nil
		},
	}

	return mrUnsubscribeCmd
}
