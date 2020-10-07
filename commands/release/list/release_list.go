package list

import (
	"fmt"

	"github.com/profclems/glab/pkg/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var factory *cmdutils.Factory

func NewCmdReleaseList(f *cmdutils.Factory) *cobra.Command {
	factory = f
	var releaseListCmd = &cobra.Command{
		Use:     "list [flags]",
		Short:   `List releases in a repository`,
		Long:    ``,
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			factory = f
			return listReleases(cmd, args)
		},
	}
	releaseListCmd.Flags().StringP("tag", "t", "", "Filter releases by tag <name>")
	return releaseListCmd
}

func listReleases(cmd *cobra.Command, args []string) error {

	l := &gitlab.ListReleasesOptions{}

	tag, err := cmd.Flags().GetString("tag")

	if err != nil {
		return err
	}
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		factory, err = factory.NewClient(r)
		if err != nil {
			return err
		}
	}
	apiClient, err := factory.HttpClient()
	if err != nil {
		return err
	}
	repo, err := factory.BaseRepo()
	if err != nil {
		return err
	}

	if tag != "" {
		release, err := api.GetRelease(apiClient, repo.FullName(), tag)
		if err != nil {
			return err
		}

		cfg, _ := factory.Config()
		glamourStyle, _ := cfg.Get(repo.RepoHost(), "glamour_style")
		fmt.Fprintln(utils.ColorableOut(cmd), releaseutils.DisplayRelease(release, glamourStyle))
	} else {
		releases, err := api.ListReleases(apiClient, repo.FullName(), l)
		if err != nil {
			return err
		}
		fmt.Fprintln(utils.ColorableOut(cmd), releaseutils.DisplayAllReleases(releases, repo.FullName()))
	}
	return nil
}
