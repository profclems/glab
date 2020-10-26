package revoke

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdRevoke(f *cmdutils.Factory) *cobra.Command {
	var mrRevokeCmd = &cobra.Command{
		Use:     "revoke [<id> | <branch>]",
		Short:   `Revoke approval on a merge request <id>`,
		Long:    ``,
		Aliases: []string{"unapprove"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
				WorkInProgress: true,
				Closed:         true,
				Merged:         true,
			}); err != nil {
				return err
			}

			fmt.Fprintf(out, "- Revoking approval for Merge Request #%d...\n", mr.IID)

			err = api.UnapproveMR(apiClient, repo.FullName(), mr.IID)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, utils.GreenCheck(), "Merge Request approval revoked")

			return nil
		},
	}

	return mrRevokeCmd
}
