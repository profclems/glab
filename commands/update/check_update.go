package update

import (
	"errors"
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/spf13/cobra"
)

func NewCheckUpdateCmd(s *iostreams.IOStreams, version string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "check-update",
		Short:   "Check for latest glab releases",
		Long:    ``,
		Aliases: []string{"update"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return CheckUpdate(s, version, false)
		},
	}

	return cmd
}

func CheckUpdate(s *iostreams.IOStreams, version string, silentErr bool) error {
	latestRelease, err := GetUpdateInfo()
	c := s.Color()
	if err != nil {
		if silentErr {
			return nil
		}
		return errors.New("could not check for update! Make sure you have a stable internet connection")
	}

	if isOlderVersion(latestRelease.Version, version) {
		fmt.Fprintf(s.StdOut, "%s %s → %s\n%s\n",
			c.Yellow("A new version of glab has been released:"),
			c.Red(version), c.Green(latestRelease.Version),
			latestRelease.URL)
	} else {
		if silentErr {
			return nil
		}
		fmt.Fprintf(s.StdOut, "%v %v", c.GreenCheck(),
			c.Green("You are already using the latest version of glab"))
	}
	return nil
}
