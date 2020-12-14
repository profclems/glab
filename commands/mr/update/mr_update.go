package update

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdUpdate(f *cmdutils.Factory) *cobra.Command {
	var mrUpdateCmd = &cobra.Command{
		Use:   "update [<id> | <branch>]",
		Short: `Update merge requests`,
		Long:  ``,
		Example: heredoc.Doc(`
	$ glab mr update 23 --ready
	$ glab mr update 23 --draft
	$ glab mr update --draft  # Updates MR related to current branch
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

			l := &gitlab.UpdateMergeRequestOptions{}
			var mergeTitle string

			isDraft, _ := cmd.Flags().GetBool("draft")
			isWIP, _ := cmd.Flags().GetBool("wip")
			if m, _ := cmd.Flags().GetString("title"); m != "" {
				mergeTitle = m
			}
			if mergeTitle == "" {
				opts := &gitlab.GetMergeRequestsOptions{}
				mr, err := api.GetMR(apiClient, repo.FullName(), mr.IID, opts)
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

			if assignees, _ := cmd.Flags().GetStringSlice("assignees"); len(assignees) > 0 {
				users, err := api.UsersByNames(apiClient, assignees)
				if err != nil {
					return err
				}
				l.AssigneeIDs = cmdutils.IDsFromUsers(users)
			}

			if removeSource, _ := cmd.Flags().GetBool("remove-source-branch"); removeSource {
				l.RemoveSourceBranch = gitlab.Bool(true)
			}

			mr, err = api.UpdateMR(apiClient, repo.FullName(), mr.IID, l)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(mr))
			return nil
		},
	}

	mrUpdateCmd.Flags().BoolP("draft", "", false, "Mark merge request as a draft")
	mrUpdateCmd.Flags().BoolP("ready", "r", false, "Mark merge request as ready to be reviewed and merged")
	mrUpdateCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrUpdateCmd.Flags().StringP("title", "t", "", "Title of merge request")
	mrUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on merge request")
	mrUpdateCmd.Flags().StringP("description", "d", "", "merge request description")
	mrUpdateCmd.Flags().StringSliceP("assignees", "a", []string{}, "Assign merge request to people by their `usernames`")
	mrUpdateCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")

	return mrUpdateCmd
}
