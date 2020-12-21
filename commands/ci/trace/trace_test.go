package trace

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/google/shlex"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"

	"github.com/profclems/glab/commands/cmdtest"
)

var (
	stubFactory *cmdutils.Factory
	cmd         *cobra.Command
	stdout      *bytes.Buffer
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "ci_trace_test")
}

func TestNewCmdTrace(t *testing.T) {
	defer config.StubConfig(`---
git_protocol: https
hosts:
  gitlab.com:
    username: root
`, "")()

	var io *utils.IOStreams
	io, _, stdout, _ = utils.IOTest()
	stubFactory, _ = cmdtest.StubFactoryWithConfig("https://gitlab.com/glab-cli/test.git")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	repo := cmdtest.CopyTestRepo(t, "ci_trace_test")
	gitCmd := exec.Command("git", "fetch", "origin")
	gitCmd.Dir = repo
	if _, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatal(err)
	}

	gitCmd = exec.Command("git", "checkout", "test-cli")
	gitCmd.Dir = repo
	if _, err := gitCmd.CombinedOutput(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		args     string
		wantOpts *TraceOpts
	}{
		{
			name: "Has no arg",
			args: ``,
			wantOpts: &TraceOpts{
				Branch: "test-cli",
				JobID:  0,
			},
		},
		{
			name: "Has arg with job-id",
			args: `224356863`,
			wantOpts: &TraceOpts{
				Branch: "test-cli",
				JobID:  224356863,
			},
		},
		{
			name: "On a specified repo with job ID",
			args: "224356863 -X glab-cli/test",
			wantOpts: &TraceOpts{
				Branch: "test-cli",
				JobID:  224356863,
			},
		},
	}

	var actualOpts *TraceOpts
	cmd = NewCmdTrace(stubFactory, func(opts *TraceOpts) error {
		actualOpts = opts
		return nil
	})
	cmd.Flags().StringP("repo", "X", "", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantOpts.IO = stubFactory.IO

			argv, err := shlex.Split(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			cmd.SetArgs(argv)
			_, err = cmd.ExecuteC()
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.wantOpts.JobID, actualOpts.JobID)
			assert.Equal(t, tt.wantOpts.Branch, actualOpts.Branch)
			assert.Equal(t, tt.wantOpts.Branch, actualOpts.Branch)
			assert.Equal(t, tt.wantOpts.IO, actualOpts.IO)
		})
	}

}

func TestTraceRun(t *testing.T) {
	var io *utils.IOStreams
	io, _, stdout, _ = utils.IOTest()
	stubFactory = cmdtest.StubFactory("https://gitlab.com/glab-cli/test.git")
	stubFactory.IO = io
	stubFactory.IO.IsaTTY = true
	stubFactory.IO.IsErrTTY = true

	tests := []struct {
		desc           string
		args           string
		assertContains func(t *testing.T, out string)
	}{
		{
			desc: "Has no arg",
			args: ``,
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...")
				assert.Contains(t, out, "Showing logs for ")
				assert.Contains(t, out, "Preparing the \"docker+machine\"")
				assert.Contains(t, out, "$ echo \"After script section\"")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "Has arg with job-id",
			args: `886379752`,
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...\n")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "On a specified repo with job ID",
			args: "886379752 -X glab-cli/test",
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...\n")
			},
		},
	}

	cmd = NewCmdTrace(stubFactory, nil)
	cmd.Flags().StringP("repo", "X", "", "")

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.args == "" {
				as, teardown := prompt.InitAskStubber()
				defer teardown()

				as.StubOne("cleanup4 (886379752) - success")
			}
			argv, err := shlex.Split(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			cmd.SetArgs(argv)
			_, err = cmd.ExecuteC()
			if err != nil {
				t.Fatal(err)
			}
			tt.assertContains(t, stdout.String())
			stdout.Reset()
		})
	}

}
