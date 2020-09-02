package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var mrApproveCmd = &cobra.Command{
	Use:     "approve <id> [flags]",
	Short:   `Approve merge requests`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run:     approveMergeRequest,
}

func approveMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		l := &gitlab.ApproveMergeRequestOptions{}
		if s, _ := cmd.Flags().GetString("sha"); s != "" {
			l.SHA = gitlab.String(s)
		}
		//if s, _ := cmd.Flags().GetString("password"); s  {
		// ToDo:
		//}

		fmt.Println(color.Yellow.Sprint("Approving Merge Request #" + mergeID + "..."))
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		_, resp, _ := gitlabClient.MergeRequestApprovals.ApproveMergeRequest(repo, manip.StringToInt(mergeID), l)
		if resp != nil {
			if resp.StatusCode == 201 {
				fmt.Println(color.Green.Sprint("Merge Request approved successfully"))
			} else if resp.StatusCode == 405 {
				fmt.Println("Merge request cannot be approved")
			} else if resp.StatusCode == 401 {
				fmt.Println("Merge request already approved or you don't have enough permission to approve this merge request")
			} else {
				fmt.Println(resp.Status)
			}
		} else {
			fmt.Println(resp)
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrApproveCmd.Flags().StringP("sha", "s", "", "The HEAD of the merge request")
	//mrApproveCmd.Flags().StringP("password", "p", "", "Current user’s password. Required if 'Require user password to approve' is enabled in the project settings.")
	mrCmd.AddCommand(mrApproveCmd)
}
