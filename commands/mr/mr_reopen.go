package mr

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var mrReopenCmd = &cobra.Command{
	Use:     "reopen <id>",
	Short:   `Reopen merge requests`,
	Long:    ``,
	Aliases: []string{"open"},
	Args:    cobra.ExactArgs(1),
	Run:     reopenMergeRequestState,
}

func reopenMergeRequestState(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		l := &gitlab.UpdateMergeRequestOptions{}
		l.StateEvent = gitlab.String("reopen")
		arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Printf("Updating Merge request #%s...\n", i2)
			mr, resp, _ := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(i2), l)
			if resp.StatusCode == 200 {
				fmt.Println(color.Green.Sprint("You have reopened merge request #" + i2))
				displayMergeRequest(mr)
			} else if resp.StatusCode == 404 {
				er("MergeRequest does not exist")
			} else {
				er("Could not complete request: " + resp.Status)
			}
		}
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrReopenCmd)
}
