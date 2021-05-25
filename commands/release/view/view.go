package view

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/profclems/glab/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type ViewOpts struct {
	TagName       string
	OpenInBrowser bool

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Config     func() (config.Config, error)
}

func NewCmdView(f *cmdutils.Factory, runE func(opts *ViewOpts) error) *cobra.Command {
	opts := &ViewOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "view <tag>",
		Short: "View information about a GitHub Release",
		Long: heredoc.Docf(`View information about a GitHub Release.

			Without an explicit tag name argument, the latest release in the project is shown.
		`, "`"),
		Example: heredoc.Doc(`
			View the latest release of a GitLab repository
			$ glab release view

			View a release with specified tag name
			$ glab release view v1.0.1 
`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			if len(args) == 1 {
				opts.TagName = args[0]
			}

			if runE != nil {
				return runE(opts)
			}

			return deleteRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.OpenInBrowser, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

func deleteRun(opts *ViewOpts) error {
	client, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}

	cfg, _ := opts.Config()

	var resp *gitlab.Response

	if opts.TagName == "" {
		tags, resp, err := client.Tags.ListTags(repo.FullName(), &gitlab.ListTagsOptions{
			OrderBy: gitlab.String("updated"),
			Sort:    gitlab.String("desc"),
		})
		if err != nil && resp != nil && resp.StatusCode != 404 {
			return cmdutils.WrapError(err, "could not fetch tag")
		}
		if len(tags) < 1 {
			return cmdutils.WrapError(err, "no tags found. Create a new tag, push to GitLab and try this command again.")
		}
		opts.TagName = tags[0].Name
	}

	release, resp, err := client.Releases.GetRelease(repo.FullName(), opts.TagName)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
			return cmdutils.WrapError(err, "release does not exist.")
		}
		return cmdutils.WrapError(err, "failed to fetch release")
	}

	if opts.OpenInBrowser { //open in browser if --web flag is specified
		url := fmt.Sprintf("%s://%s/%s/-/releases/%s",
			glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault(),
			repo.FullName(), release.TagName)

		if opts.IO.IsOutputTTY() {
			opts.IO.Logf("Opening %s in your browser.\n", url)
		}

		browser, _ := cfg.Get(repo.RepoHost(), "browser")
		return utils.OpenInBrowser(url, browser)
	}

	glamourStyle, _ := cfg.Get(repo.RepoHost(), "glamour_style")
	opts.IO.ResolveBackgroundColor(glamourStyle)

	err = opts.IO.StartPager()
	if err != nil {
		return err
	}
	defer opts.IO.StopPager()

	opts.IO.LogInfo(releaseutils.DisplayRelease(opts.IO, release))
	return nil
}
