package commands

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var issueListCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   `List project issues`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var state string
		if lb, _ := cmd.Flags().GetBool("all"); lb {
			state = "all"
		} else if lb, _ := cmd.Flags().GetBool("closed"); lb {
			state = "closed"
		} else {
			state = "opened"
		}

		l := &gitlab.ListProjectIssuesOptions{
			State: gitlab.String(state),
		}
		if lb, _ := cmd.Flags().GetString("assignee"); lb != "" {
			l.AssigneeUsername = gitlab.String(lb)
		}
		if lb, _ := cmd.Flags().GetString("label"); lb != "" {
			label := gitlab.Labels{
				lb,
			}
			l.Labels = label
		}
		if lb, _ := cmd.Flags().GetString("milestone"); lb != "" {
			l.Milestone = gitlab.String(lb)
		}
		if lb, _ := cmd.Flags().GetBool("confidential"); lb {
			l.Confidential = gitlab.Bool(lb)
		}
		if p, _ := cmd.Flags().GetInt("page"); p != 0 {
			l.Page = p
		}
		if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
			l.PerPage = p
		}
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		if lb, _ := cmd.Flags().GetBool("mine"); lb {
			u, _, _ := gitlabClient.Users.CurrentUser()
			l.AssigneeUsername = gitlab.String(u.Username)
		}
		issues, _, err := gitlabClient.Issues.ListProjectIssues(repo, l)
		if err != nil {
			return err
		}
		displayAllIssues(issues)
		return nil

	},
}

func init() {
	issueListCmd.Flags().StringP("assignee", "", "", "Filter issue by assignee <username>")
	issueListCmd.Flags().StringP("label", "l", "", "Filter issue by label <name>")
	issueListCmd.Flags().StringP("milestone", "", "", "Filter issue by milestone <id>")
	issueListCmd.Flags().BoolP("mine", "", false, "Filter only issues issues assigned to me")
	issueListCmd.Flags().BoolP("all", "a", false, "Get all issues")
	issueListCmd.Flags().BoolP("closed", "c", false, "Get only closed issues")
	issueListCmd.Flags().BoolP("opened", "o", false, "Get only opened issues")
	issueListCmd.Flags().BoolP("confidential", "", false, "Filter by confidential issues")
	issueListCmd.Flags().IntP("page", "p", 1, "Page number")
	issueListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	issueCmd.AddCommand(issueListCmd)
}
