package auth

import (
	authLoginCmd "github.com/profclems/glab/commands/auth/login"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Manage glab's authentication state",
	}

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f))

	return cmd
}
