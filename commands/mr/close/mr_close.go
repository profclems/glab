package close

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdClose(f *cmdutils.Factory) *cobra.Command {
	var mrCloseCmd = &cobra.Command{
		Use:   "close [<id> | <branch>]",
		Short: `Close merge requests`,
		Long:  ``,
		Example: heredoc.Doc(`
			$ glab mr close 1
			$ glab mr close 1 2 3 4  # close multiple branches at once
			$ glab mr close  # use checked out branch
			$ glab mr close branch
			$ glab mr close username:branch
			$ glab mr close branch -R another/repo
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args, "opened")
			if err != nil {
				return err
			}

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("close")
			for _, mr := range mrs {
				if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
					Closed: true,
					Merged: true,
				}); err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "- Closing Merge request...\n")
				_, err := api.UpdateMR(apiClient, repo.FullName(), mr.IID, l)
				if err != nil {
					return err
				}
				// Update the value of the merge request to closed so that mrutils.DisplayMR
				// prints it as red
				mr.State = "closed"

				fmt.Fprintf(f.IO.StdOut, "%s Closed Merge request !%d\n", c.RedCheck(), mr.IID)
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr, f.IO.IsaTTY))
			}

			return nil
		},
	}

	return mrCloseCmd
}
