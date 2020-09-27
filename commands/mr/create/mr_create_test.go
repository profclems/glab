package create

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
	"time"
)

func TestMrCmd(t *testing.T) {
	oldCreateMR := api.CreateMR
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.CreateMR = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		return &gitlab.MergeRequest{
			ID:          1,
			IID:         1,
			Title:       *opts.Title,
			Labels:      opts.Labels,
			State:       "opened",
			Description: *opts.Description,
			Author: &gitlab.BasicUser{
				ID:       1,
				Name:     "John Dev Wick",
				Username: "jdwick",
			},
			WebURL:    "https://gitlab.com/glab-cli/test/-/merge_requests/1",
			CreatedAt: &timer,
		}, nil
	}

	cmd := NewCmdCreate(cmdtest.StubFactory())
	cmd.Flags().StringP("repo", "R", "", "")

	cliStr := []string{"-t", "myMRtitle",
		"-d", "myMRbody",
		"-l", "test,bug",
		"--milestone", "1",
		"--assignee", "testuser",
		"-R", "glab-cli/test",
	}

	cli := strings.Join(cliStr, " ")
	t.Log(cli)
	output, err := cmdtest.RunCommand(cmd, cli)
	if err != nil {
		t.Error(err)
	}

	out := stripansi.Strip(output.String())
	outErr := stripansi.Strip(output.Stderr())

	assert.Contains(t, cmdtest.FirstLine([]byte(out)), `#1 myMRtitle`)
	cmdtest.Eq(t, outErr, "")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/1")

	api.CreateMR = oldCreateMR
}
