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

var mrCloseCmd = &cobra.Command{
	Use:     "close <id>",
	Short:   `Close merge requests`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(1),
	Run:     closeMergeRequestState,
}

func closeMergeRequestState(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		l := &gitlab.UpdateMergeRequestOptions{}
		l.StateEvent = gitlab.String("close")
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Closing Merge request...")
			mr, resp, _ := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(i2), l)
			if resp.StatusCode == 200 {
				fmt.Println(aurora.Green("You have closed merge request #" + i2))
				displayMergeRequest(mr)
			} else if resp.StatusCode == 404 {
				fmt.Println(aurora.Red("MergeRequest does not exist"))
			} else {
				fmt.Println("Could not complete request: ", resp.Status)
			}
		}
	} else {
		cmd.Usage()
	}
}

func init() {
	mrCmd.AddCommand(mrCloseCmd)
}
