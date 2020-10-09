package list

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var issueListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List project issues`,
		Long:    ``,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				state string
				err   error
			)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			if lb, _ := cmd.Flags().GetBool("all"); lb {
				state = "all"
			} else if lb, _ := cmd.Flags().GetBool("closed"); lb {
				state = "closed"
			} else {
				state = "opened"
			}

			opts := &gitlab.ListProjectIssuesOptions{
				State: gitlab.String(state),
			}
			if lb, _ := cmd.Flags().GetString("assignee"); lb != "" {
				opts.AssigneeUsername = gitlab.String(lb)
			}
			if lb, _ := cmd.Flags().GetString("label"); lb != "" {
				label := gitlab.Labels{
					lb,
				}
				opts.Labels = label
			}
			if lb, _ := cmd.Flags().GetString("milestone"); lb != "" {
				opts.Milestone = gitlab.String(lb)
			}
			if lb, _ := cmd.Flags().GetBool("confidential"); lb {
				opts.Confidential = gitlab.Bool(lb)
			}
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				opts.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				opts.PerPage = p
			}

			if lb, _ := cmd.Flags().GetBool("mine"); lb {
				u, _ := api.CurrentUser(nil)
				opts.AssigneeUsername = gitlab.String(u.Username)
			}
			issues, err := api.ListIssues(apiClient, repo.FullName(), opts)
			if err != nil {
				return err
			}
			fmt.Fprintln(utils.ColorableOut(cmd), issueutils.DisplayAllIssues(issues, repo.FullName()))
			return nil

		},
	}
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

	return issueListCmd
}
