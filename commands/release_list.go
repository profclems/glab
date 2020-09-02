package commands

import (
	"github.com/profclems/glab/internal/git"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var releaseListCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   `List releases`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.MaximumNArgs(3),
	RunE:    listReleases,
}

func listReleases(cmd *cobra.Command, args []string) error {

	l := &gitlab.ListReleasesOptions{}

	tag, err := cmd.Flags().GetString("tag")

	if err != nil {
		return err
	}

	gitlabClient, repo := git.InitGitlabClient()

	if tag != "" {
		release, _, err := gitlabClient.Releases.GetRelease(repo, tag)
		if err != nil {
			return err
		}
		displayRelease(release)
	} else {
		releases, _, err := gitlabClient.Releases.ListReleases(repo, l)
		if err != nil {
			return err
		}
		displayAllReleases(releases)
	}
	return nil
}

func init() {
	releaseListCmd.Flags().StringP("tag", "t", "", "Filter releases by tag <name>")
	releaseCmd.AddCommand(releaseListCmd)
}
