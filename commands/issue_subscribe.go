package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var issueSubscribeCmd = &cobra.Command{
	Use:     "subscribe <id>",
	Short:   `Subscribe to an issue`,
	Long:    ``,
	Aliases: []string{"sub"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return
		}
		if len(args) > 0 {
			mergeID := strings.TrimSpace(args[0])
			gitlabClient, repo := git.InitGitlabClient()
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				repo = r
			}
			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("Subscribing to Issue #" + i2)
				issue, resp, _ := gitlabClient.Issues.SubscribeToIssue(repo, manip.StringToInt(i2), nil)

				if isSuccessful(resp.StatusCode) {
					fmt.Println("Subscribed to issue #" + i2)
					displayIssue(issue)
				} else if resp.StatusCode == 404 {
					er("Issue does not exist")
				} else {
					er("Could not complete request; " + resp.Status)
				}
			}
		} else {
			cmdErr(cmd, args)
		}
	},
}

func init() {
	issueCmd.AddCommand(issueSubscribeCmd)
}
