package set

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/alecthomas/assert"
	"github.com/google/shlex"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/httpmock"
	"github.com/xanzy/go-gitlab"
)

func Test_isValidKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key is empty",
			args: args{
				key: "",
			},
			want: false,
		},
		{
			name: "key is valid",
			args: args{
				key: "abc123_",
			},
			want: true,
		},
		{
			name: "key is invalid",
			args: args{
				key: "abc-123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidKey(tt.args.key); got != tt.want {
				t.Errorf("isValidKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	tests := []struct {
		name     string
		valueArg string
		want     string
		stdin    string
	}{
		{
			name:     "literal value",
			valueArg: "a secret",
			want:     "a secret",
		},
		{
			name:  "from stdin",
			want:  "a secret",
			stdin: "a secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, stdin, _, _ := iostreams.Test()

			io.IsInTTY = false

			_, err := stdin.WriteString(tt.stdin)
			assert.NoError(t, err)

			args := []string{tt.valueArg}
			value, err := getValue(&SetOpts{
				Value: tt.valueArg,
				IO:    io,
			}, args)
			assert.NoError(t, err)

			assert.Equal(t, value, tt.want)
		})
	}
}

func Test_NewCmdSet(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    SetOpts
		stdinTTY bool
		wantsErr bool
	}{
		{
			name:     "invalid type",
			cli:      "cool_secret -g coolGroup -t 'mV'",
			wantsErr: true,
		},
		{
			name:     "value argument and value flag specified",
			cli:      `cool_secret value -v"another value"`,
			wantsErr: true,
		},
		{
			name:     "no name",
			cli:      "",
			wantsErr: true,
		},
		{
			name:     "no value, stdin is terminal",
			cli:      "cool_secret",
			stdinTTY: true,
			wantsErr: true,
		},
		{
			name: "protected var",
			cli:  `cool_secret -v"a secret" -p`,
			wants: SetOpts{
				Key:       "cool_secret",
				Protected: true,
				Value:     "a secret",
				Group:     "",
			},
		},
		{
			name: "protected var in group",
			cli:  `cool_secret --group coolGroup -v"cool"`,
			wants: SetOpts{
				Key:       "cool_secret",
				Protected: false,
				Value:     "cool",
				Group:     "coolGroup",
			},
		},
		{
			name: "leading numbers in name",
			cli:  `123_TOKEN -v"cool"`,
			wants: SetOpts{
				Key:       "123_TOKEN",
				Protected: false,
				Value:     "cool",
				Group:     "",
			},
		},
		{
			name:     "invalid characters in name",
			cli:      `BAD-SECRET -v"cool"`,
			wantsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			f := &cmdutils.Factory{
				IO: io,
			}

			io.IsInTTY = tt.stdinTTY

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *SetOpts
			cmd := NewCmdSet(f, func(opts *SetOpts) error {
				gotOpts = opts
				return nil
			})
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

			assert.Equal(t, tt.wants.Key, gotOpts.Key)
			assert.Equal(t, tt.wants.Value, gotOpts.Value)
			assert.Equal(t, tt.wants.Group, gotOpts.Group)
		})
	}
}

func Test_setRun_project(t *testing.T) {
	reg := &httpmock.Mocker{}
	defer reg.Verify(t)

	reg.RegisterResponder("POST", "/projects/owner/repo/variables",
		httpmock.NewStringResponse(201, `
			{
    			"key": "NEW_VARIABLE",
    			"value": "new value",
    			"variable_type": "env_var",
    			"protected": false,
    			"masked": false
			}
		`),
	)

	io, _, stdout, _ := iostreams.Test()

	opts := &SetOpts{
		HTTPClient: func() (*gitlab.Client, error) {
			a, _ := api.TestClient(&http.Client{Transport: reg}, "", "gitlab.com", false)
			return a.Lab(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.FromFullName("owner/repo")
		},
		IO:    io,
		Key:   "NEW_VARIABLE",
		Value: "new value",
	}
	_, _ = opts.HTTPClient()

	err := setRun(opts)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "✓ Created variable NEW_VARIABLE for owner/repo\n")
}

func Test_setRun_group(t *testing.T) {
	reg := &httpmock.Mocker{}
	defer reg.Verify(t)

	reg.RegisterResponder("POST", "/groups/mygroup/variables",
		httpmock.NewStringResponse(201, `
			{
    			"key": "NEW_VARIABLE",
    			"value": "new value",
    			"variable_type": "env_var",
    			"protected": false,
    			"masked": false
			}
		`),
	)

	io, _, stdout, _ := iostreams.Test()

	opts := &SetOpts{
		HTTPClient: func() (*gitlab.Client, error) {
			a, _ := api.TestClient(&http.Client{Transport: reg}, "", "gitlab.com", false)
			return a.Lab(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.FromFullName("owner/repo")
		},
		IO:    io,
		Key:   "NEW_VARIABLE",
		Value: "new value",
		Group: "mygroup",
	}
	_, _ = opts.HTTPClient()

	err := setRun(opts)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "✓ Created variable NEW_VARIABLE for group mygroup\n")
}
