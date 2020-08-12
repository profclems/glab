package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var issueUnsubscribeCmd = &cobra.Command{
	Use:     "unsubscribe <id>",
	Short:   `Unsubscribe to an issue`,
	Long:    ``,
	Aliases: []string{"unsub"},
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
				fmt.Println("Unsubscribing to Issue #" + i2)
				issue, resp, _ := gitlabClient.Issues.UnsubscribeFromIssue(repo, manip.StringToInt(i2))

				if isSuccessful(resp.StatusCode) {
					fmt.Println("Unsubscribed to issue #" + i2)
					displayIssue(issue)
				} else if resp.StatusCode == 404 {
					er("Issue does not exist")
				} else {
					er(resp.Status)
				}
			}
		} else {
			cmdErr(cmd, args)
		}
	},
}

func init() {
	issueCmd.AddCommand(issueUnsubscribeCmd)
}
