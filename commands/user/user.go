package user

import (
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
)

func NewCmdUser(f *cmdutils.Factory) *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "user <command> [flags]",
		Short: "Interact with user",
		Long:  "",
	}

	//TODO @zemzale 20/10/20 Add command for user events
	return userCmd
}
