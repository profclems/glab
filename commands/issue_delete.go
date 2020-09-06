package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
)

var issueDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   `Delete an issue`,
	Long:    ``,
	Aliases: []string{"del"},
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
		return nil
	},
}

func init() {
	issueCmd.AddCommand(issueDeleteCmd)
}
