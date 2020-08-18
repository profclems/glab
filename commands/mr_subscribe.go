package commands

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"glab/internal/git"
	"glab/internal/manip"
)

var mrSubscribeCmd = &cobra.Command{
	Use:     "subscribe <id>",
	Short:   `Subscribe to merge requests`,
	Long:    ``,
	Aliases: []string{"sub"},
	Args:    cobra.ExactArgs(1),
	Run:     subscribeSubscribeRequest,
}

func subscribeSubscribeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Subscribing to merge Request #" + i2)
			issue, resp, _ := gitlabClient.MergeRequests.SubscribeToMergeRequest(repo, manip.StringToInt(i2), nil)

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("You have successfully subscribed to merge request #" + i2))
				displayMergeRequest(issue)
			} else if resp.StatusCode == 404 {
				er("Merge Request does not exist")
			} else {
				er("Could not complete request." + resp.Status)
			}
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrSubscribeCmd)
}
