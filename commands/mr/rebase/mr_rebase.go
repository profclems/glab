package rebase

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdRebase(f *cmdutils.Factory) *cobra.Command {
	var mrRebaseCmd = &cobra.Command{
		Use:     "rebase <id> [flags]",
		Short:   `Automatically rebase the source_branch of the merge request against its target_branch.`,
		Long:    `If you don’t have permissions to push to the merge request’s source branch - you’ll get a 403 Forbidden response.`,
		Aliases: []string{"accept"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mergeID := strings.TrimSpace(args[0])

			fmt.Fprintln(out, "- Sending request...")
			err = api.RebaseMR(apiClient, repo.FullName(), utils.StringToInt(mergeID))
			if err != nil {
				return err
			}

			opts := &gitlab.GetMergeRequestsOptions{}
			opts.IncludeRebaseInProgress = gitlab.Bool(true)
			fmt.Fprintln(out, "- Checking rebase status...")
			i := 0
			for {
				mr, err := api.GetMR(apiClient, repo.FullName(), utils.StringToInt(mergeID), opts)
				if err != nil {
					return err
				}
				if mr.RebaseInProgress {
					if i == 0 {
						fmt.Fprintln(out, "- Rebase in progress...")
					}
				} else {
					if mr.MergeError != "" && mr.MergeError != "null" {
						fmt.Fprintln(utils.ColorableErr(cmd), mr.MergeError)
						break
					}
					fmt.Fprintln(out, utils.GreenCheck(), "Rebase successful")
					break
				}
				i++
			}

			return nil
		},
	}

	return mrRebaseCmd
}
