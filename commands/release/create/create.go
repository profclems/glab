package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/profclems/glab/commands/release/releaseutils"
	"github.com/profclems/glab/commands/release/releaseutils/upload"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/profclems/glab/pkg/surveyext"
	"github.com/profclems/glab/pkg/utils"

	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/glinstance"
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
	RepoOverride     string

	NoteProvided       bool
	ReleaseNotesAction string

	AssetLinks []*upload.ReleaseAsset
	AssetFiles []*upload.ReleaseFile

	IO         *iostreams.IOStreams
	HTTPClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Config     func() (config.Config, error)
}

func NewCmdCreate(f *cmdutils.Factory, runE func(opts *CreateOpts) error) *cobra.Command {
	opts := &CreateOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "create <tag> [<files>...]",
		Short: "Create a new or update a GitLab Release for a repository",
		Long: heredoc.Docf(`Create a new or update a GitLab Release for a repository.

				If the release already exists, glab updates the release with the new info provided.

				If a git tag specified does not yet exist, the release will automatically get created
				from the latest state of the default branch and tagged with the specified tag name.
				Use %[1]s--ref%[1]s to override this.
				The %[1]sref%[1]s can be a commit SHA, another tag name, or a branch name.
				To fetch the new tag locally after the release, do %[1]sgit fetch --tags origin%[1]s.

				To create a release from an annotated git tag, first create one locally with
				git, push the tag to GitLab, then run this command.

				NB: Developer level access to the project is required to create a release.
		`, "`"),
		Args: cmdutils.MinimumArgs(1, "no tag name provided"),
		Example: heredoc.Doc(`
			Interactively create a release
			$ glab release create v1.0.1

			Non-interactively create a release by specifying a note
			$ glab release create v1.0.1 --notes "bugfix release"

			Use release notes from a file
			$ glab release create v1.0.1 -F changelog.md

			Upload a release asset with a display name
			$ glab release create v1.0.1 '/path/to/asset.zip#My display label'

			Upload a release asset with a display name and type
			$ glab release create v1.0.1 '/path/to/asset.png#My display label#image'

			Upload all assets in a specified folder
			$ glab release create v1.0.1 ./dist/*

			Upload all tarballs in a specified folder
			$ glab release create v1.0.1 ./dist/*.tar.gz

			Create a release with assets specified as JSON object
			$ glab release create v1.0.1 --assets-links='
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
			opts.RepoOverride, _ = cmd.Flags().GetString("repo")
			opts.HTTPClient = f.HttpClient
			opts.BaseRepo = f.BaseRepo

			opts.TagName = args[0]

			opts.AssetFiles, err = releaseutils.AssetsFromArgs(args[1:])
			if err != nil {
				return err
			}

			if opts.AssetLinksAsJson != "" {
				err := json.Unmarshal([]byte(opts.AssetLinksAsJson), &opts.AssetLinks)
				if err != nil {
					return fmt.Errorf("failed to parse JSON string: %w", err)
				}
			}

			opts.NoteProvided = cmd.Flags().Changed("notes")
			if opts.NotesFile != "" {
				var b []byte
				var err error
				if opts.NotesFile == "-" {
					b, err = ioutil.ReadAll(opts.IO.In)
					_ = opts.IO.In.Close()
				} else {
					b, err = ioutil.ReadFile(opts.NotesFile)
				}

				if err != nil {
					return err
				}

				opts.Notes = string(b)
				opts.NoteProvided = true
			}

			if runE != nil {
				return runE(opts)
			}

			return createRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "The release name or title")
	cmd.Flags().StringVarP(&opts.Ref, "ref", "r", "", "If a tag specified doesn't exist, the release is created from ref and tagged with the specified tag name. It can be a commit SHA, another tag name, or a branch name.")
	cmd.Flags().StringVarP(&opts.Notes, "notes", "N", "", "The release notes/description. You can use Markdown")
	cmd.Flags().StringVarP(&opts.NotesFile, "notes-file", "F", "", "Read release notes `file`. Specify `-` as value to read from stdin")
	cmd.Flags().StringVarP(&opts.ReleasedAt, "released-at", "D", "", "The `date` when the release is/was ready. Defaults to the current datetime. Expected in ISO 8601 format (2019-03-15T08:00:00Z)")
	cmd.Flags().StringSliceVarP(&opts.Milestone, "milestone", "m", []string{}, "The title of each milestone the release is associated with")
	cmd.Flags().StringVarP(&opts.AssetLinksAsJson, "assets-links", "a", "", "`JSON` string representation of assets links (e.g. `--assets='[{\"name\": \"Asset1\", \"url\":\"https://<domain>/some/location/1\", \"link_type\": \"other\", \"filepath\": \"path/to/file\"}]')`")

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
	var tag *gitlab.Tag
	var resp *gitlab.Response

	if opts.Ref == "" {
		opts.IO.Log(color.ProgressIcon(), "Validating tag", opts.TagName)
		tag, resp, err = client.Tags.GetTag(repo.FullName(), opts.TagName)
		if err != nil && resp != nil && resp.StatusCode != 404 {
			return cmdutils.WrapError(err, "could not fetch tag")
		}
		if tag == nil && resp != nil && resp.StatusCode == 404 {
			opts.IO.Log(color.DotWarnIcon(), "Tag does not exist.")
			opts.IO.Log(color.DotWarnIcon(), "No ref was provided. Tag will be created from the latest state of the default branch")
			project, err := repo.Project(client)
			if err == nil {
				opts.IO.Logf("%s using default branch %q as ref\n", color.ProgressIcon(), project.DefaultBranch)
				opts.Ref = project.DefaultBranch
			}
		}
		// new line
		opts.IO.Log()
	}

	if opts.IO.PromptEnabled() && !opts.NoteProvided {
		editorCommand, err := cmdutils.GetEditor(opts.Config)
		if err != nil {
			return err
		}

		var tagDescription string
		var generatedChangelog string
		if tag == nil {
			tag, _, _ = client.Tags.GetTag(repo.FullName(), opts.TagName)
		}
		if tag != nil {
			tagDescription = tag.Message
		}
		if opts.RepoOverride == "" {
			headRef := opts.TagName
			if tagDescription == "" {
				if opts.Ref != "" {
					headRef = opts.Ref
					brCfg := git.ReadBranchConfig(opts.Ref)
					if brCfg.MergeRef != "" {
						headRef = brCfg.MergeRef
					}
				} else {
					headRef = "HEAD"
				}
			}

			if prevTag, err := detectPreviousTag(headRef); err == nil {
				commits, _ := changelogForRange(fmt.Sprintf("%s..%s", prevTag, headRef))
				generatedChangelog = generateChangelog(commits)
			}
		}

		editorOptions := []string{"Write my own"}
		if generatedChangelog != "" {
			editorOptions = append(editorOptions, "Write using commit log as template")
		}
		if tagDescription != "" {
			editorOptions = append(editorOptions, "Write using git tag message as template")
		}
		editorOptions = append(editorOptions, "Leave blank")

		qs := []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "Release Title (optional)",
					Default: opts.Name,
				},
			},
			{
				Name: "releaseNotesAction",
				Prompt: &survey.Select{
					Message: "Release notes",
					Options: editorOptions,
				},
			},
		}
		err = prompt.Ask(qs, opts)
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}

		var openEditor bool
		var editorContents string

		switch opts.ReleaseNotesAction {
		case "Write my own":
			openEditor = true
		case "Write using commit log as template":
			openEditor = true
			editorContents = generatedChangelog
		case "Write using git tag message as template":
			openEditor = true
			editorContents = tagDescription
		case "Leave blank":
			openEditor = false
		default:
			return fmt.Errorf("invalid action: %v", opts.ReleaseNotesAction)
		}

		if openEditor {
			txt, err := surveyext.Edit(editorCommand, "*.md", editorContents, opts.IO.In, opts.IO.StdOut, opts.IO.StdErr, nil)
			if err != nil {
				return err
			}
			opts.Notes = txt
		}
	}
	start := time.Now()

	opts.IO.Logf("%s creating or updating release %s=%s %s=%s\n",
		color.ProgressIcon(),
		color.Blue("repo"), repo.FullName(),
		color.Blue("tag"), opts.TagName)

	release, resp, err := client.Releases.GetRelease(repo.FullName(), opts.TagName)
	if err != nil && (resp == nil || (resp.StatusCode != 403 && resp.StatusCode != 404)) {
		return releaseFailedErr(err, start)
	}

	var releasedAt time.Time

	if opts.ReleasedAt != "" {
		// Parse the releasedAt to the expected format of the API
		// From the API docs "Expected in ISO 8601 format (2019-03-15T08:00:00Z)".
		releasedAt, err = time.Parse(time.RFC3339[:len(opts.ReleasedAt)], opts.ReleasedAt)
		if err != nil {
			return releaseFailedErr(err, start)
		}
	}

	if opts.Name == "" {
		opts.Name = opts.TagName
	}

	if (resp.StatusCode == 403 || resp.StatusCode == 404) || release == nil {
		createOpts := &gitlab.CreateReleaseOptions{
			Name:    &opts.Name,
			TagName: &opts.TagName,
		}

		if opts.Notes != "" {
			createOpts.Description = &opts.Notes
		}

		if opts.Ref != "" {
			createOpts.Ref = &opts.Ref
		}

		if opts.ReleasedAt != "" {
			createOpts.ReleasedAt = &releasedAt
		}

		if len(opts.Milestone) > 0 {
			createOpts.Milestones = opts.Milestone
		}

		release, _, err = client.Releases.CreateRelease(repo.FullName(), createOpts)

		if err != nil {
			return releaseFailedErr(err, start)
		}
		opts.IO.Logf("%s release created\t%s=%s\n", color.GreenCheck(),
			color.Blue("url"), fmt.Sprintf("%s://%s/%s/-/releases/%s",
				glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault(),
				repo.FullName(), release.TagName))
	} else {
		updateOpts := &gitlab.UpdateReleaseOptions{
			Name: &opts.Name,
		}
		if opts.Notes != "" {
			updateOpts.Description = &opts.Notes
		}

		if opts.ReleasedAt != "" {
			updateOpts.ReleasedAt = &releasedAt
		}

		if len(opts.Milestone) > 0 {
			updateOpts.Milestones = opts.Milestone
		}

		release, _, err = client.Releases.UpdateRelease(repo.FullName(), opts.TagName, updateOpts)
		if err != nil {
			return releaseFailedErr(err, start)
		}

		opts.IO.Logf("%s release updated\t%s=%s\n", color.GreenCheck(),
			color.Blue("url"), fmt.Sprintf("%s://%s/%s/-/releases/%s",
				glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault(),
				repo.FullName(), release.TagName))
	}

	// upload files and create asset link
	if opts.AssetFiles != nil || opts.AssetLinks != nil {
		opts.IO.Logf("\n%s Uploading release assets\n", color.ProgressIcon())
		uploadCtx := upload.Context{
			IO:          opts.IO,
			Client:      client,
			AssetsLinks: opts.AssetLinks,
			AssetFiles:  opts.AssetFiles,
		}
		if err = uploadCtx.UploadFiles(repo.FullName(), release.TagName); err != nil {
			return releaseFailedErr(err, start)
		}

		// create asset link for assets provided as json
		if err = uploadCtx.CreateReleaseAssetLinks(repo.FullName(), release.TagName); err != nil {
			return releaseFailedErr(err, start)
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
				opts.IO.Log(color.FailedIcon(), err.Error())
			} else {
				opts.IO.Logf("%s closed milestone %q\n", color.GreenCheck(), milestone)
			}
		}
	}
	opts.IO.Logf(color.Bold("%s release succeeded after %0.2fs\n"), color.GreenCheck(), time.Since(start).Seconds())
	return nil
}

func releaseFailedErr(err error, start time.Time) error {
	return cmdutils.WrapError(err, fmt.Sprintf("release failed after %0.2fs", time.Since(start).Seconds()))
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

func detectPreviousTag(headRef string) (string, error) {
	cmd := git.GitCommand("describe", "--tags", "--abbrev=0", fmt.Sprintf("%s^", headRef))
	b, err := run.PrepareCmd(cmd).Output()
	return strings.TrimSpace(string(b)), err
}

type logEntry struct {
	Subject string
	Body    string
}

func changelogForRange(refRange string) ([]logEntry, error) {
	cmd := git.GitCommand("-c", "log.ShowSignature=false", "log", "--first-parent", "--reverse", "--pretty=format:%B%x00", refRange)

	b, err := run.PrepareCmd(cmd).Output()
	if err != nil {
		return nil, err
	}

	var entries []logEntry
	for _, cb := range bytes.Split(b, []byte{'\000'}) {
		c := strings.ReplaceAll(string(cb), "\r\n", "\n")
		c = strings.TrimPrefix(c, "\n")
		if c == "" {
			continue
		}
		parts := strings.SplitN(c, "\n\n", 2)
		var body string
		subject := strings.ReplaceAll(parts[0], "\n", " ")
		if len(parts) > 1 {
			body = parts[1]
		}
		entries = append(entries, logEntry{
			Subject: subject,
			Body:    body,
		})
	}

	return entries, nil
}

func generateChangelog(commits []logEntry) string {
	var parts []string
	for _, c := range commits {
		parts = append(parts, fmt.Sprintf("* %s", c.Subject))
		if c.Body != "" {
			parts = append(parts, utils.Indent(c.Body, "  "))
		}
	}
	return strings.Join(parts, "\n\n")
}
