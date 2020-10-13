package note

import (
	"fmt"

	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdNote(f *cmdutils.Factory) *cobra.Command {
	var mrCreateNoteCmd = &cobra.Command{
		Use:     "note <merge-request-id>",
		Aliases: []string{"comment"},
		Short:   "Add a comment or note to merge request",
		Long:    ``,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			body, err := cmd.Flags().GetString("message")
			if err != nil {
				return err
			}

			mr, err = api.GetMR(apiClient, repo.FullName(), mr.IID, &gitlab.GetMergeRequestsOptions{})
			if err != nil {
				return err
			}
			if body == "" {
				body = utils.Editor(utils.EditorOptions{
					Label:    "Note Message:",
					Help:     "Enter the note message for the merge request. ",
					FileName: "*_MR_NOTE_EDITMSG.md",
				})
			}
			if body == "" {
				return fmt.Errorf("aborted... Note has an empty message")
			}

			noteInfo, err := api.CreateMRNote(apiClient, repo.FullName(), mr.IID, &gitlab.CreateMergeRequestNoteOptions{
				Body: &body,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s#note_%d\n", mr.WebURL, noteInfo.ID)
			return nil
		},
	}

	mrCreateNoteCmd.Flags().StringP("message", "m", "", "Comment/Note message")
	return mrCreateNoteCmd
}
