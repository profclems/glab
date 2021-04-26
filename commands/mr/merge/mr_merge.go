package merge

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type MergeOpts struct {
	MergeWhenPipelineSucceeds bool
	SquashBeforeMerge         bool
	RemoveSourceBranch        bool

	SquashMessage      string
	MergeCommitMessage string
	SHA                string
}

func NewCmdMerge(f *cmdutils.Factory) *cobra.Command {
	var opts = &MergeOpts{}

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

			mr, repo, err := mrutils.MRFromArgs(f, args, "opened")
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

			if !cmd.Flags().Changed("when-pipeline-succeeds") && f.IO.IsOutputTTY() && mr.Pipeline != nil {
				_ = prompt.Confirm(&opts.MergeWhenPipelineSucceeds, "Merge when pipeline succeeds?", true)
			}

			mergeOpts := &gitlab.AcceptMergeRequestOptions{}
			if opts.MergeCommitMessage != "" {
				mergeOpts.MergeCommitMessage = gitlab.String(opts.MergeCommitMessage)
			}
			if opts.SquashMessage != "" {
				mergeOpts.SquashCommitMessage = gitlab.String(opts.SquashMessage)
			}
			if opts.SquashBeforeMerge {
				mergeOpts.Squash = gitlab.Bool(true)
			}
			if opts.RemoveSourceBranch {
				mergeOpts.ShouldRemoveSourceBranch = gitlab.Bool(true)
			}
			if opts.MergeWhenPipelineSucceeds && mr.Pipeline != nil {
				mergeOpts.MergeWhenPipelineSucceeds = gitlab.Bool(true)
			}
			if opts.SHA != "" {
				mergeOpts.SHA = gitlab.String(opts.SHA)
			}

			fmt.Fprintf(f.IO.StdOut, "- Merging merge request !%d\n", mr.IID)

			mr, err = api.MergeMR(apiClient, repo.FullName(), mr.IID, mergeOpts)

			if err != nil {
				return err
			}

			isMerged := true
			if opts.MergeWhenPipelineSucceeds {
				if mr.Pipeline == nil {
					fmt.Fprintln(f.IO.StdOut, c.WarnIcon(), "No pipeline running on", mr.SourceBranch)
				} else if mr.Pipeline.Status == "success" {
					fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Pipeline Succeeded")
				} else {
					fmt.Fprintln(f.IO.StdOut, c.WarnIcon(), "Pipeline Status:", mr.Pipeline.Status)
					fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Will merge when pipeline succeeds")
					isMerged = false
				}
			}
			if isMerged {
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Merged")
			}
			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(c, mr))

			return nil
		},
	}

	mrMergeCmd.Flags().StringVarP(&opts.SHA, "sha", "", "", "Merge Commit sha")
	mrMergeCmd.Flags().BoolVarP(&opts.RemoveSourceBranch, "remove-source-branch", "d", false, "Remove source branch on merge")
	mrMergeCmd.Flags().BoolVarP(&opts.MergeWhenPipelineSucceeds, "when-pipeline-succeeds", "", true, "Merge only when pipeline succeeds")
	mrMergeCmd.Flags().StringVarP(&opts.MergeCommitMessage, "message", "m", "", "Custom merge commit message")
	mrMergeCmd.Flags().StringVarP(&opts.SquashMessage, "squash-message", "", "", "Custom Squash commit message")
	mrMergeCmd.Flags().BoolVarP(&opts.SquashBeforeMerge, "squash", "s", false, "Squash commits on merge")

	return mrMergeCmd
}
