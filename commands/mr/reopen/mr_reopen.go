package reopen

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var mrReopenCmd = &cobra.Command{
		Use:     "reopen [<id> | <branch>]",
		Short:   `Reopen merge requests`,
		Long:    ``,
		Aliases: []string{"open"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
				Opened: true,
				Merged: true,
			}); err != nil {
				return err
			}

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("reopen")

			fmt.Fprintf(f.IO.StdOut, "- Reopening Merge request !%d...\n", mr.IID)

			mr, err = api.UpdateMR(apiClient, repo.FullName(), mr.IID, l)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.StdOut, "%s Merge request !%d reopened\n", utils.GreenCheck(), mr.IID)
			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(mr))

			return nil
		},
	}

	return mrReopenCmd
}
