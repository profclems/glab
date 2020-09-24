package project

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func Test_Project(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "project")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Use \"glab repo [command] --help\" for more information about a command.")
}
