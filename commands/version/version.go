package version

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/spf13/cobra"
)

func NewCmdVersion(s *iostreams.IOStreams, version, buildDate string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "show glab version information",
		Long:    ``,
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(s.StdOut, Scheme(version, buildDate))
			return nil
		},
	}
	return versionCmd
}

func Scheme(version, buildDate string) string {
	version = strings.TrimPrefix(version, "v")

	if buildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, buildDate)
	}

	return fmt.Sprintf("glab version %s\n", version)
}
