package update

import (
	"errors"
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdUpdate(f *cmdutils.Factory) *cobra.Command {
	var issueUpdateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: `Update issue`,
		Long:  ``,
		Example: heredoc.Doc(`
	$ glab issue update 42 --label ui,ux
	$ glab issue update 42 --unlabel working
	`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var actions []string
			var assigneesToAdd, assigneesToRemove, assigneesToReplace []string
			out := f.IO.StdOut

			// Parse assignees Early so we can fail early in case of conflicts
			if cmd.Flags().Changed("assignee") {
				givenAssignees, err := cmd.Flags().GetStringSlice("assignee")
				if err != nil {
					return err
				}
				assigneesToAdd, assigneesToRemove, assigneesToReplace = cmdutils.ParseAssignees(givenAssignees)

				// Fail if relative and absolute assignees were given, there is no reason to mix them.
				if len(assigneesToReplace) != 0 && (len(assigneesToAdd) != 0 || len(assigneesToRemove) != 0) {
					return &cmdutils.FlagError{
						Err: errors.New("--assignee doesn't allow mixing relative (+,!,-) and absolute assignments"),
					}
				}

				if m := utils.CommonElementsInStringSlice(assigneesToAdd, assigneesToRemove); len(m) != 0 {
					return &cmdutils.FlagError{
						Err: fmt.Errorf("--assignee has %s %q that are present in both add and remove",
							utils.Pluralize(len(m), "element"),
							strings.Join(m, " ")),
					}
				}
			}

			if cmd.Flags().Changed("lock-discussion") && cmd.Flags().Changed("unlock-discussion") {
				return &cmdutils.FlagError{Err: errors.New("--lock-discussion and --unlock-discussion can't be used together")}
			}
			if cmd.Flags().Changed("confidential") && cmd.Flags().Changed("public") {
				return &cmdutils.FlagError{Err: errors.New("--public and --confidential can't be used together")}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			issueID := utils.StringToInt(args[0])
			l := &gitlab.UpdateIssueOptions{}

			if m, _ := cmd.Flags().GetString("title"); m != "" {
				actions = append(actions, fmt.Sprintf("updated title to %q", m))
				l.Title = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetBool("lock-discussion"); m {
				actions = append(actions, "locked discussion")
				l.DiscussionLocked = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetBool("unlock-discussion"); m {
				actions = append(actions, "unlocked dicussion")
				l.DiscussionLocked = gitlab.Bool(false)
			}

			if m, _ := cmd.Flags().GetString("description"); m != "" {
				actions = append(actions, "updated description")
				l.Description = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetStringSlice("label"); len(m) != 0 {
				actions = append(actions, fmt.Sprintf("added labels %s", strings.Join(m, " ")))
				l.AddLabels = gitlab.Labels(m)
			}
			if m, _ := cmd.Flags().GetStringSlice("unlabel"); len(m) != 0 {
				actions = append(actions, fmt.Sprintf("removed labels %s", strings.Join(m, " ")))
				l.RemoveLabels = gitlab.Labels(m)
			}
			if m, _ := cmd.Flags().GetBool("public"); m {
				actions = append(actions, "made public")
				l.Confidential = gitlab.Bool(false)
			}
			if m, _ := cmd.Flags().GetBool("confidential"); m {
				actions = append(actions, "made confidential")
				l.Confidential = gitlab.Bool(true)
			}
			if ok := cmd.Flags().Changed("milestone"); ok {
				if m, _ := cmd.Flags().GetString("milestone"); m != "" || m == "0" {
					mID, err := cmdutils.ParseMilestone(apiClient, repo, m)
					if err != nil {
						return err
					}
					actions = append(actions, fmt.Sprintf("added milestone %q", m))
					l.MilestoneID = gitlab.Int(mID)
				} else {
					// Unassign the Milestone
					actions = append(actions, "unassigned milestone")
					l.MilestoneID = gitlab.Int(0)
				}
			}
			if len(assigneesToReplace) != 0 {
				users, err := api.UsersByNames(apiClient, assigneesToReplace)
				if err != nil {
					return err
				}
				l.AssigneeIDs = cmdutils.IDsFromUsers(users)
				var usernames []string
				for i := range users {
					usernames = append(usernames, fmt.Sprintf("@%s", users[i].Username))
				}
				actions = append(actions, fmt.Sprintf("assigned to %s", strings.Join(usernames, " ")))
			} else if len(assigneesToAdd) != 0 || len(assigneesToRemove) != 0 {
				// Get List of assignees and store all of their IDs except for the ones
				// that have their `Username` match one of the usernames present in
				// `assigneesToRemove`
				issue, err := api.GetIssue(apiClient, repo.FullName(), issueID)
				if err != nil {
					return err
				}
				var assignedIDs []int
				for i := range issue.Assignees {
					// Only store them in assigneedIDs if they are not marked for removal
					if !utils.PresentInStringSlice(assigneesToRemove, issue.Assignees[i].Username) {
						assignedIDs = append(assignedIDs, issue.Assignees[i].ID)
					}
				}
				if len(assigneesToRemove) != 0 {
					actions = append(actions, fmt.Sprintf("unassigned %s", strings.Join(assigneesToRemove, "@ ")))
				}

				if len(assigneesToAdd) != 0 {
					users, err := api.UsersByNames(apiClient, assigneesToAdd)
					if err != nil {
						return err
					}
					assignedIDs = append(assignedIDs, cmdutils.IDsFromUsers(users)...)
					actions = append(actions, fmt.Sprintf("assigned %s", strings.Join(assigneesToAdd, "@ ")))
				}

				if len(assignedIDs) == 0 {
					l.AssigneeIDs = []int{0}
				} else {
					l.AssigneeIDs = assignedIDs
				}
			}

			fmt.Fprintf(out, "- Updating issue #%d\n", issueID)

			issue, err := api.UpdateIssue(apiClient, repo.FullName(), issueID, l)
			if err != nil {
				return err
			}

			for _, s := range actions {
				fmt.Fprintln(out, utils.GreenCheck(), s)
			}

			fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			return nil
		},
	}

	issueUpdateCmd.Flags().StringP("title", "t", "", "Title of issue")
	issueUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on issue")
	issueUpdateCmd.Flags().BoolP("unlock-discussion", "", false, "Unlock discussion on issue")
	issueUpdateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueUpdateCmd.Flags().StringSliceP("label", "l", []string{}, "add labels")
	issueUpdateCmd.Flags().StringSliceP("unlabel", "u", []string{}, "remove labels")
	issueUpdateCmd.Flags().BoolP("public", "p", false, "Make issue public")
	issueUpdateCmd.Flags().BoolP("confidential", "c", false, "Make issue confidential")
	issueUpdateCmd.Flags().StringP("milestone", "m", "", "title of the milestone to assign, pass \"\" or 0 to unassign")
	issueUpdateCmd.Flags().StringSliceP("assignee", "a", []string{}, "assign users via username, prefix with '!' or '-' to remove from existing assignees, '+' to add, otherwise replace existing assignees with given users")

	return issueUpdateCmd
}
