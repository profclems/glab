package create

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/profclems/glab/commands/release/upload"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type CreateOpts struct {
	Name             string
	Ref              string
	TagName          string
	Notes            string
	NotesFile        string
	Milestone        []string
	AssetLinksAsJson string
	ReleasedAt       string

	AssetLinks []*upload.ReleaseAsset
	AssetFiles []*upload.ReleaseFile

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
}

func NewCmdCreate(f *cmdutils.Factory, runE func(opts *CreateOpts) error) *cobra.Command {
	opts := &CreateOpts{
		IO: f.IO,
	}

	cmd := &cobra.Command{
		Use:   "create <tag> [<files>...]",
		Short: "Create a new GitLab Release for a repository",
		Long: `Create a new GitLab Release for a repository.

You need push access to the repository to create a Release.`,
		Args: cmdutils.MinimumArgs(1, "no tag name provided"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.TagName = args[0]

			opts.AssetFiles, err = AssetsFromArgs(args[1:])
			if err != nil {
				return err
			}

			if opts.AssetLinksAsJson != "" {
				err := json.Unmarshal([]byte(opts.AssetLinksAsJson), &opts.AssetLinks)
				if err != nil {
					return fmt.Errorf("failed to parse JSON string: %w", err)
				}
			}

			if opts.NotesFile != "" {
				if opts.NotesFile == "-" {
					b, err := ioutil.ReadAll(opts.IO.In)
					_ = opts.IO.In.Close()
					if err != nil {
						return err
					}

					opts.Notes = string(b)
				}
			}

			if runE != nil {
				return runE(opts)
			}

			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "The release name or title")
	cmd.Flags().StringVarP(&opts.Ref, "ref", "r", "", "If a tag specified doesn’t exist, the release is created from ref and tagged with the specified tag name. It can be a commit SHA, another tag name, or a branch name.")
	cmd.Flags().StringVarP(&opts.Notes, "notes", "N", "", "The release notes/description. You can use Markdown")
	cmd.Flags().StringVarP(&opts.NotesFile, "notes-file", "F", "", "Read release notes `file`")
	cmd.Flags().StringVarP(&opts.ReleasedAt, "released-at", "D", "", "The `date` when the release is/was ready. Defaults to the current datetime. Expected in ISO 8601 format (2019-03-15T08:00:00Z)")
	cmd.Flags().StringSliceVarP(&opts.Milestone, "milestone", "m", []string{}, "The title of each milestone the release is associated with")
	cmd.Flags().StringVarP(&opts.AssetLinksAsJson, "assets", "a", "", "`JSON` string representation of assets links (e.g. --assets='[{\"name\": \"Asset1\", \"url\":\"https://<domain>/some/location/1\", \"link_type\": \"other\", \"filepath\": \"path/to/file\" }]'")

	return cmd
}

func createRun(opts *CreateOpts) error {
	client, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := opts.BaseRepo()
	if err != nil {
		return err
	}
	color := opts.IO.Color()

	fmt.Fprintf(opts.IO.StdErr, "%s creating or updating release %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)

	release, resp, err := client.Releases.GetRelease(repo.FullName(), opts.TagName)
	if err != nil && (resp == nil || resp.StatusCode != 403) {
		return err
	}

	var releasedAt time.Time

	if opts.ReleasedAt != "" {
		// Parse the releasedAt to the expected format of the API
		// From the API docs "Expected in ISO 8601 format (2019-03-15T08:00:00Z)".
		releasedAt, err = time.Parse(time.RFC3339[:len(opts.ReleasedAt)], opts.ReleasedAt)
		if err != nil {
			return err
		}
	}

	if resp.StatusCode == 403 || release == nil {
		release, _, err = client.Releases.CreateRelease(repo.FullName(), &gitlab.CreateReleaseOptions{
			Name:        &opts.Name,
			Description: &opts.Notes,
			Ref:         &opts.Ref,
			TagName:     &opts.TagName,
			ReleasedAt:  &releasedAt,
			Milestones:  opts.Milestone,
		})

		if err != nil {
			return err
		}
		fmt.Fprintf(opts.IO.StdErr, "%s, release created\t%s=%s\n", color.ProgressIcon(),
			color.Blue("url"), fmt.Sprintf("%s://%s/%s/releases/tag/%s",
				glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault(),
				repo.FullName(), release.TagName))
	} else {
		apiOpts := &gitlab.UpdateReleaseOptions{}
		if opts.Notes != "" {
			apiOpts.Description = &opts.Notes
		}
		if opts.Name != "" {
			apiOpts.Name = &opts.Name
		}

		if opts.ReleasedAt != "" {
			apiOpts.ReleasedAt = &releasedAt
		}

		if len(opts.Milestone) > 0 {
			apiOpts.Milestones = opts.Milestone
		}

		release, _, err = client.Releases.UpdateRelease(repo.FullName(), opts.TagName, apiOpts)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdErr, "%s release updated\t%s=%s\n", color.ProgressIcon(),
			color.Blue("url"), fmt.Sprintf("%s://%s/%s/-/releases/tag/%s",
				glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault(),
				repo.FullName(), release.TagName))
	}

	// upload files and create asset link
	if opts.AssetFiles != nil || opts.AssetLinks != nil {
		fmt.Fprintf(opts.IO.StdErr, "\n%s Uploading release assets\n", color.ProgressIcon())
		uploadCtx := upload.Context{
			IO:          opts.IO,
			Client:      client,
			AssetsLinks: opts.AssetLinks,
			AssetFiles:  opts.AssetFiles,
		}
		if err = uploadCtx.UploadFiles(repo.FullName(), release.TagName); err != nil {
			return err
		}

		// create asset link for assets provided as json
		if err = uploadCtx.CreateReleaseAssetLinks(repo.FullName(), release.TagName); err != nil {
			return err
		}
	}
	if len(opts.Milestone) > 0 {
		// close all associated milestones
		for _, milestone := range opts.Milestone {
			// run loading msg
			opts.IO.StartSpinner("closing milestone %q", milestone)
			// close milestone
			err := closeMilestone(opts, milestone)
			// stop loading
			opts.IO.StopSpinner("")
			if err != nil {
				fmt.Fprintln(opts.IO.StdErr, color.FailedIcon(), err.Error())
			} else {
				fmt.Fprintf(opts.IO.StdErr, "%s closed milestone %q\n", color.GreenCheck(), milestone)
			}
		}
	}
	fmt.Fprintf(opts.IO.StdErr, "%s release succeeded", color.GreenCheck())
	return nil
}

func getMilestoneByTitle(c *CreateOpts, title string) (*gitlab.Milestone, error) {
	opts := &gitlab.ListMilestonesOptions{
		Title: &title,
	}

	client, err := c.HTTPClient()
	if err != nil {
		return nil, err
	}

	repo, err := c.BaseRepo()
	if err != nil {
		return nil, err
	}

	for {
		milestones, resp, err := client.Milestones.ListMilestones(repo.FullName(), opts)
		if err != nil {
			return nil, err
		}

		for _, milestone := range milestones {
			if milestone != nil && milestone.Title == title {
				return milestone, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

// CloseMilestone closes a given milestone.
func closeMilestone(c *CreateOpts, title string) error {
	client, err := c.HTTPClient()
	if err != nil {
		return err
	}

	repo, err := c.BaseRepo()
	if err != nil {
		return err
	}

	milestone, err := getMilestoneByTitle(c, title)
	if err != nil {
		return err
	}

	if milestone == nil {
		return fmt.Errorf("could not find milestone: %q", title)
	}

	closeStateEvent := "close"

	opts := &gitlab.UpdateMilestoneOptions{
		Description: &milestone.Description,
		DueDate:     milestone.DueDate,
		StartDate:   milestone.StartDate,
		StateEvent:  &closeStateEvent,
		Title:       &milestone.Title,
	}

	_, _, err = client.Milestones.UpdateMilestone(
		repo.FullName(),
		milestone.ID,
		opts,
	)

	return err
}

func AssetsFromArgs(args []string) (assets []*upload.ReleaseFile, err error) {
	for _, arg := range args {
		var label string
		fn := arg
		if idx := strings.IndexRune(arg, '#'); idx > 0 {
			fn = arg[0:idx]
			label = arg[idx+1:]
		}

		var fi os.FileInfo
		fi, err = os.Stat(fn)
		if err != nil {
			return
		}

		if label == "" {
			label = fi.Name()
		}

		assets = append(assets, &upload.ReleaseFile{
			Open: func() (io.ReadCloser, error) {
				return os.Open(fn)
			},
			Name:  fi.Name(),
			Label: label,
			Path:  fn,
		})
	}
	return
}
