package create

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/profclems/glab/pkg/prompt"

	"github.com/google/shlex"
	"github.com/profclems/glab/test"

	"github.com/profclems/glab/internal/utils"

	"github.com/acarl005/stripansi"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

var (
	stubFactory *cmdutils.Factory
	cmd         *cobra.Command
	stdout      *bytes.Buffer
	stderr      *bytes.Buffer
)

func runCommand(cmd *cobra.Command, cli string) (*test.CmdOut, error) {
	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)
	_, err = cmd.ExecuteC()

	return &test.CmdOut{
		OutBuf: stdout,
		ErrBuf: stderr,
	}, err
}

func TestMain(m *testing.M) {
	defer config.StubConfig(`---
git_protocol: https
hosts:
  gitlab.com:
    username: monalisa
`, "")()

	var io *utils.IOStreams
	io, _, stdout, stderr = utils.IOTest()
	stubFactory, _ = cmdtest.StubFactoryWithConfig("https://gitlab.com/glab-cli/test.git")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

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
			ID:           1,
			IID:          1,
			Title:        *opts.Title,
			Labels:       opts.Labels,
			State:        "opened",
			Description:  *opts.Description,
			SourceBranch: *opts.SourceBranch,
			TargetBranch: *opts.TargetBranch,
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
	ask, teardown := prompt.InitAskStubber()
	defer teardown()

	ask.Stub([]*prompt.QuestionStub{
		{
			Name:  "confirmation",
			Value: 0,
		},
	})

	cliStr := []string{"-t", "myMRtitle",
		"-d", "myMRbody",
		"-l", "test,bug",
		"--milestone", "1",
		"--assignee", "testuser",
		"-R", "glab-cli/test",
		"-s", "test-cli",
	}

	cli := strings.Join(cliStr, " ")
	t.Log(cli)
	_, err := runCommand(cmd, cli)
	if err != nil {
		t.Error(err)
		return
	}

	out := stripansi.Strip(stdout.String())
	outErr := stripansi.Strip(stderr.String())
	stdout.Reset()
	stderr.Reset()

	assert.Contains(t, cmdtest.FirstLine([]byte(out)), `!1 myMRtitle (test-cli)`)
	assert.Contains(t, outErr, "\nCreating merge request for test-cli into master in glab-cli/test\n\n")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests/1")
	stdout.Reset()
	stderr.Reset()
}

func TestNewCmdCreate_autofill(t *testing.T) {
	ask, teardown := prompt.InitAskStubber()
	defer teardown()

	ask.Stub([]*prompt.QuestionStub{
		{
			Name:  "confirmation",
			Value: 0,
		},
	})

	t.Run("create_autofill", func(t *testing.T) {
		testRepo := cmdtest.CopyTestRepo(t, "mr_cmd_autofill")
		gitCmd := exec.Command("git", "checkout", "mr-autofill-test-br")
		gitCmd.Dir = testRepo
		b, err := gitCmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}
		_, err = runCommand(cmd, "-f -b master -s mr-autofill-test-br")
		if err != nil {
			t.Error(err)
			return
		}

		out := stripansi.Strip(stdout.String())
		outErr := stripansi.Strip(stderr.String())

		assert.Equal(t, `!1 docs: add some changes to txt file (mr-autofill-test-br)
 https://gitlab.com/glab-cli/test/-/merge_requests/1

`, out)
		assert.Contains(t, outErr, "\nCreating merge request for mr-autofill-test-br into master in glab-cli/test\n\n")
		stdout.Reset()
		stderr.Reset()
	})
}

func TestMRCreate_nontty_insufficient_flags(t *testing.T) {
	stubFactory.IO.SetPrompt("true")
	cmd = NewCmdCreate(stubFactory)
	_, err := runCommand(cmd, "")
	if err == nil {
		t.Fatal("expected error")
	}

	assert.Equal(t, "--title or --fill required for non-interactive mode", err.Error())

	assert.Equal(t, "", stdout.String())
}

func TestMrBodyAndTitle(t *testing.T) {
	testRepo := cmdtest.CopyTestRepo(t, "mr_cmd_autofill")
	gitCmd := exec.Command("git", "checkout", "mr-autofill-test-br")
	gitCmd.Dir = testRepo
	b, err := gitCmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	opts := &CreateOpts{
		SourceBranch:         "mr-autofill-test-br",
		TargetBranch:         "master",
		TargetTrackingBranch: "origin/master",
	}
	t.Run("", func(t *testing.T) {
		if err = mrBodyAndTitle(opts); err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		assert.Equal(t, "docs: add some changes to txt file", opts.Title)
		assert.Equal(t, `Here, I am adding some commit body.
Little longer

Resolves #1
`, opts.Description)
	})
}
