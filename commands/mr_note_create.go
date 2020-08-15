package commands

import (
	"fmt"
	"glab/internal/git"
	"glab/internal/manip"
	"log"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var mrCreateNoteCmd = &cobra.Command{
	Use:     "note <merge-request-id>",
	Aliases: []string{"comment"},
	Short:   "Add a comment or note to merge request",
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
		mr, _, err := gitlabClient.MergeRequests.GetMergeRequest(repo, manip.StringToInt(mID), &gitlab.GetMergeRequestsOptions{})
		if err != nil {
			er(err)
			return
		}
		if body == "" {
			body = manip.Editor(manip.EditorOptions{
				Label: "Note Message:",
				Help : "Enter the note message for the merge request. ",
				FileName : "*_MR_NOTE_EDITMSG.md",
			})
		}
		if body == "" {
			log.Fatal("Aborted... Note has an empty message")
		}

		noteInfo,_, err := gitlabClient.Notes.CreateMergeRequestNote(repo, manip.StringToInt(mID), &gitlab.CreateMergeRequestNoteOptions{
			Body: &body,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s#note_%d\n",mr.WebURL, noteInfo.ID)
	},
}

func init() {
	mrCreateNoteCmd.Flags().StringP("message", "m", "", "Comment/Note message")
	mrCmd.AddCommand(mrCreateNoteCmd)
}