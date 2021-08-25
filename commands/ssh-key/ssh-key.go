package ssh

import (
	"github.com/profclems/glab/commands/cmdutils"
	cmdAdd "github.com/profclems/glab/commands/ssh-key/add"
	cmdGet "github.com/profclems/glab/commands/ssh-key/get"
	cmdList "github.com/profclems/glab/commands/ssh-key/list"
	"github.com/spf13/cobra"
)

func NewCmdSSHKey(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key <command>",
		Short: "Manage SSH keys",
		Long:  "Manage SSH keys registered with your GitLab account",
	}
	
	cmdutils.EnableRepoOverride(cmd, f)

	cmd.AddCommand(cmdAdd.NewCmdAdd(f, nil))
	cmd.AddCommand(cmdGet.NewCmdGet(f, nil))
	cmd.AddCommand(cmdList.NewCmdList(f, nil))

	return cmd
}
