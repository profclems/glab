package subscribe

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
)

func NewCmdSubscribe(f *cmdutils.Factory) *cobra.Command {
	var mrSubscribeCmd = &cobra.Command{
		Use:     "subscribe [<id> | <branch>]",
		Short:   `Subscribe to merge requests`,
		Long:    ``,
		Aliases: []string{"sub"},
		Example: heredoc.Doc(`
			$ glab mr subscribe 123
			$ glab mr sub 123
			$ glab mr subscribe branch
			$ glab mr subscribe 123 branch  # subscribe to multiple MRs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args, "any")
			if err != nil {
				return err
			}

			for _, mr := range mrs {
				if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
					Subscribed: true,
				}); err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "- Subscribing to merge request !%d\n", mr.IID)

				mr, err = api.SubscribeToMR(apiClient, repo.FullName(), mr.IID, nil)
				if err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "%s You have successfully subscribed to merge request !%d\n", c.GreenCheck(), mr.IID)
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr, f.IO.IsaTTY))
			}

			return nil
		},
	}

	return mrSubscribeCmd
}
