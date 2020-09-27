package lint

import (
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/require"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m)
}

func Test_pipelineCILint(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t)
	cmd := exec.Command(cmdtest.GlabBinaryPath, "pipeline", "ci", "lint")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	require.Contains(t, string(b), "CI yml is Valid!")
}
