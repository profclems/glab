package unsubscribe

import (
	"fmt"
	mr2 "github.com/profclems/glab/commands/mr"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var mrUnsubscribeCmd = &cobra.Command{
	Use:     "unsubscribe <id>",
	Short:   `Unsubscribe to merge requests`,
	Long:    ``,
	Aliases: []string{"unsub"},
	Args:    cobra.ExactArgs(1),
	Run:     unsubscribeMergeRequest,
}

func unsubscribeMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Unsubscribing Merge Request #" + i2)
			mr, resp, _ := gitlabClient.MergeRequests.UnsubscribeFromMergeRequest(repo, manip.StringToInt(i2))

			if resp.StatusCode == 204 {
				bodyString := resp.Body
				fmt.Println(bodyString)
				fmt.Println(color.Green.Sprint("You have successfully unsubscribed to merge request #" + i2))
				mr2.displayMergeRequest(mr)
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
	mr2.mrCmd.AddCommand(mrUnsubscribeCmd)
}
