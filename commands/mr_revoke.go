package commands

import (
	"fmt"
	"github.com/gookit/color"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
)

var mrRevokeCmd = &cobra.Command{
	Use:     "revoke <id>",
	Short:   `Revoke approval on a merge request <id>`,
	Long:    ``,
	Aliases: []string{"unapprove"},
	Args:    cobra.ExactArgs(1),
	Run:     revokeMergeRequest,
}

func revokeMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")

		fmt.Println(color.Yellow.Sprint("Revoking approval for Merge Request #" + mergeID + "..."))
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		resp, _ := gitlabClient.MergeRequestApprovals.UnapproveMergeRequest(repo, manip.StringToInt(mergeID))
		if resp != nil {
			if resp.StatusCode == 201 {
				fmt.Println(color.Green.Sprint("Merge Request approval revoked successfully"))
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
