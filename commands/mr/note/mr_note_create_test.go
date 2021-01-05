package note

import (
	"fmt"
	"testing"
	"time"

	"github.com/profclems/glab/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "mr_note_create_test")
}

type author struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func Test_mrNoteCreate(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()

	io, _, stdout, stderr := utils.IOTest()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.CreateMRNote = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.CreateMergeRequestNoteOptions) (*gitlab.Note, error) {
		if projectID == "PROJECT_MR_WITH_EMPTY_NOTE" {
			return &gitlab.Note{}, nil
		}
		return &gitlab.Note{
			ID:    1,
			Body:  *opts.Body,
			Title: *opts.Body,
			Author: author{
				ID:       1,
				Username: "johnwick",
				Name:     "John Wick",
			},
			System:     false,
			CreatedAt:  &timer,
			NoteableID: 0,
		}, nil
	}
	api.GetMR = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.GetMergeRequestsOptions) (*gitlab.MergeRequest, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		repo, err := stubFactory.BaseRepo()
		if err != nil {
			return nil, err
		}
		return &gitlab.MergeRequest{
			ID:          mrID,
			IID:         mrID,
			Title:       "mrTitle",
			Labels:      gitlab.Labels{"test", "bug"},
			State:       "opened",
			Description: "mrBody",
			Author: &gitlab.BasicUser{
				ID:       mrID,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    fmt.Sprintf("https://%s/%s/-/merge_requests/%d", repo.RepoHost(), repo.FullName(), mrID),
			CreatedAt: &timer,
		}, nil
	}
	cmd := NewCmdNote(stubFactory)
	cmdutils.EnableRepoOverride(cmd, stubFactory)

	tests := []struct {
		name          string
		args          string
		want          bool
		assertionFunc func(*testing.T, string, string, error)
	}{
		{
			name: "Has -m flag",
			args: `223 -m "Some test note"`,
			assertionFunc: func(t *testing.T, out, outErr string, err error) {
				require.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/223#note_1")
			},
		},
		{
			name: "Has no flag",
			args: "11",
			assertionFunc: func(t *testing.T, out, outErr string, err error) {
				// TODO: better test survey package
				//require.Equal(t, "aborted... Note has an empty message", err.Error())
			},
		},
		{
			name: "With --repo flag",
			args: `225 -m "Some test note" -R profclems/test`,
			assertionFunc: func(t *testing.T, out, outErr string, err error) {
				require.Contains(t, out, "https://gitlab.com/profclems/test/-/merge_requests/225#note_1")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmdtest.RunCommand(cmd, tt.args)
			if err != nil {
				t.Error(err)
				return
			}

			out := stripansi.Strip(stdout.String())
			outErr := stripansi.Strip(stderr.String())

			tt.assertionFunc(t, out, outErr, err)
		})
	}
}
