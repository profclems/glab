package todo

import (
	"fmt"

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
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mID := args[0]

			_, err = api.MRTodo(apiClient, repo.FullName(), utils.StringToInt(mID), nil)
			if err != nil {
				return err
			}
			fmt.Fprintln(out, utils.GreenCheck(), "Done!!")
			return nil
		},
	}

	return mrToDoCmd
}
