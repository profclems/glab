package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/release/releaseutils/upload"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	ForceDelete bool
	DeleteTag   bool
	TagName     string

	AssetLinks []*upload.ReleaseAsset
	AssetFiles []*upload.ReleaseFile

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Config     func() (config.Config, error)
}

func NewCmdDelete(f *cmdutils.Factory, runE func(opts *CreateOpts) error) *cobra.Command {
	opts := &CreateOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "delete <tag>",
		Short: "Delete a  GitLab Release",
		Long: heredoc.Docf(`Delete release assets to GitLab Release

			Deleting a release does not delete the associated tag unless %[1]s--with-tag%[1]s is specified.
			Maintainer level access to the project is required to delete a release.
		`, "`"),
		Args: cmdutils.MinimumArgs(1, "no tag name provided"),
		Example: heredoc.Doc(`
			Delete a release (with a confirmation prompt')
			$ glab release delete v1.1.0'

			Skip the confirmation prompt and force delete
			$ glab release delete v1.0.1 -y

			Delete release and associated tag
			$ glab release delete v1.0.1 --with-tag
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.TagName = args[0]

			if !opts.ForceDelete && !opts.IO.PromptEnabled() {
				return &cmdutils.FlagError{Err: fmt.Errorf("--yes or -y flag is required when not running interactively")}
			}

			if runE != nil {
				return runE(opts)
			}

			return deleteRun(opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.ForceDelete, "yes", "y", false, "Skip confirmation prompt")
	cmd.Flags().BoolVarP(&opts.DeleteTag, "with-tag", "t", false, "Delete associated tag")

	return cmd
}

func deleteRun(opts *CreateOpts) error {
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
			return cmdutils.WrapError(err, "release does not exist.")
		}
		return cmdutils.WrapError(err, "failed to fetch release")
	}

	if !opts.ForceDelete && opts.IO.PromptEnabled() {
		opts.IO.Logf("This action will permanently delete release %q immediately.\n\n", release.TagName)
		err = prompt.Confirm(&opts.ForceDelete, fmt.Sprintf("Are you ABSOLUTELY SURE you wish to delete this release %q?", release.Name), false)
		if err != nil {
			return cmdutils.WrapError(err, "could not prompt")
		}
	}

	if !opts.ForceDelete {
		return cmdutils.CancelError()
	}

	opts.IO.Logf("%s deleting release %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)

	release, _, err = client.Releases.DeleteRelease(repo.FullName(), release.TagName)
	if err != nil {
		return cmdutils.WrapError(err, "failed to delete release")
	}

	opts.IO.Logf(color.Bold("%s release %q deleted\n"), color.RedCheck(), release.Name)

	if opts.DeleteTag {

		opts.IO.Logf("%s deleting associated tag %q\n",
			color.ProgressIcon(), opts.TagName)

		_, err = client.Tags.DeleteTag(repo.FullName(), release.TagName)
		if err != nil {
			return cmdutils.WrapError(err, "failed to delete tag")
		}

		opts.IO.Logf(color.Bold("%s tag %q deleted\n"), color.RedCheck(), release.Name)
	}
	return nil
}
