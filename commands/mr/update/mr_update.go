package update

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdUpdate(f *cmdutils.Factory) *cobra.Command {
	var mrUpdateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: `Update merge requests`,
		Long:  ``,
		Example: heredoc.Doc(`
	$ glab mr update 23 --ready
	$ glab mr update 23 --draft
	`),
		Args: cobra.ExactArgs(1),
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

			mergeID := utils.StringToInt(args[0])
			l := &gitlab.UpdateMergeRequestOptions{}
			var mergeTitle string

			isDraft, _ := cmd.Flags().GetBool("draft")
			isWIP, _ := cmd.Flags().GetBool("wip")
			if m, _ := cmd.Flags().GetString("title"); m != "" {
				mergeTitle = m
			}
			if mergeTitle == "" {
				opts := &gitlab.GetMergeRequestsOptions{}
				mr, err := api.GetMR(apiClient, repo.FullName(), mergeID, opts)
				if err != nil {
					return err
				}
				mergeTitle = mr.Title
			}
			if isDraft || isWIP {
				if isDraft {
					mergeTitle = "Draft: " + mergeTitle
				} else {
					mergeTitle = "WIP: " + mergeTitle
				}
			} else if isReady, _ := cmd.Flags().GetBool("ready"); isReady {
				mergeTitle = strings.TrimPrefix(mergeTitle, "Draft:")
				mergeTitle = strings.TrimPrefix(mergeTitle, "draft:")
				mergeTitle = strings.TrimPrefix(mergeTitle, "DRAFT:")
				mergeTitle = strings.TrimPrefix(mergeTitle, "WIP:")
				mergeTitle = strings.TrimPrefix(mergeTitle, "wip:")
				mergeTitle = strings.TrimPrefix(mergeTitle, "Wip:")
				mergeTitle = strings.TrimSpace(mergeTitle)
			}
			l.Title = gitlab.String(mergeTitle)
			if m, _ := cmd.Flags().GetBool("lock-discussion"); m {
				l.DiscussionLocked = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetString("description"); m != "" {
				l.Description = gitlab.String(m)
			}
			mr, err := api.UpdateMR(apiClient, repo.FullName(), mergeID, l)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, mrutils.DisplayMR(mr))
			return nil
		},
	}

	mrUpdateCmd.Flags().BoolP("draft", "", false, "Mark merge request as a draft")
	mrUpdateCmd.Flags().BoolP("ready", "r", false, "Mark merge request as ready to be reviewed and merged")
	mrUpdateCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrUpdateCmd.Flags().StringP("title", "t", "", "Title of merge request")
	mrUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on merge request")
	mrUpdateCmd.Flags().StringP("description", "d", "", "merge request description")

	return mrUpdateCmd
}
