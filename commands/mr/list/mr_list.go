package list

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var mrListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List merge requests`,
		Long:    ``,
		Aliases: []string{"ls"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var state string
			var err error
			var listType string
			var titleQualifier string
			var mergeRequests []*gitlab.MergeRequest

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
				titleQualifier = state
			} else if lb, _ := cmd.Flags().GetBool("merged"); lb {
				state = "merged"
				titleQualifier = state
			} else {
				state = "opened"
				titleQualifier = "open"
			}

			l := &gitlab.ListProjectMergeRequestsOptions{
				State: gitlab.String(state),
			}
			l.Page = 1
			l.PerPage = 30

			if lb, _ := cmd.Flags().GetString("label"); lb != "" {
				label := gitlab.Labels{
					lb,
				}
				l.Labels = label
				listType = "search"
			}
			if lb, _ := cmd.Flags().GetString("milestone"); lb != "" {
				l.Milestone = gitlab.String(lb)
				listType = "search"
			}
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				l.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				l.PerPage = p
			}

			if mine, _ := cmd.Flags().GetBool("mine"); mine {
				l.Scope = gitlab.String("assigned_to_me")
				listType = "search"
			}

			assigneeIds := make([]int, 0)
			if assigneeNames, _ := cmd.Flags().GetStringSlice("assignee"); len(assigneeNames) > 0 {
				users, err := api.UsersByNames(apiClient, assigneeNames)
				if err != nil {
					return err
				}
				for _, user := range users {
					assigneeIds = append(assigneeIds, user.ID)
				}
			}

			if len(assigneeIds) > 0 {
				mergeRequests, err = api.ListMRsWithAssignees(apiClient, repo.FullName(), l, assigneeIds)

			} else {
				mergeRequests, err = api.ListMRs(apiClient, repo.FullName(), l)
			}
			if err != nil {
				return err
			}

			title := utils.NewListTitle(titleQualifier + " merge request")
			title.RepoName = repo.FullName()
			title.Page = l.Page
			title.ListActionType = listType
			title.CurrentPageTotal = len(mergeRequests)

			if err = f.IO.StartPager(); err != nil {
				return err
			}
			defer f.IO.StopPager()
			fmt.Fprintf(f.IO.StdOut, "%s\n%s\n", title.Describe(), mrutils.DisplayAllMRs(mergeRequests, repo.FullName()))

			return nil
		},
	}

	mrListCmd.Flags().StringP("label", "l", "", "Filter merge request by label <name>")
	mrListCmd.Flags().StringP("milestone", "", "", "Filter merge request by milestone <id>")
	mrListCmd.Flags().BoolP("all", "a", false, "Get all merge requests")
	mrListCmd.Flags().BoolP("closed", "c", false, "Get only closed merge requests")
	mrListCmd.Flags().BoolP("opened", "o", false, "Get only open merge requests")
	mrListCmd.Flags().BoolP("merged", "m", false, "Get only merged merge requests")
	mrListCmd.Flags().IntP("page", "p", 1, "Page number")
	mrListCmd.Flags().IntP("per-page", "P", 30, "Number of items to list per page. (default 30)")
	mrListCmd.Flags().BoolP("mine", "", false, "Get only merge requests assigned to me")
	mrListCmd.Flags().StringSliceP("assignee", "", []string{}, "Get only merge requests assigned to users")

	return mrListCmd
}
