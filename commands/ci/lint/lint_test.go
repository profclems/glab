package lint

import (
	"os/exec"
	"testing"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "pipeline_ci_lint_test")
}

func Test_pipelineCILint(t *testing.T) {
	t.Parallel()
	repo := cmdtest.CopyTestRepo(t, "pipeline_ci_lint_test")
	cmd := exec.Command(cmdtest.GlabBinaryPath, "pipeline", "ci", "lint")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	require.Contains(t, string(b), "CI yml is Valid!")
}
