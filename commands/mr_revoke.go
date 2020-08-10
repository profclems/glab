package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/logrusorgru/aurora"

	"glab/internal/git"
	"glab/internal/manip"
)

var mrRevokeCmd = &cobra.Command{
	Use:     "revoke <id>",
	Short:   `Revoke approval on a merge request <id>`,
	Long:    ``,
	Aliases: []string{"unapprove"},
	Run:     revokeMergeRequest,
}

func revokeMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")

		fmt.Println(aurora.Yellow("Revoking approval for Merge Request #" + mergeID + "..."))
		gitlabClient, repo := git.InitGitlabClient()
		resp, _ := gitlabClient.MergeRequestApprovals.UnapproveMergeRequest(repo, manip.StringToInt(mergeID))
		if resp != nil {
			if resp.StatusCode == 201 {
				fmt.Println(aurora.Green("Merge Request approval revoked successfully"))
			} else if resp.StatusCode == 405 {
				er("Merge request cannot be unapproved")
			} else if resp.StatusCode == 401 {
				er("Merge request already unapproved or you don't have enough permission to unapprove this merge request")
			} else {
				er(resp.Status)
			}
		} else {
			er(resp)
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrRevokeCmd)
}
