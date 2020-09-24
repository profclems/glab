package list

import (
	"bytes"
	"github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCmdReleaseList(t *testing.T) {

	tests := []struct {
		name string
		repo string
		tag string
		stdOutFunc func(t *testing.T, out string)
		stdErr string
		wantErr bool
	}{
		{
			name: "releases list on test repo",
			wantErr: false,
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "Showing releases")
				assert.Contains(t, out, "on glab-cli/test")
			},
		},
		{
			name: "get release by tag on test repo",
			wantErr: false,
			tag: "v0.0.1-beta",
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "5d3de07d - v0.0.1-beta")
			},
		},
		{
			name: "releases list on custom repo",
			wantErr: false,
			repo: "profclems/glab",
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "Showing releases")
				assert.Contains(t, out, "on profclems/glab")
			},
		},
		{
			name: "ERR - wrong repo",
			wantErr: true,
			repo: "profclems/gla",
		},
		{
			name: "ERR - wrong repo with tag",
			wantErr: true,
			repo: "profclems/gla",
			tag: "v0.0.1-beta",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			var stdout bytes.Buffer
			cmd := NewCmdReleaseList(test.StubFactory())
			if tt.repo != "" {
				cmd.Flags().StringP("repo", "R", "", "")
				assert.Nil(t, cmd.Flags().Set("repo", tt.repo))
			}
			if tt.tag != "" {
				assert.Nil(t, cmd.Flags().Set("tag", tt.tag))
			}
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			_, err := cmd.ExecuteC()
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			tt.stdOutFunc(t, stdout.String())
			assert.Contains(t, stderr.String(), tt.stdErr)
		})
	}
}
