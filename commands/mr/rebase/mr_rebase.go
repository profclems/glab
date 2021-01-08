package rebase

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
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

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, "- Sending request...")
			err = api.RebaseMR(apiClient, repo.FullName(), mr.IID)
			if err != nil {
				return err
			}

			opts := &gitlab.GetMergeRequestsOptions{}
			opts.IncludeRebaseInProgress = gitlab.Bool(true)
			fmt.Fprintln(f.IO.StdOut, "- Checking rebase status...")
			i := 0
			for {
				mr, err := api.GetMR(apiClient, repo.FullName(), mr.IID, opts)
				if err != nil {
					return err
				}
				if mr.RebaseInProgress {
					if i == 0 {
						fmt.Fprintln(f.IO.StdOut, "- Rebase in progress...")
					}
				} else {
					if mr.MergeError != "" && mr.MergeError != "null" {
						fmt.Fprintln(f.IO.StdErr, mr.MergeError)
						break
					}
					fmt.Fprintln(f.IO.StdOut, utils.GreenCheck(), "Rebase successful")
					break
				}
				i++
			}

			return nil
		},
	}

	return mrRebaseCmd
}
