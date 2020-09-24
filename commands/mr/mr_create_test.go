package mr

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

// MRCreate is tested in mr_test

func TestMrCmdWithArgs(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "mr", "create", "someargs")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "accepts 0 arg(s), received 1")
}
