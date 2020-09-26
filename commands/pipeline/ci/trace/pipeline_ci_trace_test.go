package trace

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ciTrace(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = repo
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "checkout", "origin/test-ci")
	cmd.Dir = repo
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	cmd = exec.Command("git", "checkout", "-b", "test-ci")
	cmd.Dir = repo
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Log(string(b))
		t.Fatal(err)
	}

	tests := []struct {
		desc           string
		args           []string
		assertContains func(t *testing.T, out string)
	}{
		{
			desc: "Has no arg",
			args: []string{},
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...")
				assert.Contains(t, out, "Showing logs for build1 job #732481769")
				assert.Contains(t, out, "Checking out 6caeb21d as test-ci...")
				assert.Contains(t, out, "Do your build here")
				assert.Contains(t, out, "$ echo \"Let's do some cleanup\"")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "Has arg with job-id",
			args: []string{"732481782"},
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...")
				assert.Contains(t, out, "Job succeeded")
			},
		},
		{
			desc: "On a specified repo with job ID",
			args: []string{"-Rglab-cli/test", "716449943"},
			assertContains: func(t *testing.T, out string) {
				assert.Contains(t, out, "Getting job trace...")
				assert.Contains(t, out, "Job succeeded")
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			cmd = exec.Command(glabBinaryPath, append([]string{"pipe", "ci", "trace"}, test.args...)...)
			cmd.Dir = repo

			b, err := cmd.CombinedOutput()
			if err != nil {
				t.Log(string(b))
				t.Fatal(err)
			}
			out := string(b)
			test.assertContains(t, out)
		})
	}

}
