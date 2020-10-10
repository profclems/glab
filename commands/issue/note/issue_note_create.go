package note

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	gitlab "github.com/xanzy/go-gitlab"
)

func NewCmdNote(f *cmdutils.Factory) *cobra.Command {
	var issueNoteCreateCmd = &cobra.Command{
		Use:     "note <issue-id>",
		Aliases: []string{"comment"},
		Short:   "Add a comment or note to an issue on Gitlab",
		Long:    ``,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mID := args[0]

			body, err := cmd.Flags().GetString("message")
			if err != nil {
				return err
			}

			mr, err := api.GetIssue(apiClient, repo.FullName(), utils.StringToInt(mID))
			if err != nil {
				return err
			}

			if body == "" {
				body = utils.Editor(utils.EditorOptions{
					Label:    "Note Message:",
					Help:     "Enter the note message. ",
					FileName: "ISSUE_NOTE_EDITMSG",
				})
			}

			if body == "" {
				return errors.New("aborted... Note is empty")
			}

			noteInfo, err := api.CreateIssueNote(apiClient, repo.FullName(), utils.StringToInt(mID), &gitlab.CreateIssueNoteOptions{
				Body: &body,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s#note_%d\n", mr.WebURL, noteInfo.ID)
			return nil
		},
	}
	issueNoteCreateCmd.Flags().StringP("message", "m", "", "Comment/Note message")

	return issueNoteCreateCmd
}
