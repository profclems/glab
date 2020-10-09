package label

import (
	"github.com/profclems/glab/commands/cmdutils"
	labelListCmd "github.com/profclems/glab/commands/label/list"
	"github.com/spf13/cobra"
)

func NewCmdLabel(f *cmdutils.Factory) *cobra.Command {
	var labelCmd = &cobra.Command{
		Use:   "label <command> [flags]",
		Short: `Manage labels on remote`,
		Long:  ``,
	}

	cmdutils.EnableRepoOverride(labelCmd, f)

	labelCmd.AddCommand(labelListCmd.NewCmdList(f))
	return labelCmd
}
