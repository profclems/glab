package note

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/commands/issue/issueutils"

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
		Short:   "Add a comment or note to an issue on GitLab",
		Long:    ``,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := f.IO.StdOut

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			issue, repo, err := issueutils.IssueFromArg(apiClient, f.BaseRepo, args[0])
			if err != nil {
				return err
			}

			body, _ := cmd.Flags().GetString("message")

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

			noteInfo, err := api.CreateIssueNote(apiClient, repo.FullName(), issue.IID, &gitlab.CreateIssueNoteOptions{
				Body: &body,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s#note_%d\n", issue.WebURL, noteInfo.ID)
			return nil
		},
	}
	issueNoteCreateCmd.Flags().StringP("message", "m", "", "Comment/Note message")

	return issueNoteCreateCmd
}
