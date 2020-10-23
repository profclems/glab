package login

import (
	"bytes"
	"testing"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/internal/utils"

	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "auth_login_test")
}

func Test_NewCmdLogin(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		stdin    string
		wants    LoginOptions
		stdinTTY bool
		wantsErr bool
	}{
		{
			name:  "nontty, stdin",
			stdin: "abc123\n",
			cli:   "--stdin",
			wants: LoginOptions{
				Hostname: "gitlab.com",
				Token:    "abc123",
			},
		},
		{
			name:  "tty, stdin",
			stdin: "def456",
			cli:   "--stdin",
			wants: LoginOptions{
				Hostname: "gitlab.com",
				Token:    "def456",
			},
			stdinTTY: true,
		},
		{
			name:     "nontty, hostname",
			cli:      "--hostname salsa.debian.org",
			wantsErr: true,
			stdinTTY: false,
		},
		{
			name:     "nontty",
			cli:      "",
			wantsErr: true,
			stdinTTY: false,
		},
		{
			name:  "nontty, stdin, hostname",
			cli:   "--hostname db.org --stdin",
			stdin: "abc123\n",
			wants: LoginOptions{
				Hostname: "db.org",
				Token:    "abc123",
			},
		},
		{
			name:  "tty, stdin, hostname",
			stdin: "gli789",
			cli:   "--stdin --hostname gl.io",
			wants: LoginOptions{
				Hostname: "gl.io",
				Token:    "gli789",
			},
			stdinTTY: true,
		},
		// TODO: how to test survey
		//{
		//	name:     "tty, hostname",
		//	cli:      "--hostname local.dev",
		//	wants: LoginOptions{
		//		Hostname:    "local.dev",
		//		Token:       "",
		//		Interactive: true,
		//	},
		//	stdinTTY: true,
		//},
		//{
		//	name:     "tty",
		//	cli:      "",
		//	wants: LoginOptions{
		//		Hostname:    "",
		//		Token:       "",
		//		Interactive: true,
		//	},
		//	stdinTTY: true,
		//},
		{
			name:     "token and stdin",
			cli:      "--token xxxx --stdin",
			wantsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, stdin, _, _ := utils.IOTest()
			f := cmdtest.StubFactory("https://gitlab.com/glab-cli/test")

			f.IO = io
			io.IsaTTY = true
			io.IsErrTTY = true
			io.IsInTTY = tt.stdinTTY

			if tt.stdin != "" {
				stdin.WriteString(tt.stdin)
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			cmd := NewCmdLogin(f)
			// TODO cobra hack-around
			cmd.Flags().BoolP("help", "x", false, "")

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()

			if tt.wantsErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.wants.Token, opts.Token)
			assert.Equal(t, tt.wants.Hostname, opts.Hostname)
			assert.Equal(t, tt.wants.Interactive, opts.Interactive)
		})
	}
}
