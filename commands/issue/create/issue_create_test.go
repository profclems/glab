package create

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

func Test_IssueCreate(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "issue", "create",
		"-t", fmt.Sprintf("Testing Issue Title %v", rand.Intn(200)),
		"-d", "This issue is created as a test",
		"-l", "test,bug",
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

	assert.Contains(t, out, "Testing Issue Title")
	assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/issues/")
}
