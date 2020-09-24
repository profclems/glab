package mr

import (
	"fmt"
	"log"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var mrRebaseCmd = &cobra.Command{
	Use:     "rebase <id> [flags]",
	Short:   `Automatically rebase the source_branch of the merge request against its target_branch.`,
	Long:    `If you don’t have permissions to push to the merge request’s source branch - you’ll get a 403 Forbidden response.`,
	Aliases: []string{"accept"},
	Args:    cobra.ExactArgs(1),
	Run:     acceptRebaseRequest,
}

func acceptRebaseRequest(cmd *cobra.Command, args []string) {
	mergeID := strings.Trim(args[0], " ")
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo, _ = fixRepoNamespace(r)
	}
	fmt.Println("Sending request...")
	_, err := gitlabClient.MergeRequests.RebaseMergeRequest(repo, manip.StringToInt(mergeID))
	if err != nil {
		er(err)
		return
	}

	opts := &gitlab.GetMergeRequestsOptions{}
	opts.IncludeRebaseInProgress = gitlab.Bool(true)
	fmt.Println("Checking rebase status...")
	i := 0
	for {
		mr, _, err := gitlabClient.MergeRequests.GetMergeRequest(repo, manip.StringToInt(mergeID), opts)
		if err != nil {
			log.Fatal(err)
		}
		if mr.RebaseInProgress {
			if i == 0 {
				fmt.Println("Rebase in progress...")
			}
		} else {
			if mr.MergeError != "" && mr.MergeError != "null" {
				fmt.Println(mr.MergeError)
				break
			}
			fmt.Println("Rebase successful")
			break
		}
		i++
	}

}

func init() {
	mrCmd.AddCommand(mrRebaseCmd)
}
