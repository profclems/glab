package trace

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/assert"
)

var (
	stubFactory *cmdutils.Factory
	cmd         *cobra.Command
	stdout      *bytes.Buffer
	stderr      *bytes.Buffer
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "pipeline_ci_trace_test")
}

func Test_ciTrace(t *testing.T) {
	t.Parallel()
	defer config.StubConfig(`---
git_protocol: https
hosts:
  gitlab.com:
    username: root
`, "")()

	var io *utils.IOStreams
	io, _, stdout, stderr = utils.IOTest()
	stubFactory, _ = cmdtest.StubFactoryWithConfig("https://gitlab.com/glab-cli/test.git")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	repo := cmdtest.CopyTestRepo(t, "pipeline_ci_trace_test")
	gitCmd := exec.Command("git", "fetch", "origin")
	gitCmd.Dir = repo
	if b, err := gitCmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	gitCmd = exec.Command("git", "checkout", "origin/test-ci")
	gitCmd.Dir = repo
	if b, err := gitCmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		//t.Fatal(err)
	}

	gitCmd = exec.Command("git", "checkout", "test-ci")
	gitCmd.Dir = repo
	if b, err := gitCmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	tests := []struct {
		desc           string
		args           string
		assertContains func(t *testing.T, out string)
	}{
		// TODO: better test for survey prompt when no argument is provided
		{
			desc: "Has no arg",
			args: ``,
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...")
				assert.Contains(t, out, "Showing logs for ")
				assert.Contains(t, out, "Preparing the \"docker+machine\"")
				assert.Contains(t, out, "Checking out 6caeb21d as test-ci...")
				assert.Contains(t, out, "$ echo \"Let's do some cleanup\"")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "Has arg with job-id",
			args: `716449943`,
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...\n")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "On a specified repo with job ID",
			args: "716449943 -X glab-cli/test",
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...\n")
				assert.Contains(t, out, "Job succeeded")
			},
		},
	}

	cmd = NewCmdTrace(stubFactory)
	cmd.Flags().StringP("repo", "X", "", "")

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := cmdtest.RunCommand(cmd, tt.args)
			if err != nil {
				t.Fatal(err)
			}
			tt.assertContains(t, stdout.String())
		})
	}

}
