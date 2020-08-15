package commands

import (
	"fmt"
	"glab/internal/git"
	"glab/internal/manip"
	"log"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var issueNoteCreateCmd = &cobra.Command{
	Use:     "note <issue-id>",
	Aliases: []string{"comment"},
	Short:   "Add a comment or note to an issue on Gitlab",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		gitlabClient, repo := git.InitGitlabClient()
		mID := args[0]
		body, err := cmd.Flags().GetString("message")
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		if err != nil {
			er(err)
			return
		}
		mr, _, err := gitlabClient.Issues.GetIssue(repo, manip.StringToInt(mID))
		if err != nil {
			er(err)
			return
		}
		if body == "" {
			body = manip.Editor(manip.EditorOptions{
				Label:    "Note Message:",
				Help:     "Enter the note message. ",
				FileName: "ISSUE_NOTE_EDITMSG",
			})
		}
		if body == "" {
			log.Fatal("Aborted... Note is empty")
		}

		noteInfo, _, err := gitlabClient.Notes.CreateIssueNote(repo, manip.StringToInt(mID), &gitlab.CreateIssueNoteOptions{
			Body: &body,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s#note_%d\n", mr.WebURL, noteInfo.ID)
	},
}

func init() {
	issueNoteCreateCmd.Flags().StringP("message", "m", "", "Comment/Note message")
	issueCmd.AddCommand(issueNoteCreateCmd)
}
