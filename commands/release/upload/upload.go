package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils"
	"github.com/profclems/glab/commands/release/releaseutils/upload"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type UploadOpts struct {
	TagName          string
	AssetLinksAsJson string

	AssetLinks []*upload.ReleaseAsset
	AssetFiles []*upload.ReleaseFile

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Config     func() (config.Config, error)
}

func NewCmdUpload(f *cmdutils.Factory, runE func(opts *UploadOpts) error) *cobra.Command {
	opts := &UploadOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "upload <tag> [<files>...]",
		Short: "Upload release asset files or links to GitLab Release",
		Long: heredoc.Doc(`Upload release assets to GitLab Release

				You can define the display name by appending '#' after the file name. 
				The link type comes after the display name (eg. 'myfile.tar.gz#My display name#package')
		`),
		Args: func() cobra.PositionalArgs {
			return func(cmd *cobra.Command, args []string) error {
				if len(args) < 1 {
					return &cmdutils.FlagError{Err: errors.New("no tag name provided")}
				}
				if len(args) < 2 && opts.AssetLinksAsJson == "" {
					return cmdutils.FlagError{Err: errors.New("no files specified")}
				}
				return nil
			}
		}(),
		Example: heredoc.Doc(`
			Upload a release asset with a display name
			$ glab release upload v1.0.1 '/path/to/asset.zip#My display label'

			Upload a release asset with a display name and type
			$ glab release upload v1.0.1 '/path/to/asset.png#My display label#image'

			Upload all assets in a specified folder
			$ glab release upload v1.0.1 ./dist/*

			Upload all tarballs in a specified folder
			$ glab release upload v1.0.1 ./dist/*.tar.gz

			Upload release assets links specified as JSON string
			$ glab release upload v1.0.1 --assets-links='
				[
					{
						"name": "Asset1", 
						"url":"https://<domain>/some/location/1", 
						"link_type": "other", 
						"filepath": "path/to/file"
					}
				]'
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.TagName = args[0]

			opts.AssetFiles, err = releaseutils.AssetsFromArgs(args[1:])
			if err != nil {
				return err
			}

			if opts.AssetFiles == nil && opts.AssetLinksAsJson == "" {
				return cmdutils.FlagError{Err: errors.New("no files specified")}
			}

			if opts.AssetLinksAsJson != "" {
				err := json.Unmarshal([]byte(opts.AssetLinksAsJson), &opts.AssetLinks)
				if err != nil {
					return fmt.Errorf("failed to parse JSON string: %w", err)
				}
			}

			if runE != nil {
				return runE(opts)
			}

			return deleteRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.AssetLinksAsJson, "assets-links", "a", "", "`JSON` string representation of assets links (e.g. `--assets='[{\"name\": \"Asset1\", \"url\":\"https://<domain>/some/location/1\", \"link_type\": \"other\", \"filepath\": \"path/to/file\"}]')`")

	return cmd
}

func deleteRun(opts *UploadOpts) error {
	start := time.Now()

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

	opts.IO.Logf("%s validating tag %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)

	release, resp, err := client.Releases.GetRelease(repo.FullName(), opts.TagName)
	if err != nil {
		if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
			return cmdutils.WrapError(err, "release does not exist. Create a new release with `glab release create "+opts.TagName+"`")
		}
		return cmdutils.WrapError(err, "failed to fetch release")
	}

	opts.IO.Logf("%s uploading release assets %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)
	// upload files and create asset link
	if opts.AssetFiles != nil || opts.AssetLinks != nil {
		uploadCtx := upload.Context{
			IO:          opts.IO,
			Client:      client,
			AssetsLinks: opts.AssetLinks,
			AssetFiles:  opts.AssetFiles,
		}
		if err = uploadCtx.UploadFiles(repo.FullName(), release.TagName); err != nil {
			return cmdutils.WrapError(err, "upload failed")
		}

		// create asset link for assets provided as json
		if err = uploadCtx.CreateReleaseAssetLinks(repo.FullName(), release.TagName); err != nil {
			return cmdutils.WrapError(err, "failed to create release link")
		}
	}

	opts.IO.Logf(color.Bold("%s upload succeeded after %0.2fs\n"), color.GreenCheck(), time.Since(start).Seconds())
	return nil
}
