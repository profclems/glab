package create

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m)
}

func TestMrCmd(t *testing.T) {
	defer config.StubConfig(`---
hosts:
  gitlab.com:
    username: monalisa
    token: OTOKEN
`, "")()
	stubFactory, _ := cmdtest.StubFactoryWithConfig("")
	oldCreateMR := api.CreateMR
	timer, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	api.CreateMR = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreateMergeRequestOptions) (*gitlab.MergeRequest, error) {
		if projectID == "" || projectID == "WRONG_REPO" || projectID == "expected_err" {
			return nil, fmt.Errorf("error expected")
		}
		repo, err := stubFactory.BaseRepo()
		if err != nil {
			return nil, err
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
			WebURL:    "https://" + repo.RepoHost() + "/" + repo.FullName() + "/-/merge_requests/1",
			CreatedAt: &timer,
		}, nil
	}

	cmd := NewCmdCreate(stubFactory)
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
