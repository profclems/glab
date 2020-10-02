package note

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/xanzy/go-gitlab"
	"testing"
	"time"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/require"
)

// TODO: test by mocking the appropriate api function
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

	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
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
	cmd.Flags().StringP("repo", "R", "", "")

	tests := []struct {
		name          string
		args          string
		want          bool
		assertionFunc func(*testing.T, string, string)
	}{
		{
			name: "Has -m flag",
			args: "223 -m \"Some test note\"",
			assertionFunc: func(t *testing.T, out, outErr string) {
				require.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/223#note_1")
			},
		},
		/*
		{
			name: "Has no flag",
			args: "11",
			assertionFunc: func(t *testing.T, out, outErr string) {
				require.Contains(t, out, "aborted... Note has an empty message")
			},
		},
		 */
		{
			name: "With --repo flag",
			args: "225 -m \"Some test note\" -R profclems/test",
			assertionFunc: func(t *testing.T, out, outErr string) {
				require.Contains(t, out, "https://gitlab.com/profclems/test/-/merge_requests/225#note_1")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := cmdtest.RunCommand(cmd, tt.args)
			if err != nil {
				t.Error(err)
				return
			}

			out := stripansi.Strip(output.String())
			outErr := stripansi.Strip(output.Stderr())

			tt.assertionFunc(t, out, outErr)
		})
	}
}
