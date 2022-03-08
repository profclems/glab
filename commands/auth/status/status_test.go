package status

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/pkg/httpmock"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func Test_NewCmdStatus(t *testing.T) {
	tests := []struct {
		name  string
		cli   string
		wants StatusOpts
	}{
		{
			name:  "no arguments",
			cli:   "",
			wants: StatusOpts{},
		},
		{
			name: "hostname set",
			cli:  "--hostname gitlab.gnome.org",
			wants: StatusOpts{
				Hostname: "gitlab.gnome.org",
			},
		},
		{
			name: "show token",
			cli:  "--show-token",
			wants: StatusOpts{
				ShowToken: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &cmdutils.Factory{}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *StatusOpts
			cmd := NewCmdStatus(f, func(opts *StatusOpts) error {
				gotOpts = opts
				return nil
			})

			// TODO cobra hack-around
			cmd.Flags().BoolP("help", "x", false, "")

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.Hostname, gotOpts.Hostname)
		})
	}
}

func Test_statusRun(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.alpinelinux.org:
    token: xxxxxxxxxxxxxxxxxxxx
    git_protocol: ssh
    api_protocol: https
  gitlab.foo.bar:
    token: glpat-xxxxxxxxxxxxxxxxxxxx
    git_protocol: ssh
    api_protocol: https
  another.host:
    token: isinvalid
  gl.io:
    token: 
`, "")()

	cfgFile := config.ConfigFile()
	configs, err := config.ParseConfig("config.yml")
	assert.Nil(t, err)

	tests := []struct {
		name    string
		opts    *StatusOpts
		wantErr bool
		stderr  string
	}{
		{
			name: "hostname set with old token format",
			opts: &StatusOpts{
				Hostname: "gitlab.alpinelinux.org",
			},
			wantErr: false,
			stderr: fmt.Sprintf(`gitlab.alpinelinux.org
  ✓ Logged in to gitlab.alpinelinux.org as john_smith (%s)
  ✓ Git operations for gitlab.alpinelinux.org configured to use ssh protocol.
  ✓ API calls for gitlab.alpinelinux.org are made over https protocol
  ✓ REST API Endpoint: https://gitlab.alpinelinux.org/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.alpinelinux.org/api/graphql/
  ✓ Token: **************************
`, cfgFile),
		},
		{
			name: "hostname set with new token format",
			opts: &StatusOpts{
				Hostname: "gitlab.foo.bar",
			},
			wantErr: false,
			stderr: fmt.Sprintf(`gitlab.foo.bar
  ✓ Logged in to gitlab.foo.bar as john_doe (%s)
  ✓ Git operations for gitlab.foo.bar configured to use ssh protocol.
  ✓ API calls for gitlab.foo.bar are made over https protocol
  ✓ REST API Endpoint: https://gitlab.foo.bar/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.foo.bar/api/graphql/
  ✓ Token: **************************
`, cfgFile),
		},
		{
			name: "instance not authenticated",
			opts: &StatusOpts{
				Hostname: "invalid.instance",
			},
			wantErr: true,
			stderr:  "x invalid.instance not authenticated with glab. Run `glab auth login --hostname invalid.instance` to authenticate",
		},
	}

	fakeHTTP := &httpmock.Mocker{
		MatchURL: httpmock.HostAndPath,
	}
	defer fakeHTTP.Verify(t)

	fakeHTTP.RegisterResponder("GET", "https://gitlab.alpinelinux.org/api/v4/user", httpmock.NewStringResponse(200, `
		{
  			"username": "john_smith"
		}
	`))
	fakeHTTP.RegisterResponder("GET", "https://gitlab.foo.bar/api/v4/user", httpmock.NewStringResponse(200, `
		{
  			"username": "john_doe"
		}
	`))

	client := func(token, hostname string) (*api.Client, error) {
		return api.TestClient(&http.Client{Transport: fakeHTTP}, token, hostname, false)
	}
	// FIXME: something fishy is occurring here as without making a first call to client function, httpMock does not work
	_, _ = client("", "gitlab.com")

	for _, tt := range tests {
		io, _, stdout, stderr := iostreams.Test()
		tt.opts.Config = func() (config.Config, error) {
			return configs, nil
		}
		tt.opts.IO = io
		tt.opts.HttpClientOverride = client
		t.Run(tt.name, func(t *testing.T) {
			if err := statusRun(tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("statusRun() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, stdout.String(), "")
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}

func Test_statusRun_noHostnameSpecified(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.alpinelinux.org:
    token: xxxxxxxxxxxxxxxxxxxx
    git_protocol: ssh
    api_protocol: https
  another.host:
    token: isinvalid
  gl.io:
    token: 
`, "")()

	fakeHTTP := &httpmock.Mocker{
		MatchURL: httpmock.HostAndPath,
	}
	defer fakeHTTP.Verify(t)

	cfgFile := config.ConfigFile()

	fakeHTTP.RegisterResponder("GET", "https://gitlab.alpinelinux.org/api/v4/user", httpmock.NewStringResponse(200, `
		{
  			"username": "john_smith"
		}
	`))

	fakeHTTP.RegisterResponder("GET", "https://another.host/api/v4/user?u=1", httpmock.NewStringResponse(401, `
		{
  			"message": "invalid token"
		}
	`))

	fakeHTTP.RegisterResponder("GET", "https://gl.io/api/v4/user?u=1", httpmock.NewStringResponse(401, `
		{
  			"message": "no token provided"
		}
	`))

	expectedOutput := fmt.Sprintf(`gitlab.alpinelinux.org
  ✓ Logged in to gitlab.alpinelinux.org as john_smith (%s)
  ✓ Git operations for gitlab.alpinelinux.org configured to use ssh protocol.
  ✓ API calls for gitlab.alpinelinux.org are made over https protocol
  ✓ REST API Endpoint: https://gitlab.alpinelinux.org/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.alpinelinux.org/api/graphql/
  ✓ Token: **************************
another.host
  x another.host: api call failed: GET https://another.host/api/v4/user: 401 {message: invalid token}
  ✓ Git operations for another.host configured to use ssh protocol.
  ✓ API calls for another.host are made over https protocol
  ✓ REST API Endpoint: https://another.host/api/v4/
  ✓ GraphQL Endpoint: https://another.host/api/graphql/
  ✓ Token: **************************
  ! Invalid token provided
gl.io
  x gl.io: api call failed: GET https://gl.io/api/v4/user: 401 {message: no token provided}
  ✓ Git operations for gl.io configured to use ssh protocol.
  ✓ API calls for gl.io are made over https protocol
  ✓ REST API Endpoint: https://gl.io/api/v4/
  ✓ GraphQL Endpoint: https://gl.io/api/graphql/
  x No token provided
`, cfgFile)

	configs, err := config.ParseConfig("config.yml")
	assert.Nil(t, err)
	io, _, stdout, stderr := iostreams.Test()

	client := func(token, hostname string) (*api.Client, error) {
		return api.TestClient(&http.Client{Transport: fakeHTTP}, token, hostname, false)
	}
	// FIXME: something fishy is occurring here as without making a first call to client function, httpMock does not work
	_, _ = client("", "gitlab.com")

	opts := &StatusOpts{
		Config: func() (config.Config, error) {
			return configs, nil
		},
		HttpClientOverride: client,
		IO:                 io,
	}

	err = statusRun(opts)
	assert.Equal(t, err, nil)
	assert.Equal(t, stdout.String(), "")
	assert.Equal(t, expectedOutput, stderr.String())
}

func Test_statusRun_noInstance(t *testing.T) {
	defer config.StubConfig(`---
git_protocol: ssh
`, "")()

	configs, err := config.ParseConfig("config.yml")
	assert.Nil(t, err)
	io, _, stdout, stderr := iostreams.Test()

	opts := &StatusOpts{
		Config: func() (config.Config, error) {
			return configs, nil
		},
		IO: io,
	}
	t.Run("no instance authenticated", func(t *testing.T) {

		err := statusRun(opts)
		assert.Equal(t, err, cmdutils.SilentError)
		assert.Equal(t, stdout.String(), "")
		assert.Equal(t, "No GitLab instance has been authenticated with glab. Run `glab auth login` to authenticate.\n", stderr.String())
	})
}
