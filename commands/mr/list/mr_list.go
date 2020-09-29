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

			if lb, _ := cmd.Flags().GetBool("all"); lb {
				state = "all"
			} else if lb, _ := cmd.Flags().GetBool("closed"); lb {
				state = "closed"
			} else if lb, _ := cmd.Flags().GetBool("merged"); lb {
				state = "merged"
			} else {
				state = "opened"
			}

			l := &gitlab.ListProjectMergeRequestsOptions{
				State: gitlab.String(state),
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
			if p, _ := cmd.Flags().GetInt("page"); p != 0 {
				l.Page = p
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				l.PerPage = p
			}

			if mine, _ := cmd.Flags().GetBool("mine"); mine {
				l.Scope = gitlab.String("assigned_to_me")
			}

			mergeRequests, err := api.ListMRs(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, mrutils.DisplayAllMRs(mergeRequests, repo.FullName()))
			return nil
		},
	}

	mrListCmd.Flags().StringP("label", "l", "", "Filter merge request by label <name>")
	mrListCmd.Flags().StringP("milestone", "", "", "Filter merge request by milestone <id>")
	mrListCmd.Flags().BoolP("all", "a", false, "Get all merge requests")
	mrListCmd.Flags().BoolP("closed", "c", false, "Get only closed merge requests")
	mrListCmd.Flags().BoolP("opened", "o", false, "Get only opened merge requests")
	mrListCmd.Flags().BoolP("merged", "m", false, "Get only merged merge requests")
	mrListCmd.Flags().IntP("page", "p", 1, "Page number")
	mrListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	mrListCmd.Flags().BoolP("mine", "", false, "Get only merge requests assigned to me")

	return mrListCmd
}
