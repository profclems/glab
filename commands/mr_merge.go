package commands

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"
)

var mrMergeCmd = &cobra.Command{
	Use:     "merge <id> [flags]",
	Short:   `Merge/Accept merge requests`,
	Long:    ``,
	Aliases: []string{"accept"},
	Args:    cobra.ExactArgs(1),
	Run:     acceptMergeRequest,
}

func acceptMergeRequest(cmd *cobra.Command, args []string) {
	mergeID := strings.Trim(args[0], " ")
	l := &gitlab.AcceptMergeRequestOptions{}
	if m, _ := cmd.Flags().GetString("message"); m != "" {
		l.MergeCommitMessage = gitlab.String(m)
	}
	if m, _ := cmd.Flags().GetString("squash-message"); m != "" {
		l.SquashCommitMessage = gitlab.String(m)
	}
	if m, _ := cmd.Flags().GetBool("squash"); m {
		l.Squash = gitlab.Bool(m)
	}
	if m, _ := cmd.Flags().GetBool("remove-source-branch"); m {
		l.ShouldRemoveSourceBranch = gitlab.Bool(m)
	}
	if m, _ := cmd.Flags().GetBool("when-pipeline-succeeds"); m {
		l.MergeWhenPipelineSucceeds = gitlab.Bool(m)
	}
	if m, _ := cmd.Flags().GetString("sha"); m != "" {
		l.SHA = gitlab.String(m)
	}
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo, _ = fixRepoNamespace(r)
	}
	mr, resp, _ := gitlabClient.MergeRequests.AcceptMergeRequest(repo, manip.StringToInt(mergeID), l)

	fmt.Println(aurora.Yellow("Accepting Merge Request #" + mergeID + "..."))

	if resp.StatusCode == 200 {
		fmt.Println(aurora.Green("Merge Request accepted successfully"))
		displayMergeRequest(mr)
	} else if resp.StatusCode == 405 {
		fmt.Println("Merge request cannot be merged")
	} else if resp.StatusCode == 401 {
		fmt.Println("You don't have enough permission to accept this merge request")
	} else if resp.StatusCode == 406 {
		fmt.Println("Branch cannot be merged. There are merge conflicts.")
	} else {
		fmt.Println(resp)
	}
}

func init() {
	mrMergeCmd.Flags().StringP("sha", "", "", "Merge Commit sha")
	mrMergeCmd.Flags().BoolP("remove-source-branch", "d", false, "Remove source branch on merge")
	mrMergeCmd.Flags().BoolP("when-pipeline-succeeds", "", true, "Merge only when pipeline succeeds. Default to true")
	mrMergeCmd.Flags().StringP("message", "m", "", "Get only closed merge requests")
	mrMergeCmd.Flags().StringP("squash-message", "", "", "Squash commit message")
	mrMergeCmd.Flags().BoolP("squash", "s", false, "Squash commits on merge")
	mrCmd.AddCommand(mrMergeCmd)
}
