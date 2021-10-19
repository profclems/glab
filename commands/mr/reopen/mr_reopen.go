package reopen

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var mrReopenCmd = &cobra.Command{
		Use:   "reopen [<id>... | <branch>...]",
		Short: `Reopen merge requests`,
		Example: heredoc.Doc(`
			$ glab mr reopen 123
			$ glab mr reopen 123 456 789
			$ glab mr reopen branch-1 branch-2
			$ glab mr reopen  # use checked out branch
		`),
		Aliases: []string{"open"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := f.IO.Color()
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args, "closed")
			if err != nil {
				return err
			}

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("reopen")
			for _, mr := range mrs {
				if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
					Opened: true,
					Merged: true,
				}); err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "- Reopening Merge request !%d...\n", mr.IID)
				mr, err = api.UpdateMR(apiClient, repo.FullName(), mr.IID, l)
				if err != nil {
					return err
				}

				fmt.Fprintf(f.IO.StdOut, "%s Reopened Merge request !%d\n", c.GreenCheck(), mr.IID)
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(f.IO.Color(), mr, f.IO.IsaTTY))
			}

			return nil
		},
	}

	return mrReopenCmd
}
