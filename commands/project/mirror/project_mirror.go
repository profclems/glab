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
		Long:  `Mirrors a project/repository to the specified location using pull or push method.`,
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			# Mirror a project/repository
			# TODO: Add Examples
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

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
					"argument direction value should be pull or push",
				)
			}

			if opts.Direction == "pull" && opts.AllowDivergence {
				fmt.Fprintf(
					f.IO.StdOut,
					"[Warning] allow-divergence flag has no effect for pull mirror and is ignored.\n",
				)
			}

			opts.URL = strings.TrimSpace(opts.URL)

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

	_ = projectMirrorCmd.MarkFlagRequired("url")
	_ = projectMirrorCmd.MarkFlagRequired("direction")

	return projectMirrorCmd
}

func runProjectMirror(f *cmdutils.Factory, opts *MirrorOptions) error {

	if opts.Direction == "push" {
		return createPushMirror(f, opts)
	} else {
		return createPullMirror(f, opts)
	}
}

func createPushMirror(f *cmdutils.Factory, opts *MirrorOptions) error {
	var pm *gitlab.ProjectMirror
	var err error
	var pushOptions = api.CreatePushMirrorOptions{
		Url:                   opts.URL,
		Enabled:               opts.Enabled,
		OnlyProtectedBranches: opts.ProtectedBranchesOnly,
		KeepDivergentRefs:     opts.AllowDivergence,
	}
	pm, err = api.CreatePushMirror(
		opts.APIClient,
		opts.ProjectID,
		&pushOptions,
	)
	if err != nil {
		return cmdutils.WrapError(err, "Failed to create Push Mirror")
	}
	greenCheck := f.IO.Color().Green("✓")
	fmt.Fprintf(
		f.IO.StdOut,
		"%s Created Push Mirror for %s (%d) on GitLab at %s (%d)\n",
		greenCheck, pm.URL, pm.ID, opts.ProjectName, opts.ProjectID,
	)
	return err
}

func createPullMirror(f *cmdutils.Factory, opts *MirrorOptions) error {
	var pullOptions = api.CreatePullMirrorOptions{
		Url:                   opts.URL,
		Enabled:               opts.Enabled,
		OnlyProtectedBranches: opts.ProtectedBranchesOnly,
	}
	err := api.CreatePullMirror(
		opts.APIClient,
		opts.ProjectID,
		&pullOptions,
	)
	if err != nil {
		return cmdutils.WrapError(err, "Failed to create Pull Mirror")
	}
	greenCheck := f.IO.Color().Green("✓")
	fmt.Fprintf(
		f.IO.StdOut,
		"%s Created Pull Mirror for %s on GitLab at %s (%d)\n",
		greenCheck, opts.URL, opts.ProjectName, opts.ProjectID,
	)
	return err
}
