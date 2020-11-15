package update

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

func NewCheckUpdateCmd(s *utils.IOStreams, version, build string) *cobra.Command {
	// versionCmd represents the version command
	var cmd = &cobra.Command{
		Use:     "check-update",
		Short:   "Check for latest glab releases",
		Long:    ``,
		Aliases: []string{"update"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return CheckUpdate(cmd, s, version, build, false)
		},
	}

	return cmd
}

func CheckUpdate(cmd *cobra.Command, s *utils.IOStreams, version, build string, silentErr bool) error {
	releaseInfo, err := GetUpdateInfo()
	if err != nil {
		if silentErr {
			return nil
		}
		return errors.New("could not check for update! Make sure you have a stable internet connection")
	}

	latestVersion := strings.TrimSpace(releaseInfo.Name)
	version = strings.TrimSpace(version)

	if isLatestVersion(latestVersion, version) {
		if silentErr {
			return nil
		}
		fmt.Fprintf(s.StdOut, "%v %v", utils.GreenCheck(),
			utils.Green("You are already using the latest version of glab"))
	} else {
		fmt.Fprintf(s.StdOut, "%s %s → %s\n%s\n",
			utils.Yellow("A new version of glab has been released:"),
			utils.Red(version), utils.Green(latestVersion),
			releaseInfo.HTMLUrl)
	}
	return nil
}

func isLatestVersion(latestVersion, appVersion string) bool {
	latestVersion = strings.TrimSpace(latestVersion)
	appVersion = strings.TrimSpace(appVersion)
	vo, v1e := version.NewVersion(appVersion)
	vn, v2e := version.NewVersion(latestVersion)
	return v1e == nil && v2e == nil && vo.LessThan(vn)
}
