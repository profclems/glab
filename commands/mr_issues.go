package commands

import (
	"github.com/spf13/cobra"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/xanzy/go-gitlab"
)

var mrIssuesCmd = &cobra.Command{
	Use:     "issues <id>",
	Short:   `Get issues related to a particular merge request.`,
	Long:    ``,
	Aliases: []string{"issue"},
	Args:    cobra.ExactArgs(1),
	Example: "$ glab mr issues 46",
	Run:     issuesRelatedMergeRequest,
}

func issuesRelatedMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		l := &gitlab.GetIssuesClosedOnMergeOptions{}
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		mr, _, err := gitlabClient.MergeRequests.GetIssuesClosedOnMerge(repo, manip.StringToInt(mergeID), l)
		if err != nil {
			er(err)
		}
		displayAllIssues(mr)
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrCmd.AddCommand(mrIssuesCmd)
}
