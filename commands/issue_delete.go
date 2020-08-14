package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var issueDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   `Delete an issue`,
	Long:    ``,
	Aliases: []string{"del"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return
		}
		if len(args) > 0 {
			issueID := strings.TrimSpace(args[0])
			gitlabClient, repo := git.InitGitlabClient()
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				repo = r
			}
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("Deleting Issue #" + i2)
				issue, _ := gitlabClient.Issues.DeleteIssue(repo, manip.StringToInt(i2))
				if isSuccessful(issue.StatusCode) {
					fmt.Println("Issue Deleted")
				} else if issue.StatusCode == 404 {
					er("Issue does not exist")
				} else {
					er(issue.Status)
				}
			}
		} else {
			cmdErr(cmd, args)
		}
	},
}

func init() {
	issueCmd.AddCommand(issueDeleteCmd)
}
