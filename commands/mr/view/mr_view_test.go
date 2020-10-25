package view

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/profclems/glab/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/internal/run"
	mainTest "github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
)

var (
	stubFactory *cmdutils.Factory
	stdout      *bytes.Buffer
	stderr      *bytes.Buffer
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

	var io *utils.IOStreams
	io, _, stdout, stderr = utils.IOTest()
	stubFactory, _ = cmdtest.StubFactoryWithConfig("")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
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
			Assignees: []*gitlab.BasicUser{
				{
					Username: "mona",
				},
				{
					Username: "lisa",
				},
			},
			WebURL:         fmt.Sprintf("https://%s/%s/-/merge_requests/%d", repo.RepoHost(), repo.FullName(), mrID),
			CreatedAt:      &timer,
			UserNotesCount: 2,
			Milestone: &gitlab.Milestone{
				Title: "MilestoneTitle",
			},
		}, nil
	}
	cmdtest.InitTest(m, "mr_view_test")
}

func TestMRView_web_numberArg(t *testing.T) {
	cmd := NewCmdView(stubFactory)
	cmdutils.EnableRepoOverride(cmd, stubFactory)

	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &mainTest.OutputStub{}
	})
	defer restoreCmd()

	_, err := cmdtest.RunCommand(cmd, "225 -w -R glab-cli/test")
	if err != nil {
		t.Error(err)
		return
	}

	out := stripansi.Strip(stdout.String())
	outErr := stripansi.Strip(stderr.String())
	stdout.Reset()
	stderr.Reset()

	assert.Contains(t, outErr, "Opening gitlab.com/glab-cli/test/-/merge_requests/225 in your browser.")
	assert.Equal(t, out, "")

	if seenCmd == nil {
		t.Log("expected a command to run")
	}
}

func TestMRView(t *testing.T) {
	oldListMrNotes := api.ListMRNotes
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.ListMRNotes = func(client *gitlab.Client, projectID interface{}, mrID int, opts *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error) {
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
				Body:  "Marked PR as ready",
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

	t.Run("show", func(t *testing.T) {
		cmd := NewCmdView(stubFactory)
		cmdutils.EnableRepoOverride(cmd, stubFactory)

		_, err := cmdtest.RunCommand(cmd, "13 -c -s -R glab-cli/test")

		if err != nil {
			t.Error(err)
			return
		}

		out := stripansi.Strip(stdout.String())
		outErr := stripansi.Strip(stderr.String())
		stdout.Reset()
		stderr.Reset()

		require.Contains(t, out, "mrTitle !13")
		require.Equal(t, outErr, "")
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/13")
		assert.Contains(t, out, "johnwick Marked PR as ready")
	})

	t.Run("no_tty", func(t *testing.T) {
		stubFactory.IO.IsaTTY = false
		stubFactory.IO.IsErrTTY = false

		cmd := NewCmdView(stubFactory)
		cmdutils.EnableRepoOverride(cmd, stubFactory)

		_, err := cmdtest.RunCommand(cmd, "13 -c -s -R glab-cli/test")
		if err != nil {
			t.Error(err)
			return
		}

		expectedOutputs := []string{
			`title:\tmrTitle`,
			`assignees:\tmona, lisa`,
			`author:\tjdwick`,
			`state:\topen`,
			`comments:\t2`,
			`labels:\ttest, bug`,
			`milestone:\tMilestoneTitle\n`,
			`--`,
			`mrBody`,
		}

		out := stripansi.Strip(stdout.String())
		outErr := stripansi.Strip(stderr.String())

		cmdtest.Eq(t, outErr, "")
		t.Helper()
		var r *regexp.Regexp
		for _, l := range expectedOutputs {
			r = regexp.MustCompile(l)
			if !r.MatchString(out) {
				t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, out)
				return
			}
		}
	})
	api.ListMRNotes = oldListMrNotes
}
