package pipeline

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestPipelineCmd(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "pipeline")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Use \"glab pipeline [command] --help\" for more information about a command.")
}
