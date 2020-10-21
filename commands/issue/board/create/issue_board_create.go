package create

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var boardName string

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:     "create [flags]",
		Short:   `Create a project issue board.`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				boardName = args[0]
			}
			var err error
			out := f.IO.StdOut

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			if boardName == "" {
				boardName = utils.AskQuestionWithInput("Board Name:", "", true)
			}

			opts := &gitlab.CreateIssueBoardOptions{
				Name: gitlab.String(boardName),
			}

			fmt.Fprintln(out, "- Creating board")

			issueBoard, err := api.CreateIssueBoard(apiClient, repo.FullName(), opts)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s Board created: %q", utils.GreenCheck(), issueBoard.Name)

			return nil
		},
	}

	issueCmd.Flags().StringVarP(&boardName, "name", "n", "", "The name of the new board")

	return issueCmd
}
