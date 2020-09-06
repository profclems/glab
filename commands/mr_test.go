package commands

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMrCmd(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	repo := copyTestRepo(t)
	var mrID string
	t.Run("create", func(t *testing.T) {
		git := exec.Command("git", "checkout", "test-cli")
		git.Dir = repo
		b, err := git.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		cmd := exec.Command(glabBinaryPath, "mr", "create",
			"-t", fmt.Sprintf("MR Title %v", rand.Int()),
			"-d", "This MR is created as a test",
			"-l", "test,bug",
			"--assignee", "profclems",
			"--milestone", "1",
		)
		cmd.Dir = repo

		b, _ = cmd.CombinedOutput()
		out := string(b)
		t.Log(out)
		out = stripansi.Strip(out)
		assert.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests")
		r := regexp.MustCompile(`#\S+`)

		//i := strings.Index(out, "/diffs\n")
		//mrID = strings.TrimPrefix(out[:i], "https://gitlab.com/glab-cli/test/-/merge_requests/")
		mrID = strings.TrimPrefix(r.FindStringSubmatch(out)[0], "#")
		t.Log(mrID)
	})
	t.Run("show", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create likely failed")
		}
		cmd := exec.Command(glabBinaryPath, "mr", "show", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		out := string(b)
		outStripped := stripansi.Strip(out) // To remove ansi chars added by glamour
		require.Contains(t, outStripped, "This MR is created as a test")
		assert.Contains(t, out, fmt.Sprintf("https://gitlab.com/glab-cli/test/-/merge_requests/%s", mrID))
	})
	t.Run("delete", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create likely failed")
		}
		cmd := exec.Command(glabBinaryPath, "mr", "delete", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		out := stripansi.Strip(string(b))
		require.Contains(t, out, fmt.Sprintf("Deleting Merge Request #%s\n", mrID))
		require.Contains(t, out, "Merge Request Deleted Successfully")
	})
}

func Test_mrCmd_autofill(t *testing.T) {
	repo := copyTestRepo(t)
	var mrID string
	t.Run("create", func(t *testing.T) {
		git := exec.Command("git", "checkout", "test-cli")
		git.Dir = repo
		b, err := git.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		cmd := exec.Command(glabBinaryPath, "mr", "create", "-f")
		cmd.Dir = repo

		b, _ = cmd.CombinedOutput()
		out := string(b)
		t.Log(out)
		out = stripansi.Strip(out)
		require.Contains(t, out, "https://gitlab.com/glab-cli/test/-/merge_requests")
		r := regexp.MustCompile(`#\S+`)
		mrID = strings.TrimPrefix(r.FindStringSubmatch(out)[0], "#")
		t.Log(mrID)

	})
	t.Run("delete", func(t *testing.T) {
		if mrID == "" {
			t.Skip("mrID is empty, create -F likely failed")
		}
		cmd := exec.Command(glabBinaryPath, "mr", "delete", mrID)
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}
		out := stripansi.Strip(string(b))
		require.Contains(t, out, fmt.Sprintf("Deleting Merge Request #%s\n", mrID))
		require.Contains(t, out, "Merge Request Deleted Successfully")
	})

}

func TestMrCmdNoArgs(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	cmd := exec.Command(glabBinaryPath, "mr")
	cmd.Dir = repo

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}
	out := string(b)
	t.Log(out)

	assert.Contains(t, out, "Use \"glab mr [command] --help\" for more information about a command.")
}
