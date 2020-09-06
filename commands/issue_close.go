package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var issueCloseCmd = &cobra.Command{
	Use:     "close <id>",
	Short:   `Close an issue`,
	Long:    ``,
	Aliases: []string{"unsub"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return nil
		}
		if len(args) > 0 {
			issueID := strings.TrimSpace(args[0])
			gitlabClient, repo := git.InitGitlabClient()
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				repo, _ = fixRepoNamespace(r)
			}
			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("Closing Issue...")
				issue, resp, err := gitlabClient.Issues.UpdateIssue(repo, manip.StringToInt(i2), l)
				if err != nil {
					return err
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
			return nil
		}
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueCloseCmd)
}
