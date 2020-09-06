package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
)

var issueSubscribeCmd = &cobra.Command{
	Use:     "subscribe <id>",
	Short:   `Subscribe to an issue`,
	Long:    ``,
	Aliases: []string{"sub"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return nil
		}
		if len(args) > 0 {
			mergeID := strings.TrimSpace(args[0])
			gitlabClient, repo := git.InitGitlabClient()
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				repo, _ = fixRepoNamespace(r)
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
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueSubscribeCmd)
}
