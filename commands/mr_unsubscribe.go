package commands

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"glab/internal/git"
	"glab/internal/manip"
)

var mrUnsubscribeCmd = &cobra.Command{
	Use:   "unsubscribe <id>",
	Short: `Unsubscribe to merge requests`,
	Long:  ``,
	Aliases: []string{"unsub"},
	Run: unsubscribeMergeRequest,
}

func unsubscribeMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #" + i2)
			mr, resp, _ := gitlabClient.MergeRequests.UnsubscribeFromMergeRequest(repo, manip.StringToInt(i2))

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully unsubscribed to merge request #" + i2))
				displayMergeRequest(mr)
			} else if resp.StatusCode == 404 {
				er("MergeRequest does not exist")
			} else {
				er("Could not complete request." + resp.Status)
			}
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrUnsubscribeCmd)
}