package commands

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Alias(t *testing.T)  {
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "alias")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Use \"glab alias [command] --help\" for more information about a command.")
}