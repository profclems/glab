package label

import (
	"github.com/profclems/glab/commands/cmdutils"
	labelCreateCmd "github.com/profclems/glab/commands/label/create"
	labelListCmd "github.com/profclems/glab/commands/label/list"
	"github.com/spf13/cobra"
)

func NewCmdLabel(f *cmdutils.Factory) *cobra.Command {
	var labelCmd = &cobra.Command{
		Use:   "label <command> [flags]",
		Short: `Manage labels on remote`,
		Long:  ``,
	}
	labelCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format. Supports group namespaces")

	labelCmd.AddCommand(labelListCmd.NewCmdList(f))
	labelCmd.AddCommand(labelCreateCmd.NewCmdCreate(f))
	return labelCmd
}
