package clone

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/xanzy/go-gitlab"
)

type CloneOptions struct {
	GroupName         string
	IncludeSubgroups  bool
	WithMREnabled     bool
	WithIssuesEnabled bool
	WithShared        bool
	Archived          bool
	ArchivedSet       bool
	Visibility        string
	Owned             bool
	GitFlags          []string
	Dir               string
	Host              string
	Protocol          string

	RemoteArgs *glrepo.RemoteArgs

	IO        *iostreams.IOStreams
	APIClient *api.Client
	Config    func() (config.Config, error)

	CurrentUser *gitlab.User
}

type ContextOpts struct {
	Project *gitlab.Project
	Repo    string
}

func NewCmdClone(f *cmdutils.Factory, runE func(*CloneOptions, *ContextOpts) error) *cobra.Command {
	opts := &CloneOptions{
		IO:     f.IO,
		Config: f.Config,
	}

	ctxOpts := &ContextOpts{}

	var repoCloneCmd = &cobra.Command{
		Use:   "clone <repo> [<dir>] [-- [<gitflags>...]]",
		Short: `Clone a Gitlab repository/project`,
		Example: heredoc.Doc(`
	$ glab repo clone profclems/glab

	$ glab repo clone https://gitlab.com/profclems/glab

	$ glab repo clone profclems/glab mydirectory  # Clones repo into mydirectory

	$ glab repo clone glab   # clones repo glab for current user 

	$ glab repo clone 4356677   # finds the project by the ID provided and clones it

	# Clone all repos in a group
	$ glab repo clone -g everyonecancontribute  

	# Clone from a self-hosted instance
	$ GITLAB_HOST=salsa.debian.org glab repo clone myrepo  
	`),
		Long: heredoc.Doc(`
Clone a GitLab repository/project

	Clone supports these shorthands
	- repo
	- namespace/repo
	- org/group/repo
	- project ID
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if nArgs := len(args); nArgs > 0 {
				ctxOpts.Repo = args[0]
				if nArgs > 1 {
					opts.Dir = args[1]
				}
				opts.GitFlags = args[1:]
			}

			if ctxOpts.Repo == "" && opts.GroupName == "" {
				return &cmdutils.FlagError{Err: fmt.Errorf("specify repo argument or use --group flag to specify a group to clone all repos from the group")}
			}

			if runE != nil {
				return runE(opts, ctxOpts)
			}

			opts.Host = glinstance.OverridableDefault()
			opts.ArchivedSet = cmd.Flags().Changed("archived")

			cfg, err := opts.Config()
			if err != nil {
				return err
			}
			opts.APIClient, err = api.NewClientWithCfg(opts.Host, cfg, false)
			if err != nil {
				return err
			}

			opts.CurrentUser, err = api.CurrentUser(opts.APIClient.Lab())
			if err != nil {
				return err
			}

			opts.Protocol, _ = cfg.Get(opts.Host, "git_protocol")

			opts.RemoteArgs = &glrepo.RemoteArgs{
				Protocol: opts.Protocol,
				Token:    opts.APIClient.Token(),
				Url:      opts.Host,
				Username: opts.CurrentUser.Username,
			}

			if opts.GroupName != "" {
				return groupClone(opts, ctxOpts)
			}

			return cloneRun(opts, ctxOpts)
		},
	}

	repoCloneCmd.Flags().StringVarP(&opts.GroupName, "group", "g", "", "Specify group to clone repositories from")
	repoCloneCmd.Flags().BoolVarP(&opts.Archived, "archived", "a", false, "Limit by archived status. Used with --group flag")
	repoCloneCmd.Flags().BoolVarP(&opts.IncludeSubgroups, "include-subgroups", "G", true, "Include projects in subgroups of this group. Default is true. Used with --group flag")
	repoCloneCmd.Flags().BoolVarP(&opts.Owned, "mine", "m", false, "Limit by projects in the group owned by the current authenticated user. Used with --group flag")
	repoCloneCmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "", "Limit by visibility {public, internal, or private}. Used with --group flag")
	repoCloneCmd.Flags().BoolVarP(&opts.WithIssuesEnabled, "with-issues-enabled", "I", false, "Limit by projects with issues feature enabled. Default is false. Used with --group flag")
	repoCloneCmd.Flags().BoolVarP(&opts.WithMREnabled, "with-mr-enabled", "M", false, "Limit by projects with issues feature enabled. Default is false. Used with --group flag")
	repoCloneCmd.Flags().BoolVarP(&opts.WithShared, "with-shared", "S", false, "Include projects shared to this group. Default is false. Used with --group flag")

	repoCloneCmd.Flags().SortFlags = false
	repoCloneCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err == pflag.ErrHelp {
			return err
		}
		return &cmdutils.FlagError{Err: fmt.Errorf("%w\nSeparate git clone flags with '--'.", err)}
	})

	return repoCloneCmd
}

func groupClone(opts *CloneOptions, ctxOpts *ContextOpts) error {
	c := opts.IO.Color()
	ListGroupProjectOpts := &gitlab.ListGroupProjectsOptions{}
	if opts.WithShared {
		ListGroupProjectOpts.WithShared = gitlab.Bool(true)
	}
	if opts.WithMREnabled {
		ListGroupProjectOpts.WithMergeRequestsEnabled = gitlab.Bool(true)
	}
	if opts.WithIssuesEnabled {
		ListGroupProjectOpts.WithIssuesEnabled = gitlab.Bool(true)
	}
	if opts.Owned {
		ListGroupProjectOpts.Owned = gitlab.Bool(true)
	}
	if opts.ArchivedSet {
		ListGroupProjectOpts.Archived = gitlab.Bool(opts.Archived)
	}
	if opts.IncludeSubgroups {
		ListGroupProjectOpts.IncludeSubgroups = gitlab.Bool(true)
	}
	if opts.Visibility != "" {
		ListGroupProjectOpts.Visibility = gitlab.Visibility(gitlab.VisibilityValue(opts.Visibility))
	}
	ListGroupProjectOpts.PerPage = 100 //TODO: Allow user to specify the page and limit
	projects, err := api.ListGroupProjects(opts.APIClient.Lab(), opts.GroupName, ListGroupProjectOpts)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Fprintf(opts.IO.StdErr, "Group %q does not have any projects\n", opts.GroupName)
		return cmdutils.SilentError
	}
	var finalOutput []string
	for _, project := range projects {
		ctxOpt := *ctxOpts
		ctxOpt.Project = project
		ctxOpt.Repo = project.PathWithNamespace
		err = cloneRun(opts, &ctxOpt)
		if err != nil {
			finalOutput = append(finalOutput, fmt.Sprintf("%s %s - Error: %q", c.RedCheck(), project.PathWithNamespace, err.Error()))
		} else {
			finalOutput = append(finalOutput, fmt.Sprintf("%s %s", c.GreenCheck(), project.PathWithNamespace))
		}
	}

	// Print error/success msgs in human-readable formats
	for _, out := range finalOutput {
		fmt.Fprintln(opts.IO.StdOut, out)
	}
	if err != nil { // if any error came up
		return cmdutils.SilentError
	}
	return nil
}

func cloneRun(opts *CloneOptions, ctxOpts *ContextOpts) (err error) {
	if !git.IsValidURL(ctxOpts.Repo) {
		// Assuming that repo is a project ID if it is an integer
		if _, err := strconv.ParseInt(ctxOpts.Repo, 10, 64); err != nil {
			// Assuming that "/" in the project name means its owned by an organisation
			if !strings.Contains(ctxOpts.Repo, "/") {
				ctxOpts.Repo = fmt.Sprintf("%s/%s", opts.CurrentUser.Username, ctxOpts.Repo)
			}
		}
		if ctxOpts.Project == nil {
			ctxOpts.Project, err = api.GetProject(opts.APIClient.Lab(), ctxOpts.Repo)
			if err != nil {
				return
			}
		}
		ctxOpts.Repo, err = glrepo.RemoteURL(ctxOpts.Project, opts.RemoteArgs)
		if err != nil {
			return
		}
	} else if !strings.HasSuffix(ctxOpts.Repo, ".git") {
		ctxOpts.Repo += ".git"
	}
	_, err = git.RunClone(ctxOpts.Repo, opts.GitFlags)
	if err != nil {
		return
	}
	// Cloned project was a fork belonging to the user; user is
	// treating fork's ssh/https url as origin. Add upstream as remote pointing
	// to forked repo's ssh/https url depending on the users preferred protocol
	if ctxOpts.Project != nil {
		if ctxOpts.Project.ForkedFromProject != nil && strings.Contains(ctxOpts.Project.PathWithNamespace, opts.CurrentUser.Username) {
			if opts.Dir == "" {
				opts.Dir = "./" + ctxOpts.Project.Path
			}
			fProject, err := api.GetProject(opts.APIClient.Lab(), ctxOpts.Project.ForkedFromProject.PathWithNamespace)
			if err != nil {
				return err
			}
			repoURL, err := glrepo.RemoteURL(fProject, opts.RemoteArgs)
			if err != nil {
				return err
			}
			err = git.AddUpstreamRemote(repoURL, opts.Dir)
			if err != nil {
				return err
			}
		}
	}
	return
}
