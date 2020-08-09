package commands

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var mrReopenCmd = &cobra.Command{
	Use:   "reopen <id>",
	Short: `Reopen merge requests`,
	Long:  ``,
	Aliases: []string{"open"},
	Args:    cobra.MaximumNArgs(1),
	Run: reopenMergeRequestState,
}

func reopenMergeRequestState(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		l := &gitlab.UpdateMergeRequestOptions{}
		l.StateEvent = gitlab.String("reopen")
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Printf("Updating Merge request #%s...\n", i2)
			mr, resp, _ := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(i2), l)
			if resp.StatusCode == 200 {
				fmt.Println(aurora.Green("You have reopened merge request #" + i2))
				displayMergeRequest(mr)
			} else if resp.StatusCode == 404 {
				er("MergeRequest does not exist")
			} else {
				er("Could not complete request: "+resp.Status)
			}
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrReopenCmd)
}