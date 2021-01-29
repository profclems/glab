package status

import (
	"fmt"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
)

type StatusOpts struct {
	Hostname  string
	ShowToken bool

	HttpClientOverride func(token, hostname string) (*api.Client, error) // used in tests to mock http client
	IO                 *iostreams.IOStreams
	Config             func() (config.Config, error)
}

func NewCmdStatus(f *cmdutils.Factory, runE func(*StatusOpts) error) *cobra.Command {
	opts := &StatusOpts{
		IO:     f.IO,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "status",
		Args:  cobra.ExactArgs(0),
		Short: "View authentication status",
		Long: heredoc.Doc(`Verifies and displays information about your authentication state.
			
			This command tests the authentication states of all known GitLab instances in the config file and reports issues if any
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runE != nil {
				return runE(opts)
			}

			return statusRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "hostname", "h", "", "Check a specific instance's authentication status")
	cmd.Flags().BoolVarP(&opts.ShowToken, "show-token", "t", false, "Display the auth token")

	return cmd
}

func statusRun(opts *StatusOpts) error {
	c := opts.IO.Color()
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	stderr := opts.IO.StdErr

	statusInfo := map[string][]string{}

	instances, err := cfg.Hosts()
	if len(instances) == 0 || err != nil {
		fmt.Fprintf(stderr,
			"No GitLab instance has been authenticated with glab. Run `%s` to authenticate.\n", c.Bold("glab auth login"))
		return cmdutils.SilentError
	}

	var hostNotAuthenticated bool
	if opts.Hostname != "" {
		hostNotAuthenticated = true
	}

	for _, instance := range instances {
		if opts.Hostname != "" && opts.Hostname != instance {
			continue
		}
		hostNotAuthenticated = false
		statusInfo[instance] = []string{}
		addMsg := func(x string, ys ...interface{}) {
			statusInfo[instance] = append(statusInfo[instance], fmt.Sprintf(x, ys...))
		}

		token, tokenSource, _ := cfg.GetWithSource(instance, "token")
		apiClient, err := api.NewClientWithCfg(instance, cfg, false)
		if opts.HttpClientOverride != nil {
			apiClient, _ = opts.HttpClientOverride(token, instance)
		}
		if err == nil {
			user, err := api.CurrentUser(apiClient.Lab())
			if err != nil {
				addMsg("%s %s: api call failed: %s", c.FailedIcon(), instance, err)
			} else {
				addMsg("%s Logged in to %s as %s (%s)", c.GreenCheck(), instance, c.Bold(user.Username), tokenSource)
			}
		} else {
			addMsg("%s %s: failed to initialize api client: %s", c.FailedIcon(), instance, err)
		}
		proto, _ := cfg.Get(instance, "git_protocol")
		if proto != "" {
			addMsg("%s Git operations for %s configured to use %s protocol.",
				c.GreenCheck(), instance, c.Bold(proto))
		}
		apiProto, _ := cfg.Get(instance, "api_protocol")
		apiEndpoint := glinstance.APIEndpoint(instance, apiProto)
		graphQLEndpoint := glinstance.GraphQLEndpoint(instance, apiProto)
		if apiProto != "" {
			addMsg("%s API calls for %s are made over %s protocol",
				c.GreenCheck(), instance, c.Bold(apiProto))
			addMsg("%s REST API Endpoint: %s",
				c.GreenCheck(), c.Bold(apiEndpoint))
			addMsg("%s GraphQL Endpoint: %s",
				c.GreenCheck(), c.Bold(graphQLEndpoint))
		}
		if token != "" {
			tokenDisplay := "********************"
			if opts.ShowToken {
				tokenDisplay = token
			}
			addMsg("%s Token: %s", c.GreenCheck(), tokenDisplay)
			if !api.IsValidToken(token) {
				addMsg("%s Invalid token provided", c.WarnIcon())
			}
		} else {
			addMsg("%s No token provided", c.FailedIcon())
		}
	}

	if opts.Hostname != "" && hostNotAuthenticated {
		fmt.Fprintf(stderr, "%s %s not authenticated with glab. Run `%s %s` to authenticate", c.FailedIcon(), opts.Hostname, c.Bold("glab auth login --hostname"), c.Bold(opts.Hostname))
		return cmdutils.SilentError
	}

	for _, instance := range instances {
		lines, ok := statusInfo[instance]
		if !ok {
			continue
		}
		fmt.Fprintf(stderr, "%s\n", c.Bold(instance))
		for _, line := range lines {
			fmt.Fprintf(stderr, "  %s\n", line)
		}
	}
	return nil
}
