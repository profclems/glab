package version

import (
	"fmt"

	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

var VersionOutput = "DEV"

func NewCmdVersion(s *utils.IOStreams, version, build string) *cobra.Command {
	VersionOutput = fmt.Sprintf("glab %s (%s)", version, build)
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "show glab version information",
		Long:    ``,
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(s.StdOut, VersionOutput)
			return nil
		},
	}
	versionCmd.Root().SetVersionTemplate(VersionOutput)
	return versionCmd
}
