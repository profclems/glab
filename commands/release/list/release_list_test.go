package list

import (
	"errors"
	"fmt"
	"testing"
	"time"

	cmdTestUtils "github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"

	"github.com/acarl005/stripansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

type author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func TestNewCmdReleaseList(t *testing.T) {

	oldGetRelease := api.GetRelease
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.GetRelease = func(client *gitlab.Client, projectID interface{}, tag string) (*gitlab.Release, error) {
		if projectID == "" || projectID == "WRONG_REPO" {
			return nil, fmt.Errorf("error expected")
		}
		return &gitlab.Release{
			TagName:     tag,
			Name:        tag,
			Description: "Dummy description for " + tag,
			Author: author{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			CreatedAt: &timer,
		}, nil
	}

	oldListReleases := api.ListReleases
	api.ListReleases = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListReleasesOptions) ([]*gitlab.Release, error) {
		if projectID == "" || projectID == "WRONG_REPO" {
			return nil, errors.New("fatal: wrong Repository")
		}
		return append([]*gitlab.Release{}, &gitlab.Release{
			TagName:     "0.1.0",
			Name:        "Initial Release",
			Description: "Dummy description for 0.1.0",
			Author: author{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			CreatedAt: &timer,
		}), nil
	}

	tests := []struct {
		name       string
		repo       string
		tag        string
		stdOutFunc func(t *testing.T, out string)
		stdErr     string
		wantErr    bool
	}{
		{
			name:    "releases list on test repo",
			wantErr: false,
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "Showing releases 1 of 1 on glab-cli/test")
			},
		},
		{
			name:    "get release by tag on test repo",
			wantErr: false,
			tag:     "v0.0.1-beta",
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "Dummy description for v0.0.1-beta")
			},
		},
		{
			name:    "releases list on custom repo",
			wantErr: false,
			repo:    "profclems/glab",
			stdOutFunc: func(t *testing.T, out string) {
				assert.Contains(t, out, "Showing releases 1 of 1 on profclems/glab")
			},
		},
		{
			name:    "ERR - wrong repo",
			wantErr: true,
			repo:    "WRONG_REPO",
		},
		{
			name:    "ERR - wrong repo with tag",
			wantErr: true,
			repo:    "WRONG_REPO",
			tag:     "v0.0.1-beta",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdReleaseList(cmdTestUtils.StubFactory())
			if tt.repo != "" {
				cmd.Flags().StringP("repo", "R", "", "")
				assert.Nil(t, cmd.Flags().Set("repo", tt.repo))
			}
			if tt.tag != "" {
				assert.Nil(t, cmd.Flags().Set("tag", tt.tag))
			}
			output, err := cmdTestUtils.RunCommand(cmd, ``)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.Nil(t, err)
			}

			out := stripansi.Strip(output.String())
			outErr := stripansi.Strip(output.Stderr())

			tt.stdOutFunc(t, out)
			assert.Contains(t, outErr, tt.stdErr)
		})
	}

	api.GetRelease = oldGetRelease
	api.ListReleases = oldListReleases
}
