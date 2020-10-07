package clone

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/profclems/glab/commands/cmdtest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "repo_clone_test")
}

func Test_repoClone(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t, "repo_clone_test")
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
