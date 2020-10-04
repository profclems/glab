package board

import (
	"github.com/profclems/glab/commands/cmdutils"
	boardCreateCmd "github.com/profclems/glab/commands/issue/board/create"
	"github.com/spf13/cobra"
)

func NewCmdBoard(f *cmdutils.Factory) *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "board [command] [flags]",
		Short: `Work with GitLab Issue Boards in the given project.`,
		Long:  ``,
	}

	issueCmd.AddCommand(boardCreateCmd.NewCmdCreate(f))
	issueCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")

	return issueCmd
}
