package list

import (
	"fmt"

	"github.com/profclems/glab/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils"
	"github.com/profclems/glab/pkg/utils"

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
	// deprecate in favour of the `release view` command
	_ = releaseListCmd.Flags().MarkDeprecated("tag", "use `glab release view <tag>` instead")

	// make it hidden but still accessible
	// TODO: completely remove before a major release (v2.0.0+)
	_ = releaseListCmd.Flags().MarkHidden("tag")

	return releaseListCmd
}

func listReleases(cmd *cobra.Command, args []string) error {
	l := &gitlab.ListReleasesOptions{}

	tag, err := cmd.Flags().GetString("tag")

	if err != nil {
		return err
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
		factory.IO.ResolveBackgroundColor(glamourStyle)

		err = factory.IO.StartPager()
		if err != nil {
			return err
		}
		defer factory.IO.StopPager()

		fmt.Fprintln(factory.IO.StdOut, releaseutils.DisplayRelease(factory.IO, release))
	} else {
		l.PerPage = 30

		releases, err := api.ListReleases(apiClient, repo.FullName(), l)
		if err != nil {
			return err
		}

		title := utils.NewListTitle("release")
		title.RepoName = repo.FullName()
		title.Page = 0
		title.CurrentPageTotal = len(releases)
		err = factory.IO.StartPager()
		if err != nil {
			return err
		}
		defer factory.IO.StopPager()

		fmt.Fprintf(factory.IO.StdOut, "%s\n%s\n", title.Describe(), releaseutils.DisplayAllReleases(factory.IO, releases, repo.FullName()))
	}
	return nil
}
