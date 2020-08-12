package commands

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var mrUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: `Update merge requests`,
	Long:  ``,
	Args:    cobra.ExactArgs(1),
	Run:   updateMergeRequest,
}

func updateMergeRequest(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")
		l := &gitlab.UpdateMergeRequestOptions{}
		if m, _ := cmd.Flags().GetString("title"); m != "" {
			l.Title = gitlab.String(m)
		}
		if m, _ := cmd.Flags().GetBool("lock-discussion"); m {
			l.DiscussionLocked = gitlab.Bool(m)
		}
		if m, _ := cmd.Flags().GetString("description"); m != "" {
			l.Description = gitlab.String(m)
		}
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		mr, _, err := gitlabClient.MergeRequests.UpdateMergeRequest(repo, manip.StringToInt(mergeID), l)
		if err != nil {
			er(err)
		}
		displayMergeRequest(mr)
	} else {
		cmdErr(cmd, args)
	}
}

func init() {
	mrUpdateCmd.Flags().StringP("title", "t", "", "Title of merge request")
	mrUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on merge request")
	mrUpdateCmd.Flags().StringP("description", "d", "", "merge request description")
	mrCmd.AddCommand(mrUpdateCmd)
}
