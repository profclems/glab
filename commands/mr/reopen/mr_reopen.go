package reopen

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var mrReopenCmd = &cobra.Command{
		Use:     "reopen [<id> | <branch>]",
		Short:   `Reopen merge requests`,
		Long:    ``,
		Args:    cobra.MaximumNArgs(1),
		Aliases: []string{"open"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c := f.IO.Color()
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args)
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
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(f.IO.Color(), mr))
			}

			return nil
		},
	}

	return mrReopenCmd
}
