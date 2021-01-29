package merge

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdMerge(f *cmdutils.Factory) *cobra.Command {
	var mrMergeCmd = &cobra.Command{
		Use:     "merge {<id> | <branch>}",
		Short:   `Merge/Accept merge requests`,
		Long:    ``,
		Aliases: []string{"accept"},
		Example: heredoc.Doc(`
			$ glab mr merge 235
			$ glab mr accept 235
			$ glab mr merge    # Finds open merge request from current branch
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

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
				Conflict:       true,
				PipelineStatus: true,
				MergePrivilege: true,
			}); err != nil {
				return err
			}

			opts := &gitlab.AcceptMergeRequestOptions{}
			if m, _ := cmd.Flags().GetString("message"); m != "" {
				opts.MergeCommitMessage = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetString("squash-message"); m != "" {
				opts.SquashCommitMessage = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetBool("squash"); m {
				opts.Squash = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("remove-source-branch"); m {
				opts.ShouldRemoveSourceBranch = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("when-pipeline-succeeds"); m {
				opts.MergeWhenPipelineSucceeds = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetString("sha"); m != "" {
				opts.SHA = gitlab.String(m)
			}

			fmt.Fprintf(f.IO.StdOut, "- Merging merge request !%d\n", mr.IID)

			mr, err = api.MergeMR(apiClient, repo.FullName(), mr.IID, opts)

			if err != nil {
				return err
			}

			if m, _ := cmd.Flags().GetBool("when-pipeline-succeeds"); m {
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Will merge when pipeline succeeds")
			} else {
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Merged")
			}
			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr))

			return nil
		},
	}

	mrMergeCmd.Flags().StringP("sha", "", "", "Merge Commit sha")
	mrMergeCmd.Flags().BoolP("remove-source-branch", "d", false, "Remove source branch on merge")
	mrMergeCmd.Flags().BoolP("when-pipeline-succeeds", "", true, "Merge only when pipeline succeeds")
	mrMergeCmd.Flags().StringP("message", "m", "", "Custom merge commit message")
	mrMergeCmd.Flags().StringP("squash-message", "", "", "Custom Squash commit message")
	mrMergeCmd.Flags().BoolP("squash", "s", false, "Squash commits on merge")

	return mrMergeCmd
}
