package update

import (
	"errors"
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
			var actions []string

			if cmd.Flags().Changed("lock-discussion") && cmd.Flags().Changed("unlock-discussion") {
				return &cmdutils.FlagError{Err: errors.New("--lock-discussion and --unlock-discussion can't be used together")}
			}

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
				actions = append(actions, fmt.Sprintf("updated title to %q", m))
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
					actions = append(actions, "marked as Draft")
					mergeTitle = "Draft: " + mergeTitle
				} else {
					actions = append(actions, "marked as WIP")
					mergeTitle = "WIP: " + mergeTitle
				}
			} else if isReady, _ := cmd.Flags().GetBool("ready"); isReady {
				actions = append(actions, "marked as ready")
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
				actions = append(actions, "locked discussion")
				l.DiscussionLocked = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("unlock-discussion"); m {
				actions = append(actions, "unlocked discussion")
				l.DiscussionLocked = gitlab.Bool(false)
			}

			if m, _ := cmd.Flags().GetString("description"); m != "" {
				actions = append(actions, "updated description")
				l.Description = gitlab.String(m)
			}

			if assignees, _ := cmd.Flags().GetStringSlice("assignees"); len(assignees) > 0 {
				users, err := api.UsersByNames(apiClient, assignees)
				if err != nil {
					return err
				}
				l.AssigneeIDs = cmdutils.IDsFromUsers(users)
				var usernames []string
				for i := range users {
					usernames = append(usernames, fmt.Sprintf("@%s", users[i].Username))
				}
				actions = append(actions, "assigned to", strings.Join(usernames, " "))
			}

			if removeSource, _ := cmd.Flags().GetBool("remove-source-branch"); removeSource {
				actions = append(actions, "enabled removal of source branch on merge")
				l.RemoveSourceBranch = gitlab.Bool(true)
			}

			fmt.Fprintf(f.IO.StdOut, "- Updating merge request !%d\n", mr.IID)

			mr, err = api.UpdateMR(apiClient, repo.FullName(), mr.IID, l)
			if err != nil {
				return err
			}

			for _, s := range actions {
				fmt.Fprintln(f.IO.StdOut, utils.GreenCheck(), s)
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
	mrUpdateCmd.Flags().BoolP("unlock-discussion", "", false, "Unlock discussion on merge request")
	mrUpdateCmd.Flags().StringP("description", "d", "", "merge request description")
	mrUpdateCmd.Flags().StringSliceP("assignees", "a", []string{}, "Assign merge request to people by their `usernames`")
	mrUpdateCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")

	return mrUpdateCmd
}
