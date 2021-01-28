package view

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/profclems/glab/pkg/iostreams"

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

	var io *iostreams.IOStreams
	io, _, stdout, stderr = iostreams.IOTest()
	stubFactory, _ = cmdtest.StubFactoryWithConfig("")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

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
			Milestone: &gitlab.Milestone{
				Title: "MilestoneTitle",
			},
			Assignees: []*gitlab.IssueAssignee{
				{
					Username: "mona",
				},
				{
					Username: "lisa",
				},
			},
			Author: &gitlab.IssueAuthor{
				ID:       issueID,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:         fmt.Sprintf("https://%s/%s/-/issues/%d", repo.RepoHost(), repo.FullName(), issueID),
			CreatedAt:      &timer,
			UserNotesCount: 2,
		}, nil
	}
	cmdtest.InitTest(m, "mr_view_test")
}

func TestNewCmdView_web_numberArg(t *testing.T) {
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

	assert.Contains(t, stderr.String(), "Opening gitlab.com/glab-cli/test/-/issues/225 in your browser.")
	assert.Equal(t, "", stdout.String())

	if seenCmd == nil {
		t.Log("expected a command to run")
	}
	stdout.Reset()
	stderr.Reset()
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

		require.Contains(t, out, "issueTitle #13")
		require.Contains(t, out, "issueBody")
		require.Equal(t, outErr, "")
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/issues/13")
		assert.Contains(t, out, "johnwick Marked issue as stale")
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
			`title:\tissueTitle`,
			`assignees:\tmona, lisa`,
			`author:\tjdwick`,
			`state:\topen`,
			`comments:\t2`,
			`labels:\ttest, bug`,
			`milestone:\tMilestoneTitle\n`,
			`--`,
			`issueBody`,
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
	api.ListIssueNotes = oldListIssueNotes
}
