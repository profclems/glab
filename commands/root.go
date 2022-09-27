package commands

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	aliasCmd "github.com/profclems/glab/commands/alias"
	apiCmd "github.com/profclems/glab/commands/api"
	authCmd "github.com/profclems/glab/commands/auth"
	pipelineCmd "github.com/profclems/glab/commands/ci"
	"github.com/profclems/glab/commands/cmdutils"
	completionCmd "github.com/profclems/glab/commands/completion"
	configCmd "github.com/profclems/glab/commands/config"
	"github.com/profclems/glab/commands/help"
	issueCmd "github.com/profclems/glab/commands/issue"
	labelCmd "github.com/profclems/glab/commands/label"
	mrCmd "github.com/profclems/glab/commands/mr"
	projectCmd "github.com/profclems/glab/commands/project"
	releaseCmd "github.com/profclems/glab/commands/release"
	snippetCmd "github.com/profclems/glab/commands/snippet"
	sshCmd "github.com/profclems/glab/commands/ssh-key"
	updateCmd "github.com/profclems/glab/commands/update"
	userCmd "github.com/profclems/glab/commands/user"
	variableCmd "github.com/profclems/glab/commands/variable"
	versionCmd "github.com/profclems/glab/commands/version"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdRoot is the main root/parent command
func NewCmdRoot(f *cmdutils.Factory, version, buildDate string) *cobra.Command {
	c := f.IO.Color()
	var rootCmd = &cobra.Command{
		Use:           "glab <command> <subcommand> [flags]",
		Short:         "A GitLab CLI Tool",
		Long:          `GLab is an open source GitLab CLI tool bringing GitLab to your command line`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Annotations: map[string]string{
			"help:environment": heredoc.Doc(`
			GITLAB_TOKEN: an authentication token for API requests. Setting this avoids being
			prompted to authenticate and overrides any previously stored credentials.
			Can be set in the config with 'glab config set token xxxxxx'

			GITLAB_HOST or GL_HOST: specify the url of the gitlab server if self hosted (eg: https://gitlab.example.com). Default is https://gitlab.com.

			REMOTE_ALIAS or GIT_REMOTE_URL_VAR: git remote variable or alias that contains the gitlab url.
			Can be set in the config with 'glab config set remote_alias origin'

			VISUAL, EDITOR (in order of precedence): the editor tool to use for authoring text.
			Can be set in the config with 'glab config set editor vim'

			BROWSER: the web browser to use for opening links.
			Can be set in the config with 'glab config set browser mybrowser'

			GLAMOUR_STYLE: environment variable to set your desired markdown renderer style
			Available options are (dark|light|notty) or set a custom style
			https://github.com/charmbracelet/glamour#styles

			NO_PROMPT: set to 1 (true) or 0 (false) to disable and enable prompts respectively

			NO_COLOR: set to any value to avoid printing ANSI escape sequences for color output.

			FORCE_HYPERLINKS: set to 1 to force hyperlinks to be output, even when not outputing to a TTY

			GLAB_CONFIG_DIR: set to a directory path to override the global configuration location 
		`),
			"help:feedback": heredoc.Docf(`
			Encountered a bug or want to suggest a feature?
			Open an issue using '%s'
		`, c.Bold(c.Yellow("glab issue create -R profclems/glab"))),
		},
	}

	rootCmd.SetOut(f.IO.StdOut)
	rootCmd.SetErr(f.IO.StdErr)

	rootCmd.PersistentFlags().Bool("help", false, "Show help for command")
	rootCmd.SetHelpFunc(func(command *cobra.Command, args []string) {
		help.RootHelpFunc(f.IO.Color(), command, args)
	})
	rootCmd.SetUsageFunc(help.RootUsageFunc)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if errors.Is(err, pflag.ErrHelp) {
			return err
		}
		return &cmdutils.FlagError{Err: err}
	})

	formattedVersion := versionCmd.Scheme(version, buildDate)
	rootCmd.SetVersionTemplate(formattedVersion)
	rootCmd.Version = formattedVersion

	// Child commands
	rootCmd.AddCommand(aliasCmd.NewCmdAlias(f))
	rootCmd.AddCommand(configCmd.NewCmdConfig(f))
	rootCmd.AddCommand(completionCmd.NewCmdCompletion(f.IO))
	rootCmd.AddCommand(versionCmd.NewCmdVersion(f.IO, version, buildDate))
	rootCmd.AddCommand(updateCmd.NewCheckUpdateCmd(f.IO, version))
	rootCmd.AddCommand(authCmd.NewCmdAuth(f))

	// the commands below require apiClient and resolved repos
	f.BaseRepo = resolvedBaseRepo(f)
	cmdutils.HTTPClientFactory(f) // Initialize HTTP Client

	rootCmd.AddCommand(issueCmd.NewCmdIssue(f))
	rootCmd.AddCommand(labelCmd.NewCmdLabel(f))
	rootCmd.AddCommand(mrCmd.NewCmdMR(f))
	rootCmd.AddCommand(pipelineCmd.NewCmdCI(f))
	rootCmd.AddCommand(projectCmd.NewCmdRepo(f))
	rootCmd.AddCommand(releaseCmd.NewCmdRelease(f))
	rootCmd.AddCommand(sshCmd.NewCmdSSHKey(f))
	rootCmd.AddCommand(userCmd.NewCmdUser(f))
	rootCmd.AddCommand(variableCmd.NewVariableCmd(f))
	rootCmd.AddCommand(apiCmd.NewCmdApi(f, nil))
	rootCmd.AddCommand(snippetCmd.NewCmdSnippet(f))

	rootCmd.Flags().BoolP("version", "v", false, "show glab version information")
	return rootCmd
}

func resolvedBaseRepo(f *cmdutils.Factory) func() (glrepo.Interface, error) {
	return func() (glrepo.Interface, error) {
		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}
		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}
		repoContext, err := glrepo.ResolveRemotesToRepos(remotes, httpClient, "")
		if err != nil {
			return nil, err
		}
		baseRepo, err := repoContext.BaseRepo(true)
		if err != nil {
			return nil, err
		}

		return baseRepo, nil
	}
}
