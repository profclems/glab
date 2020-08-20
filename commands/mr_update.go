package commands

import (
	"github.com/MakeNowJust/heredoc"
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
	Example: heredoc.Doc(`
	$ glab mr update 23 --ready
	$ glab mr update 23 --draft
	`),
	Args:  cobra.ExactArgs(1),
	RunE:   updateMergeRequest,
}

func updateMergeRequest(cmd *cobra.Command, args []string) error {
	mergeID := manip.StringToInt(args[0])
	l := &gitlab.UpdateMergeRequestOptions{}
	var mergeTitle string
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo = r
	}
	isDraft, _ := cmd.Flags().GetBool("draft")
	isWIP, _ := cmd.Flags().GetBool("wip")
	if m, _ := cmd.Flags().GetString("title"); m != "" {
		mergeTitle = m
	}
	if mergeTitle == "" {
		opts := &gitlab.GetMergeRequestsOptions{}
		mr, _, err := gitlabClient.MergeRequests.GetMergeRequest(repo, mergeID, opts)
		if err != nil {
			return err
		}
		mergeTitle = mr.Title
	}
	if isDraft || isWIP {
		if isDraft {
			mergeTitle = "Draft: " + mergeTitle
		} else {
			mergeTitle = "WIP: " + mergeTitle
		}
	} else if isReady, _ := cmd.Flags().GetBool("ready"); isReady {
		mergeTitle = strings.TrimPrefix(mergeTitle, "Draft:")
		mergeTitle = strings.TrimPrefix(mergeTitle, "draft:")
		mergeTitle = strings.TrimPrefix(mergeTitle, "DRAFT:")
		mergeTitle = strings.TrimPrefix(mergeTitle, "WIP:")
		mergeTitle = strings.TrimPrefix(mergeTitle, "wip:")
		mergeTitle = strings.TrimPrefix(mergeTitle, "Wip:")
		mergeTitle = strings.TrimSpace(mergeTitle)
	}
	l.Title = gitlab.String(mergeTitle)
	if m, _ := cmd.Flags().GetBool("lock-discussion"); m {
		l.DiscussionLocked = gitlab.Bool(m)
	}
	if m, _ := cmd.Flags().GetString("description"); m != "" {
		l.Description = gitlab.String(m)
	}
	mr, _, err := gitlabClient.MergeRequests.UpdateMergeRequest(repo, mergeID, l)
	if err != nil {
		return err
	}
	displayMergeRequest(mr)
	return nil
}

func init() {
	mrUpdateCmd.Flags().BoolP("draft", "", false, "Mark merge request as a draft")
	mrUpdateCmd.Flags().BoolP("ready", "r", false, "Mark merge request as ready to be reviewed and merged")
	mrUpdateCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrUpdateCmd.Flags().StringP("title", "t", "", "Title of merge request")
	mrUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on merge request")
	mrUpdateCmd.Flags().StringP("description", "d", "", "merge request description")
	mrCmd.AddCommand(mrUpdateCmd)
}
