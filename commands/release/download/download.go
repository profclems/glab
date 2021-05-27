package download

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils/upload"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type DownloadOpts struct {
	TagName    string
	Asset      string
	AssetNames []string
	Dir        string
	External   bool

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Config     func() (config.Config, error)
}

func NewCmdDownload(f *cmdutils.Factory, runE func(opts *DownloadOpts) error) *cobra.Command {
	opts := &DownloadOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "download <tag>",
		Short: "Download asset files from a GitLab Release",
		Long: heredoc.Docf(`Download asset files from a GitLab Release

			If no tag is specified, assets are downloaded from the latest release.
			Use %[1]s--asset-name%[1]s to specify a file name to download from the release assets.
			%[1]s--asset-name%[1]s flag accepts glob patterns.

			Unless %[1]s--include-external%[1]s flag is specified, external files are not downloaded.
		`, "`"),
		Args: cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			Download all assets from the latest release
			$ glab release download

			Download all assets from the specified release tag
			$ glab release download v1.1.0

			Download assets with names matching the glob pattern
			$ glab release download v1.10.1 --asset-name="*.tar.gz"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			if len(args) == 1 {
				opts.TagName = args[0]
			}

			if runE != nil {
				return runE(opts)
			}

			return downloadRun(opts)
		},
	}

	cmd.Flags().StringArrayVarP(&opts.AssetNames, "asset-name", "n", []string{}, "Download only assets that match the name or a glob pattern")
	cmd.Flags().StringVarP(&opts.Dir, "dir", "D", ".", "Directory to download the release assets to")
	cmd.Flags().BoolVarP(&opts.External, "include-external", "x", false, "Include external asset files")

	return cmd
}

func downloadRun(opts *DownloadOpts) error {
	client, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}
	color := opts.IO.Color()
	var resp *gitlab.Response
	var release *gitlab.Release
	var downloadableAssets []*upload.ReleaseAsset

	if opts.TagName == "" {
		opts.IO.Logf("%s fetching latest release %s=%s\n",
			color.ProgressIcon(),
			color.Blue("repo"), repo.FullName())
		releases, _, err := client.Releases.ListReleases(repo.FullName(), &gitlab.ListReleasesOptions{})
		if err != nil {
			return cmdutils.WrapError(err, "could not fetch latest release")
		}
		if len(releases) < 1 {
			return cmdutils.WrapError(errors.New("not found"), fmt.Sprintf("no release found for %q", repo.FullName()))
		}

		release = releases[0]
		opts.TagName = release.TagName
	} else {
		opts.IO.Logf("%s fetching release %s=%s %s=%s\n",
			color.ProgressIcon(),
			color.Blue("repo"), repo.FullName(),
			color.Blue("tag"), opts.TagName)

		release, resp, err = client.Releases.GetRelease(repo.FullName(), opts.TagName)
		if err != nil {
			if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
				return cmdutils.WrapError(err, "release does not exist.")
			}
			return cmdutils.WrapError(err, "failed to fetch release")
		}
	}

	for _, link := range release.Assets.Links {
		if link.External && !opts.External {
			continue
		}
		if len(opts.AssetNames) > 0 && (!matchAny(opts.AssetNames, link.Name)) {
			continue
		}
		downloadableAssets = append(downloadableAssets, &upload.ReleaseAsset{
			Name: &link.Name,
			URL:  &link.URL,
		})
	}

	for _, source := range release.Assets.Sources {
		name := path.Base(source.URL)
		if len(opts.AssetNames) > 0 && (!matchAny(opts.AssetNames, name)) {
			continue
		}
		downloadableAssets = append(downloadableAssets, &upload.ReleaseAsset{
			Name: &name,
			URL:  &source.URL,
		})
	}

	if downloadableAssets == nil || len(downloadableAssets) < 1 {
		opts.IO.Logf("%s no release assets found\n",
			color.DotWarnIcon())
		return nil
	}
	opts.IO.Logf("%s downloading release assets %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)

	err = downloadAssets(api.GetClient(), opts.IO, downloadableAssets, opts.Dir)
	if err != nil {
		return cmdutils.WrapError(err, "failed to download release")
	}

	opts.IO.Logf(color.Bold("%s release %q downloaded\n"), color.RedCheck(), release.Name)

	return nil
}

func matchAny(patterns []string, name string) bool {
	for _, p := range patterns {
		matched, err := filepath.Match(p, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func downloadAssets(httpClient *api.Client, io *iostreams.IOStreams, toDownload []*upload.ReleaseAsset, destDir string) error {
	color := io.Color()
	for _, asset := range toDownload {
		io.Logf("%s downloading file %s=%s %s=%s\n",
			color.ProgressIcon(),
			color.Blue("name"), *asset.Name,
			color.Blue("url"), *asset.URL)
		err := downloadAsset(httpClient, *asset.URL, filepath.Join(destDir, *asset.Name))
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadAsset(client *api.Client, assetURL, destinationPath string) error {
	var body io.Reader
	//color := streams.Color()

	baseURL, _ := url.Parse(assetURL)

	req, err := api.NewHTTPRequest(client, "GET", baseURL, body, []string{"Accept:application/octet-stream"}, false)
	if err != nil {
		return err
	}
	resp, err := client.HTTPClient().Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return errors.New(resp.Status)
	}

	f, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
