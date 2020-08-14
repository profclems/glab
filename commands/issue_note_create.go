package commands


import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"glab/internal/git"
	"glab/internal/manip"
	"log"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

var issueCreateNoteCmd = &cobra.Command{
	Use:     "note <issue-id>",
	Aliases: []string{"comment"},
	Short:   "Add a comment to issue",
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
			prompt := &survey.Editor{
				Renderer:      survey.Renderer{},
				Message:       "Note Message: ",
				Help:          "Enter the note message for issue. Uses the editor defined by the $VISUAL or $EDITOR environment variables). If neither of those are present, notepad (on Windows) or vim (Linux or Mac) is used",
				FileName:      "*.md",
			}
			err = survey.AskOne(prompt, &body)
		}

		if err != nil {
			er(err)
			return
		}
		if body == "" {
			log.Fatal("Aborted... Note has an empty message")
		}

		noteInfo,_, err := gitlabClient.Notes.CreateIssueNote(repo, manip.StringToInt(mID), &gitlab.CreateIssueNoteOptions{
			Body: &body,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s#note_%d\n",mr.WebURL, noteInfo.ID)
	},
}

func init() {
	issueCreateNoteCmd.Flags().StringP("message", "m", "", "Enter note message")
	issueCmd.AddCommand(issueCreateNoteCmd)
}