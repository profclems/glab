package variable

import (
	"github.com/profclems/glab/commands/cmdutils"
	deleteCmd "github.com/profclems/glab/commands/variable/delete"
	listCmd "github.com/profclems/glab/commands/variable/list"
	setCmd "github.com/profclems/glab/commands/variable/set"
	updateCmd "github.com/profclems/glab/commands/variable/update"
	"github.com/spf13/cobra"
)

func NewVariableCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variable",
		Short:   "Manage GitLab Project and Group Variables",
		Aliases: []string{"var"},
	}

	cmdutils.EnableRepoOverride(cmd, f)

	cmd.AddCommand(setCmd.NewCmdSet(f, nil))
	cmd.AddCommand(listCmd.NewCmdSet(f, nil))
	cmd.AddCommand(deleteCmd.NewCmdSet(f, nil))
	cmd.AddCommand(updateCmd.NewCmdSet(f, nil))
	return cmd
}
