package view

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/pkg/api"
	mainTest "github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
	"os/exec"
	"testing"
	"time"
)

var (
	stubFactory *cmdutils.Factory
)

type author struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func TestMain(m *testing.M) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()

	stubFactory, _ = cmdtest.StubFactoryWithConfig("")

	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.GetIssue = func(client *gitlab.Client, projectID interface{}, issueID int) (*gitlab.Issue, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		repo, err := stubFactory.BaseRepo()
		if err != nil {
			return nil, err
		}
		return &gitlab.Issue{
			ID:          issueID,
			IID:         issueID,
			Title:       "issueTitle",
			Labels:      gitlab.Labels{"test", "bug"},
			State:       "opened",
			Description: "issueBody",
			References: &gitlab.IssueReferences{
				Full: fmt.Sprintf("%s#%d", repo.FullName(), issueID),
			},
			Author: &gitlab.IssueAuthor{
				ID:       issueID,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    fmt.Sprintf("https://%s/%s/-/issues/%d", repo.RepoHost(), repo.FullName(), issueID),
			CreatedAt: &timer,
		}, nil
	}
	cmdtest.InitTest(m, "mr_view_test")
}

func TestNewCmdView_web_numberArg(t *testing.T) {
	cmd := NewCmdView(stubFactory)
	cmd.Flags().StringP("repo", "R", "", "")

	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &mainTest.OutputStub{}
	})
	defer restoreCmd()

	output, err := cmdtest.RunCommand(cmd, "225 -w -R glab-cli/test")
	if err != nil {
		t.Error(err)
		return
	}

	out := stripansi.Strip(output.String())
	outErr := stripansi.Strip(output.Stderr())

	assert.Contains(t, out, "Opening gitlab.com/glab-cli/test/-/issues/225 in your browser.")
	assert.Equal(t, "", outErr)

	if seenCmd == nil {
		t.Log("expected a command to run")
	}
}

func TestNewCmdView(t *testing.T) {
	oldListIssueNotes := api.ListIssueNotes
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.ListIssueNotes = func(client *gitlab.Client, projectID interface{}, issueID int, opts *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
		if projectID == "PROJECT_MR_WITH_EMPTY_NOTE" {
			return []*gitlab.Note{}, nil
		}
		return []*gitlab.Note{
			{
				ID:    1,
				Body:  "Note Body",
				Title: "Note Title",
				Author: author{
					ID:       1,
					Username: "johnwick",
					Name:     "John Wick",
				},
				System:     false,
				CreatedAt:  &timer,
				NoteableID: 0,
			},
			{
				ID:    1,
				Body:  "Marked issue as stale",
				Title: "",
				Author: author{
					ID:       1,
					Username: "johnwick",
					Name:     "John Wick",
				},
				System:     true,
				CreatedAt:  &timer,
				NoteableID: 0,
			},
		}, nil
	}
	cmd := NewCmdView(stubFactory)
	cmd.Flags().StringP("repo", "R", "", "")

	t.Run("show", func(t *testing.T) {
		output, err := cmdtest.RunCommand(cmd, "13 -c -s -R glab-cli/test")
		if err != nil {
			t.Error(err)
			return
		}

		out := stripansi.Strip(output.String())
		outErr := stripansi.Strip(output.Stderr())

		require.Contains(t, out, "issueTitle #13")
		require.Contains(t, out, "issueBody")
		require.Equal(t, outErr, "")
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/issues/13")
		assert.Contains(t, out, "johnwick:\tMarked issue as stale")
	})
	api.ListIssueNotes = oldListIssueNotes
}
