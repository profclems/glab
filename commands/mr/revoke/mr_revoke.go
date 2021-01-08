package revoke

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
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
		Example: heredoc.Doc(`
			$ glab mr revoke 123
			$ glab mr unapprove 123
			$ glab mr revoke branch
			$ glab mr revoke  # use checked out branch
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

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

			fmt.Fprintf(f.IO.StdOut, "- Revoking approval for Merge Request #%d...\n", mr.IID)

			err = api.UnapproveMR(apiClient, repo.FullName(), mr.IID)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, utils.GreenCheck(), "Merge Request approval revoked")

			return nil
		},
	}

	return mrRevokeCmd
}
