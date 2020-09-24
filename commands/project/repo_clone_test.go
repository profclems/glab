package project

import (
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_repoClone(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	// profclems/test is a forked repo from glab-cli/test
	cmd := exec.Command(glabBinaryPath, "repo", "clone", "test")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Cloning into 'test'...")
	assert.Contains(t, out, "Updating upstream")
	assert.Regexp(t, regexp.MustCompile(` \* \[new branch\]\s+master\s+-> upstream/master`), out)
}
