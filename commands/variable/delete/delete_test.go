package delete

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/xanzy/go-gitlab"

	"github.com/alecthomas/assert"
	"github.com/google/shlex"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/httpmock"
)

func Test_NewCmdSet(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    DeleteOpts
		stdinTTY bool
		wantsErr bool
	}{
		{
			name:     "delete var",
			cli:      "cool_secret",
			wantsErr: false,
		},
		{
			name:     "delete scoped var",
			cli:      "cool_secret --scope prod",
			wantsErr: false,
		},
		{
			name:     "delete group var",
			cli:      "cool_secret -g mygroup",
			wantsErr: false,
		},
		{
			name:     "delete scoped group var",
			cli:      "cool_secret -g mygroup --scope prod",
			wantsErr: true,
		},
		{
			name:     "no name",
			cli:      "",
			wantsErr: true,
		},
		{
			name:     "invalid characters in name",
			cli:      "BAD-SECRET",
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

			cmd := NewCmdSet(f, func(opts *DeleteOpts) error {
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
		})
	}
}

func Test_deleteRun_project(t *testing.T) {
	reg := &httpmock.Mocker{
		MatchURL: httpmock.PathAndQuerystring,
	}
	defer reg.Verify(t)

	reg.RegisterResponder("DELETE", "projects/owner%2Frepo/variables/TEST_VAR?filter%5Benvironment_scope%5D=%2A",
		httpmock.NewStringResponse(204, ""),
	)

	io, _, stdout, _ := iostreams.Test()

	opts := &DeleteOpts{
		HTTPClient: func() (*gitlab.Client, error) {
			a, _ := api.TestClient(&http.Client{Transport: reg}, "", "gitlab.com", false)
			return a.Lab(), nil
		},
		BaseRepo: func() (glrepo.Interface, error) {
			return glrepo.FromFullName("owner/repo")
		},
		IO:    io,
		Key:   "TEST_VAR",
		Scope: "*",
	}
	_, _ = opts.HTTPClient()

	err := deleteRun(opts)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "✓ Deleted variable TEST_VAR with scope * for owner/repo\n")
}
