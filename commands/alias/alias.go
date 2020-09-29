package alias

import (
	deleteCmd "github.com/profclems/glab/commands/alias/delete"
	listCmd "github.com/profclems/glab/commands/alias/list"
	setCmd "github.com/profclems/glab/commands/alias/set"
	"github.com/profclems/glab/commands/cmdutils"

	"github.com/spf13/cobra"
)

func NewCmdAlias(f *cmdutils.Factory) *cobra.Command {
	var aliasCmd = &cobra.Command{
		Use:   "alias [command] [flags]",
		Short: `Create, list and delete aliases`,
		Long:  ``,
	}
	aliasCmd.AddCommand(deleteCmd.NewCmdDelete(f, nil))
	aliasCmd.AddCommand(listCmd.NewCmdList(f, nil))
	aliasCmd.AddCommand(setCmd.NewCmdSet(f, nil))
	return aliasCmd
}
