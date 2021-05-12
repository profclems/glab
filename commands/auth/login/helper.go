package login

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/iostreams"

	"github.com/spf13/cobra"
)

const tokenUser = "oauth2"

type configExt interface {
	GetWithSource(string, string) (string, string, error)
}

type CredentialOptions struct {
	IO     *iostreams.IOStreams
	Config func() (configExt, error)

	Operation string
}

func NewCmdCredential(f *cmdutils.Factory, runF func(*CredentialOptions) error) *cobra.Command {
	opts := &CredentialOptions{
		IO: f.IO,
		Config: func() (configExt, error) {
			return f.Config()
		},
	}

	cmd := &cobra.Command{
		Use:    "git-credential",
		Args:   cobra.ExactArgs(1),
		Short:  "Implements git credential helper manager",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Operation = args[0]

			if runF != nil {
				return runF(opts)
			}
			return helperRun(opts)
		},
	}

	return cmd
}

func helperRun(opts *CredentialOptions) error {
	if opts.Operation == "store" {
		// We pretend to implement the "store" operation, but do nothing since we already have a cached token.
		return cmdutils.SilentError
	}

	if opts.Operation != "get" {
		return fmt.Errorf("glab auth git-credential: %q is an invalid operation", opts.Operation)
	}

	expectedParams := map[string]string{}

	s := bufio.NewScanner(opts.IO.In)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		key, value := parts[0], parts[1]
		if key == "url" {
			u, err := url.Parse(value)
			if err != nil {
				return err
			}
			expectedParams["protocol"] = u.Scheme
			expectedParams["host"] = u.Host
			expectedParams["path"] = u.Path
			expectedParams["username"] = u.User.Username()
			expectedParams["password"], _ = u.User.Password()
		} else {
			expectedParams[key] = value
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	if expectedParams["protocol"] != "https" && expectedParams["protocol"] != "http" {
		return cmdutils.SilentError
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	var gotUser string
	gotToken, source, _ := cfg.GetWithSource(expectedParams["host"], "token")
	if strings.HasSuffix(source, "_TOKEN") {
		gotUser = tokenUser
	} else {
		gotUser, _, _ = cfg.GetWithSource(expectedParams["host"], "user")
	}

	if gotUser == "" || gotToken == "" {
		return cmdutils.SilentError
	}

	if expectedParams["username"] != "" && gotUser != tokenUser && !strings.EqualFold(expectedParams["username"], gotUser) {
		return cmdutils.SilentError
	}

	fmt.Fprintf(opts.IO.StdOut, "protocol=%s\n", expectedParams["protocol"])
	fmt.Fprintf(opts.IO.StdOut, "host=%s\n", expectedParams["host"])
	fmt.Fprintf(opts.IO.StdOut, "username=%s\n", gotUser)
	fmt.Fprintf(opts.IO.StdOut, "password=%s\n", gotToken)

	return nil
}
