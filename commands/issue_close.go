package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"log"
	"strings"
)

var issueCloseCmd = &cobra.Command{
	Use:     "close",
	Short:   `Close an issue`,
	Long:    ``,
	Aliases: []string{"unsub"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return
		}
		if len(args) > 0 {
			issueID := strings.TrimSpace(args[0])
			gitlabClient, repo := git.InitGitlabClient()
			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("Closing Issue...")
				issue, resp, err := gitlabClient.Issues.UpdateIssue(repo, manip.StringToInt(i2), l)
				if err != nil {
					log.Fatal(err)
				}
				if isSuccessful(resp.StatusCode) {
					fmt.Println("Issue #" + i2 + " closed")
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
	issueCmd.AddCommand(issueCloseCmd)
}