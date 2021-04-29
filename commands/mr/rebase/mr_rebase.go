package rebase

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
)

func NewCmdRebase(f *cmdutils.Factory) *cobra.Command {
	var mrRebaseCmd = &cobra.Command{
		Use:   "rebase [<id> | <branch>] [flags]",
		Short: `Automatically rebase the source_branch of the merge request against its target_branch.`,
		Long:  `If you don’t have permissions to push to the merge request’s source branch - you’ll get a 403 Forbidden response.`,
		Example: heredoc.Doc(`
			$ glab mr rebase 123
			$ glab mr rebase  # get from current branch
			$ glab mr rebase branch
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args, "opened")
			if err != nil {
				return err
			}

			if err = mrutils.RebaseMR(f.IO, apiClient, repo, mr); err != nil {
				return err
			}

			return nil
		},
	}

	return mrRebaseCmd
}
