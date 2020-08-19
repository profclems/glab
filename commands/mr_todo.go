package commands

import (
	"fmt"
	"glab/internal/git"
	"glab/internal/manip"

	"github.com/spf13/cobra"
)

var mrToDoCmd = &cobra.Command{
	Use:     "todo <merge-request-id>",
	Aliases: []string{"add-todo"},
	Short:   "Add a ToDo to merge request",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		gitlabClient, repo := git.InitGitlabClient()
		mID := args[0]

		_, _, err := gitlabClient.MergeRequests.CreateTodo(repo, manip.StringToInt(mID))
		if err != nil {
			return err
		}
		fmt.Println("Done!!")
		return nil
	},
}

func init() {
	mrCmd.AddCommand(mrToDoCmd)
}
