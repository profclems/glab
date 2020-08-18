package commands

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var mrListCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   `List merge requests`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(0),
	RunE:    listMergeRequest,
}

func listMergeRequest(cmd *cobra.Command, args []string) error {
	var state string
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
		l.Labels = &label
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

	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo = r
	}
	mergeRequests, _, err := gitlabClient.MergeRequests.ListProjectMergeRequests(repo, l)
	if err != nil {
		return err
	}
	displayAllMergeRequests(mergeRequests)
	return nil
}

func init() {
	mrListCmd.Flags().StringP("label", "l", "", "Filter merge request by label <name>")
	mrListCmd.Flags().StringP("milestone", "", "", "Filter merge request by milestone <id>")
	mrListCmd.Flags().BoolP("all", "a", false, "Get all merge requests")
	mrListCmd.Flags().BoolP("closed", "c", false, "Get only closed merge requests")
	mrListCmd.Flags().BoolP("opened", "o", false, "Get only opened merge requests")
	mrListCmd.Flags().BoolP("merged", "m", false, "Get only merged merge requests")
	mrListCmd.Flags().IntP("page", "p", 1, "Page number")
	mrListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	mrCmd.AddCommand(mrListCmd)
}
