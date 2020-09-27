package clone

import (
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"regexp"
	"testing"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m)
}

func Test_repoClone(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t)
	// profclems/test is a forked repo from glab-cli/test
	cmd := exec.Command(cmdtest.GlabBinaryPath, "repo", "clone", "test")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	assert.NotEmpty(t, out)

	assert.Contains(t, out, "Cloning into 'test'...")
	assert.Contains(t, out, "Updating upstream")
	assert.Regexp(t, regexp.MustCompile(` \* \[new branch\]\s+master\s+-> upstream/master`), out)
}
