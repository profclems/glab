package version

import (
	"fmt"

	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

func NewCmdVersion(version, build string) *cobra.Command {
	versionOutput := fmt.Sprintf("glab %s (%s)", version, build)
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "show glab version information",
		Long:    ``,
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(utils.ColorableOut(cmd), versionOutput)
			return nil
		},
	}
	versionCmd.Root().SetVersionTemplate(versionOutput)
	return versionCmd
}
