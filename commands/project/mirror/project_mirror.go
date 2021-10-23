package mirror

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type MirrorOptions struct {
	URL                   string
	Direction             string
	Enabled               bool
	ProtectedBranchesOnly bool
	AllowDivergence       bool
	ProjectID             int
	ProjectName           string

	APIClient *gitlab.Client
	IO        *iostreams.IOStreams
	Repo      glrepo.Interface
}

func NewCmdMirror(f *cmdutils.Factory) *cobra.Command {
	opts := MirrorOptions{
		IO: f.IO,
	}

	var projectMirrorCmd = &cobra.Command{
		Use:   "mirror [flags]",
		Short: "Mirror a project/repository",
		Long:  `Mirrors a project/repository to the specified location using speicified method.`,
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			# Mirror a project/repository
			# TODO: Add Examples
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) != 1 {
				return cmdutils.WrapError(
					errors.New("project/repository required"),
					"mirror takes only one positional argument project/repository",
				)
			}

			opts.ProjectName = args[0]
			if len(strings.Split(opts.ProjectName, "/")) != 2 {
				return cmdutils.WrapError(
					errors.New("ill-formatted argument project/repository"),
					"argument should be in the form of project/repository",
				)
			}

			if opts.Direction != "pull" && opts.Direction != "push" {
				return cmdutils.WrapError(
					errors.New("invalid choice for --direction"),
					"argument direction value should be pull or push, default is pull",
				)
			}

			opts.URL = strings.TrimSpace(opts.URL)

			if opts.URL == "" {
				return cmdutils.WrapError(
					errors.New("argument URL is required"),
					"argument URL must be provided",
				)
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			opts.APIClient = apiClient

			project, err := api.GetProject(apiClient, opts.ProjectName)
			if err != nil {
				return cmdutils.WrapError(err, "project/repository not found")
			}
			opts.ProjectID = project.ID
			return runProjectMirror(f, &opts)
		},
	}
	projectMirrorCmd.Flags().StringVar(&opts.URL, "url", "", "The target URL to which the repository is mirrored.")
	projectMirrorCmd.Flags().StringVar(&opts.Direction, "direction", "pull", "Mirror direction")
	projectMirrorCmd.Flags().BoolVar(&opts.Enabled, "enabled", true, "Determines if the mirror is enabled.")
	projectMirrorCmd.Flags().BoolVar(&opts.ProtectedBranchesOnly, "protected-branches-only", false, "Determines if only protected branches are mirrored.")
	projectMirrorCmd.Flags().BoolVar(&opts.AllowDivergence, "allow-divergence", false, "Determines if divergent refs are skipped.")

	return projectMirrorCmd
}

func runProjectMirror(f *cmdutils.Factory, opts *MirrorOptions) error {
	var pm *gitlab.ProjectMirror
	var err error
	if opts.Direction == "push" {
		pm, err = api.CreatePushMirror(
			opts.APIClient,
			opts.ProjectID,
			opts.URL,
			opts.Enabled,
			opts.ProtectedBranchesOnly,
			opts.AllowDivergence,
		)
		if err != nil {
			return cmdutils.WrapError(err, "Failed to create mirror")
		}
		greenCheck := f.IO.Color().Green("✓")
		fmt.Fprintf(
			f.IO.StdOut,
			"%s Created %s Mirror for %s (%d) on GitLab at %s (%d)\n",
			greenCheck, strings.ToTitle(opts.Direction), pm.URL, pm.ID, opts.ProjectName, opts.ProjectID,
		)
	} else {
		_, err = api.CreatePullMirror()
	}

	return err
}
