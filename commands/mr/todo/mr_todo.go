package todo

import (
	"fmt"

	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdTodo(f *cmdutils.Factory) *cobra.Command {
	var mrToDoCmd = &cobra.Command{
		Use:     "todo <merge-request-id>",
		Aliases: []string{"add-todo"},
		Short:   "Add a ToDo to merge request",
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

			_, err = api.MRTodo(apiClient, repo.FullName(), mr.IID, nil)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, utils.GreenCheck(), "Done!!")

			return nil
		},
	}

	return mrToDoCmd
}
