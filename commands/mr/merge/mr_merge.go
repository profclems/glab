package merge

import (
	"fmt"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdMerge(f *cmdutils.Factory) *cobra.Command {
	var mrMergeCmd = &cobra.Command{
		Use:     "merge <id> [flags]",
		Short:   `Merge/Accept merge requests`,
		Long:    ``,
		Aliases: []string{"accept"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mergeID := strings.Trim(args[0], " ")
			l := &gitlab.AcceptMergeRequestOptions{}
			if m, _ := cmd.Flags().GetString("message"); m != "" {
				l.MergeCommitMessage = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetString("squash-message"); m != "" {
				l.SquashCommitMessage = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetBool("squash"); m {
				l.Squash = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("remove-source-branch"); m {
				l.ShouldRemoveSourceBranch = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("when-pipeline-succeeds"); m {
				l.MergeWhenPipelineSucceeds = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetString("sha"); m != "" {
				l.SHA = gitlab.String(m)
			}

			fmt.Fprintf(out, "- Merging merge request #%s\n", mergeID)

			mr, err := api.MergeMR(apiClient, repo, manip.StringToInt(mergeID), l)

			if err != nil {
				return err
			}

			fmt.Fprintln(out, utils.GreenCheck(), "Merged")
			fmt.Fprintln(out, mrutils.DisplayMR(mr))

			return nil
		},
	}

	mrMergeCmd.Flags().StringP("sha", "", "", "Merge Commit sha")
	mrMergeCmd.Flags().BoolP("remove-source-branch", "d", false, "Remove source branch on merge")
	mrMergeCmd.Flags().BoolP("when-pipeline-succeeds", "", true, "Merge only when pipeline succeeds. Default to true")
	mrMergeCmd.Flags().StringP("message", "m", "", "Get only closed merge requests")
	mrMergeCmd.Flags().StringP("squash-message", "", "", "Squash commit message")
	mrMergeCmd.Flags().BoolP("squash", "s", false, "Squash commits on merge")

	return mrMergeCmd
}