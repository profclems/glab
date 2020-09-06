package commands

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var issueNoteCreateCmd = &cobra.Command{
	Use:     "note <issue-id>",
	Aliases: []string{"comment"},
	Short:   "Add a comment or note to an issue on Gitlab",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		gitlabClient, repo := git.InitGitlabClient()
		mID := args[0]
		body, err := cmd.Flags().GetString("message")
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		}
		if err != nil {
			return err
		}
		mr, _, err := gitlabClient.Issues.GetIssue(repo, manip.StringToInt(mID))
		if err != nil {
			return err
		}
		if body == "" {
			body = manip.Editor(manip.EditorOptions{
				Label:    "Note Message:",
				Help:     "Enter the note message. ",
				FileName: "ISSUE_NOTE_EDITMSG",
			})
		}
		if body == "" {
			return errors.New("aborted... Note is empty")
		}

		noteInfo, _, err := gitlabClient.Notes.CreateIssueNote(repo, manip.StringToInt(mID), &gitlab.CreateIssueNoteOptions{
			Body: &body,
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s#note_%d\n", mr.WebURL, noteInfo.ID)
		return nil
	},
}

func init() {
	issueNoteCreateCmd.Flags().StringP("message", "m", "", "Comment/Note message")
	issueCmd.AddCommand(issueNoteCreateCmd)
}
