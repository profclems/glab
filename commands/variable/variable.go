package variable

import (
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	setCmd "github.com/profclems/glab/commands/variable/set"
	"github.com/spf13/cobra"
)

func NewVariableCmd(f *cmdutils.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "variable",
		Short:   "Manage GitLab Project and Group Variables",
		Aliases: []string{"var"},
	}

	cmdutils.EnableRepoOverride(cmd, f)
	action.EnableRepoOverride(cmd, f)

	cmd.AddCommand(setCmd.NewCmdSet(f, nil))
	return cmd
}
