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
	releaseCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")

	releaseCmd.AddCommand(releaseListCmd.NewCmdReleaseList(f))
	return releaseCmd
}
