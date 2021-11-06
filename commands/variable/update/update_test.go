package update

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

func Test_NewCmdSet(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    UpdateOpts
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
			wants: UpdateOpts{
				Key:       "cool_secret",
				Protected: true,
				Value:     "a secret",
				Group:     "",
			},
		},
		{
			name: "protected var in group",
			cli:  `cool_secret --group coolGroup -v"cool"`,
			wants: UpdateOpts{
				Key:       "cool_secret",
				Protected: false,
				Value:     "cool",
				Group:     "coolGroup",
			},
		},
		{
			name: "leading numbers in name",
			cli:  `123_TOKEN -v"cool"`,
			wants: UpdateOpts{
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

			var gotOpts *UpdateOpts
			cmd := NewCmdSet(f, func(opts *UpdateOpts) error {
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

func Test_updateRun_project(t *testing.T) {
	reg := &httpmock.Mocker{}
	defer reg.Verify(t)

	reg.RegisterResponder("PUT", "/projects/owner/repo/variables/TEST_VARIABLE",
		httpmock.NewStringResponse(201, `
			{
    			"key": "TEST_VARIABLE",
    			"value": "foo",
    			"variable_type": "env_var",
    			"protected": false,
    			"masked": false
			}
		`),
	)

	io, _, stdout, _ := iostreams.Test()

	opts := &UpdateOpts{
		HTTPClient: func() (*gitlab.Client, error) {
			a, _ := api.TestClient(&http.Client{Transport: reg}, "", "gitlab.com", false)
			return a.Lab(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.FromFullName("owner/repo")
		},
		IO:    io,
		Key:   "TEST_VARIABLE",
		Value: "foo",
		Scope: "*",
	}
	_, _ = opts.HTTPClient()

	err := updateRun(opts)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "✓ Updated variable TEST_VARIABLE for project owner/repo with scope *\n")
}

func Test_updateRun_group(t *testing.T) {
	reg := &httpmock.Mocker{}
	defer reg.Verify(t)

	reg.RegisterResponder("PUT", "/groups/mygroup/variables/TEST_VARIABLE",
		httpmock.NewStringResponse(201, `
			{
    			"key": "TEST_VARIABLE",
    			"value": "blargh",
    			"variable_type": "env_var",
    			"protected": false,
    			"masked": false
			}
		`),
	)

	io, _, stdout, _ := iostreams.Test()

	opts := &UpdateOpts{
		HTTPClient: func() (*gitlab.Client, error) {
			a, _ := api.TestClient(&http.Client{Transport: reg}, "", "gitlab.com", false)
			return a.Lab(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.FromFullName("owner/repo")
		},
		IO:    io,
		Key:   "TEST_VARIABLE",
		Value: "blargh",
		Group: "mygroup",
	}
	_, _ = opts.HTTPClient()

	err := updateRun(opts)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "✓ Updated variable TEST_VARIABLE for group mygroup\n")
}
