package user

import (
	"github.com/profclems/glab/commands/cmdutils"
	userEventsCmd "github.com/profclems/glab/commands/user/events"
	"github.com/spf13/cobra"
)

func NewCmdUser(f *cmdutils.Factory) *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "user <command> [flags]",
		Short: "Interact with user",
		Long:  "",
	}

	userCmd.AddCommand(userEventsCmd.NewCmdEvents(f))

	return userCmd
}
