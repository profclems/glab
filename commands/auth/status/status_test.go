package status

import (
	"bytes"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"

	"github.com/profclems/glab/commands/cmdutils"

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
  somehost.com:
    token: yyyyyyyyyyyyyyyyyyyy
    git_protocol: ssh
    api_protocol: https
  another.host:
    token: isinvalid
`, "")()

	configs, err := config.ParseConfig("config.yml")
	assert.Nil(t, err)

	io, _, stdout, stderr := utils.IOTest()

	tests := []struct {
		name    string
		opts    *StatusOpts
		wantErr bool
		stderr  string
	}{
		{
			name: "hostname set",
			opts: &StatusOpts{
				Hostname: "gitlab.alpinelinux.org",
			},
			stderr: `gitlab.alpinelinux.org
  x gitlab.alpinelinux.org: api call failed: GET https://gitlab.alpinelinux.org/api/v4/user: 401 {message: 401 Unauthorized}
  ✓ Git operations for gitlab.alpinelinux.org configured to use ssh protocol.
  ✓ API calls for gitlab.alpinelinux.org are made over https protocol
  ✓ REST API Endpoint: https://gitlab.alpinelinux.org/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.alpinelinux.org/api/graphql/
  ✓ Token: ********************
`,
		},
	}
	for _, tt := range tests {
		tt.opts.Config = func() (config.Config, error) {
			return configs, nil
		}
		tt.opts.IO = io
		t.Run(tt.name, func(t *testing.T) {
			if err := statusRun(tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("statusRun() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, stdout.String(), "")
			assert.Equal(t, tt.stderr, stderr.String())
		})
	}
}

func Test_statusRun_noinstance(t *testing.T) {
	defer config.StubConfig(`---
git_protocol: ssh
`, "")()

	configs, err := config.ParseConfig("config.yml")
	assert.Nil(t, err)
	io, _, stdout, stderr := utils.IOTest()

	opts := &StatusOpts{
		Config: func() (config.Config, error) {
			return configs, nil
		},
		IO: io,
	}
	t.Run("no instance authenticated", func(t *testing.T) {
		if err := statusRun(opts); (err != nil) != true {
			t.Errorf("statusRun() error = %v, wantErr %v", err, true)
		}
		assert.Equal(t, stdout.String(), "")
		assert.Equal(t, "No GitLab instance has been authenticated with glab. Run glab auth login to authenticate.\n", stderr.String())
	})
}
