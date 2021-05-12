package login

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/profclems/glab/commands/auth/authutils"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/spf13/cobra"
)

type LoginOptions struct {
	IO     *iostreams.IOStreams
	Config func() (config.Config, error)

	Interactive bool

	Hostname string
	Token    string
}

var opts *LoginOptions

func NewCmdLogin(f *cmdutils.Factory) *cobra.Command {
	opts = &LoginOptions{
		IO:     f.IO,
		Config: f.Config,
	}

	var tokenStdin bool

	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Authenticate with a GitLab instance",
		Long: heredoc.Docf(`
			Authenticate with a GitLab instance.
			You can pass in a token on standard input by using %[1]s--stdin%[1]s.
			The minimum required scopes for the token are: "api", "write_repository".
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ glab auth login
			# authenticate against gitlab.com by reading the token from a file
			$ glab auth login --stdin < myaccesstoken.txt
			# authenticate with a self-hosted GitLab instance
			$ glab auth login --hostname salsa.debian.org
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.IO.PromptEnabled() && !(tokenStdin || opts.Token == "") {
				return &cmdutils.FlagError{Err: errors.New("--stdin or --token required when not running interactively")}
			}

			if opts.Token != "" && tokenStdin {
				return &cmdutils.FlagError{Err: errors.New("specify one of --token or --stdin. You cannot use both flags at the same time")}
			}

			if tokenStdin {
				defer opts.IO.In.Close()
				token, err := ioutil.ReadAll(opts.IO.In)
				if err != nil {
					return fmt.Errorf("failed to read token from STDIN: %w", err)
				}
				opts.Token = strings.TrimSpace(string(token))
			}

			if opts.IO.PromptEnabled() && opts.Token == "" && opts.IO.IsaTTY {
				opts.Interactive = true
			}

			if cmd.Flags().Changed("hostname") {
				if err := hostnameValidator(opts.Hostname); err != nil {
					return &cmdutils.FlagError{Err: fmt.Errorf("error parsing --hostname: %w", err)}
				}
			}

			if !opts.Interactive {
				if opts.Hostname == "" {
					opts.Hostname = glinstance.Default()
				}
			}

			return loginRun()
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "h", "", "The hostname of the GitLab instance to authenticate with")
	cmd.Flags().StringVarP(&opts.Token, "token", "t", "", "Your GitLab access token")
	cmd.Flags().BoolVar(&tokenStdin, "stdin", false, "Read token from standard input")

	return cmd
}

func loginRun() error {
	c := opts.IO.Color()
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	if opts.Token != "" {
		if opts.Hostname == "" {
			return errors.New("empty hostname would leak oauth_token")
		}

		err := cfg.Set(opts.Hostname, "token", opts.Token)
		if err != nil {
			return err
		}

		return cfg.Write()
	}

	hostname := opts.Hostname
	apiHostname := opts.Hostname
	isSelfHosted := false

	if hostname == "" {
		var hostType int
		err := survey.AskOne(&survey.Select{
			Message: "What GitLab instance do you want to log into?",
			Options: []string{
				"GitLab.com",
				"GitLab Self-hosted Instance",
			},
		}, &hostType)

		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}

		isSelfHosted = hostType == 1

		hostname = glinstance.Default()
		apiHostname = hostname
		if isSelfHosted {
			err := survey.AskOne(&survey.Input{
				Message: "GitLab hostname:",
			}, &hostname, survey.WithValidator(hostnameValidator))
			if err != nil {
				return fmt.Errorf("could not prompt: %w", err)
			}
			err = survey.AskOne(&survey.Input{
				Message: "API hostname:",
				Help:    "For instances with different hostname for the API endpoint",
				Default: hostname,
			}, &apiHostname, survey.WithValidator(hostnameValidator))
			if err != nil {
				return fmt.Errorf("could not prompt: %w", err)
			}
		}
	}

	fmt.Fprintf(opts.IO.StdErr, "- Logging into %s\n", hostname)

	if token := config.GetFromEnv("token"); token != "" {
		fmt.Fprintf(opts.IO.StdErr, "%s you have GITLAB_TOKEN or OAUTH_TOKEN environment variable set. Unset if you don't want to use it for glab\n", c.Yellow("!WARNING:"))
	}
	existingToken, _, _ := cfg.GetWithSource(hostname, "token")

	if existingToken != "" && opts.Interactive {
		apiClient, err := cmdutils.LabClientFunc(hostname, cfg, false)
		if err != nil {
			return err
		}

		user, err := api.CurrentUser(apiClient)
		if err == nil {
			username := user.Username
			var keepGoing bool
			err = survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf(
					"You're already logged into %s as %s. Do you want to re-authenticate?",
					hostname,
					username),
				Default: false,
			}, &keepGoing)
			if err != nil {
				return fmt.Errorf("could not prompt: %w", err)
			}

			if !keepGoing {
				return nil
			}
		}
	}

	fmt.Fprintln(opts.IO.StdErr)
	fmt.Fprintln(opts.IO.StdErr, heredoc.Doc(getAccessTokenTip(hostname)))
	var token string
	err = survey.AskOne(&survey.Password{
		Message: "Paste your authentication token:",
	}, &token, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("could not prompt: %w", err)
	}

	if hostname == "" {
		return errors.New("empty hostname would leak token")
	}

	err = cfg.Set(hostname, "token", token)
	if err != nil {
		return err
	}
	err = cfg.Set(hostname, "api_host", apiHostname)
	if err != nil {
		return err
	}

	gitProtocol := "https"
	apiProtocol := "https"

	glabExecutable := "glab"
	if exe, err := os.Executable(); err == nil {
		glabExecutable = exe
	}
	credentialFlow := &authutils.GitCredentialFlow{Executable: glabExecutable}

	if opts.Interactive {
		err = survey.AskOne(&survey.Select{
			Message: "Choose default git protocol",
			Options: []string{
				"SSH",
				"HTTPS",
				"HTTP",
			},
			Default: "HTTPS",
		}, &gitProtocol)
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}

		gitProtocol = strings.ToLower(gitProtocol)
		if opts.Interactive && gitProtocol != "ssh" {
			if err := credentialFlow.Prompt(hostname, gitProtocol); err != nil {
				return err
			}
		}

		if isSelfHosted {
			err = survey.AskOne(&survey.Select{
				Message: "Choose host API protocol",
				Options: []string{
					"HTTPS",
					"HTTP",
				},
				Default: "HTTPS",
			}, &apiProtocol)
			if err != nil {
				return fmt.Errorf("could not prompt: %w", err)
			}

			apiProtocol = strings.ToLower(apiProtocol)
		}

		fmt.Fprintf(opts.IO.StdErr, "- glab config set -h %s git_protocol %s\n", hostname, gitProtocol)
		err = cfg.Set(hostname, "git_protocol", gitProtocol)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdErr, "%s Configured git protocol\n", c.GreenCheck())

		fmt.Fprintf(opts.IO.StdErr, "- glab config set -h %s api_protocol %s\n", hostname, apiProtocol)
		err = cfg.Set(hostname, "api_protocol", apiProtocol)
		if err != nil {
			return err
		}

		fmt.Fprintf(opts.IO.StdErr, "%s Configured API protocol\n", c.GreenCheck())
	}
	apiClient, err := cmdutils.LabClientFunc(hostname, cfg, false)
	if err != nil {
		return err
	}

	user, err := api.CurrentUser(apiClient)
	if err != nil {
		return fmt.Errorf("error using api: %w", err)
	}
	username := user.Username

	err = cfg.Set(hostname, "user", username)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	if credentialFlow.ShouldSetup() {
		err := credentialFlow.Setup(hostname, gitProtocol, username, token)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(opts.IO.StdErr, "%s Logged in as %s\n", c.GreenCheck(), c.Bold(username))

	return nil
}

func hostnameValidator(v interface{}) error {
	val := fmt.Sprint(v)
	if len(strings.TrimSpace(val)) < 1 {
		return errors.New("a value is required")
	}
	re := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-.]*$`)
	if !re.MatchString(val) {
		return fmt.Errorf("invalid hostname %q", val)
	}
	return nil
}

func getAccessTokenTip(hostname string) string {
	glHostname := hostname
	if glHostname == "" {
		glHostname = glinstance.OverridableDefault()
	}
	return fmt.Sprintf(`
	Tip: you can generate a Personal Access Token here https://%s/profile/personal_access_tokens
	The minimum required scopes are 'api' and 'write_repository'.`, glHostname)
}
