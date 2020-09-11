package commands

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"

	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/test"
)

func TestMRView_web_numberArg(t *testing.T) {
	repo := copyTestRepo(t)
	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &test.OutputStub{}
	})
	defer restoreCmd()

	cmd := exec.Command(glabBinaryPath, "mr", "view", "-w", "225")
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
