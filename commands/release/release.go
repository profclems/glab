package release

import (
	"github.com/profclems/glab/commands/cmdutils"
	releaseListCmd "github.com/profclems/glab/commands/release/list"

	"github.com/spf13/cobra"
)

func NewCmdRelease(f *cmdutils.Factory) *cobra.Command {
	var releaseCmd = &cobra.Command{
		Use:   "release <command> [flags]",
		Short: `Manage GitLab releases`,
		Long:  ``,
	}

	cmdutils.EnableRepoOverride(releaseCmd, f)
	releaseCmd.AddCommand(releaseListCmd.NewCmdReleaseList(f))

	return releaseCmd
}
