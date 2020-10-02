package create

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var (
	stubFactory *cmdutils.Factory
	cmd         *cobra.Command
)

func TestMain(m *testing.M) {
	defer config.StubConfig(`---
git_protocol: https
hosts:
  gitlab.com:
    username: monalisa
`, "")()

	stubFactory, _ = cmdtest.StubFactoryWithConfig("")

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

	cmd = NewCmdCreate(stubFactory)
	cmd.Flags().StringP("repo", "R", "", "")

	cmdtest.InitTest(m, "mr_cmd_autofill")
}

func TestMrCmd(t *testing.T) {

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
		return
	}

	out := stripansi.Strip(output.String())
	outErr := stripansi.Strip(output.Stderr())

	assert.Contains(t, cmdtest.FirstLine([]byte(out)), `!1 myMRtitle`)
	cmdtest.Eq(t, outErr, "")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/1")
}

func TestNewCmdCreate_autofill(t *testing.T) {
	t.Run("create_autofill", func(t *testing.T) {
		git := exec.Command("git", "checkout", "test-cli")
		b, err := git.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		output, err := cmdtest.RunCommand(cmd, "-f -b master")
		if err != nil {
			t.Error(err)
			return
		}

		out := stripansi.Strip(output.String())
		outErr := stripansi.Strip(output.Stderr())

		assert.Contains(t, cmdtest.FirstLine([]byte(out)), `!1 SOme changes happened'`)
		cmdtest.Eq(t, outErr, "")
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/1")

	})
}
