package commands

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

func Test_MrCreate (t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "mr", "create",
		"-t", fmt.Sprintf("Testing MR Title %v", rand.Int()),
		"-d", "This MR is created as a test",
		"-l", "test,bug",
		"--assignee", "profclems",
		"--weight", "1",
		"--milestone", "1",
		"--linked-mr", "3")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Testing MR Title")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge-requests/")
}