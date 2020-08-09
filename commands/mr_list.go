package commands

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var mrListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: `List merge requests`,
	Long:  ``,
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(3),
	RunE: listMergeRequest,
}

func listMergeRequest(cmd *cobra.Command, args []string) error {
	var state string
	if lb, _ := cmd.Flags().GetBool("all"); lb  {
		state = "all"
	} else if lb, _ := cmd.Flags().GetBool("closed"); lb  {
		state = "closed"
	} else if lb, _ := cmd.Flags().GetBool("merged"); lb  {
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
	if lb, _ := cmd.Flags().GetString("milestone"); lb != ""  {
		l.Milestone = gitlab.String(lb)
	}

	gitlabClient, repo := git.InitGitlabClient()
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
	mrCmd.AddCommand(mrListCmd)
}
