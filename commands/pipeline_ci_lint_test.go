package commands

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pipelineCILint(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "pipeline", "ci", "lint")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	require.Contains(t, string(b), "CI yml is Valid!")
}
