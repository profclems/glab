package create

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_IssueCreate(t *testing.T) {
	cmdtest.CopyTestRepo(t, "issue_create")
	ask, teardown := prompt.InitAskStubber()
	defer teardown()

	ask.Stub([]*prompt.QuestionStub{
		{
			Name:  "confirmation",
			Value: 0,
		},
	})

	oldCreateIssue := api.CreateIssue
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.CreateIssue = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateIssueOptions) (*gitlab.Issue, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		return &gitlab.Issue{
			ID:          1,
			IID:         1,
			Title:       *opts.Title,
			Labels:      opts.Labels,
			State:       "opened",
			Description: *opts.Description,
			Weight:      *opts.Weight,
			Author: &gitlab.IssueAuthor{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    "https://gitlab.com/glab-cli/test/-/issues/1",
			CreatedAt: &timer,
		}, nil
	}

	io, _, stdout, stderr := utils.IOTest()
	f := cmdtest.StubFactory("https://gitlab.com/glab-cli/test")
	f.IO = io
	f.IO.IsaTTY = true
	f.IO.IsErrTTY = true

	cmd := NewCmdCreate(f)
	cmd.Flags().StringP("repo", "R", "", "")

	cliStr := []string{"-t", "myissuetitle",
		"-d", "myissuebody",
		"-l", "test,bug",
		"--weight", "1",
		"--milestone", "1",
		"--linked-mr", "3",
		"--confidential",
		"--assignee", "testuser",
		"-R", "glab-cli/test",
	}

	cli := strings.Join(cliStr, " ")
	t.Log(cli)
	_, err := cmdtest.RunCommand(cmd, cli)
	if err != nil {
		t.Error(err)
	}

	out := stripansi.Strip(stdout.String())
	outErr := stripansi.Strip(stderr.String())
	expectedOut := fmt.Sprintf("#1 myissuetitle (%s)", utils.TimeToPrettyTimeAgo(timer))
	cmdtest.Eq(t, cmdtest.FirstLine([]byte(out)), expectedOut)
	cmdtest.Eq(t, outErr, "\n- Creating issue in glab-cli/test\n")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/issues/1")

	api.CreateIssue = oldCreateIssue
}
