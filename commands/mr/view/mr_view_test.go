package view

import (
	"os/exec"
	"testing"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/internal/run"
	mainTest "github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
)

// TODO: test by mocking the appropriate api function
func TestMain(m *testing.M) {
	cmdtest.InitTest(m)
}

func TestMRView_web_numberArg(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t)
	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &mainTest.OutputStub{}
	})
	defer restoreCmd()

	cmd := exec.Command(cmdtest.GlabBinaryPath, "mr", "view", "-w", "225")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Errorf("error running command `mr view`: %v", err)
	}

	assert.Contains(t, string(b), "Opening gitlab.com/glab-cli/test/-/merge_requests/225 in your browser.")

	if seenCmd == nil {
		t.Log("expected a command to run")
	}
}
