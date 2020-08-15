package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var issueReopenCmd = &cobra.Command{
	Use:     "reopen <id>",
	Short:   `Reopen a closed issue`,
	Long:    ``,
	Aliases: []string{"open"},
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
				repo = r
			}
			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("reopen")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("Reopening Issue...")
				issue, resp, err := gitlabClient.Issues.UpdateIssue(repo, manip.StringToInt(i2), l)
				if err != nil {
					return err
				}
				if isSuccessful(resp.StatusCode) {
					fmt.Println("Issue #" + i2 + " eopened")
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
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueReopenCmd)
}
