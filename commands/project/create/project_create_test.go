package create

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_projectCreateCmd(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)
	expectedPath := fmt.Sprintf("glab-cli/%s", filepath.Base(repo))

	// remove the .git/config so no remotes exist
	err := os.Remove(filepath.Join(repo, ".git/config"))
	if err != nil {
		t.Errorf("could not remove .git/config: %v", err)
	}
	_, err = deleteProject(expectedPath)
	if err != nil {
		t.Logf("unable to delete project %s: %v", expectedPath, err)
	}
	t.Run("create", func(t *testing.T) {
		cmd := exec.Command(glabBinaryPath, "repo", "create", "-g", "glab-cli", "--public")
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		rp := strings.Split(expectedPath, "/")[1]
		require.Contains(t, string(b),
			"✓ Created repository glab / "+rp+" on GitLab: https://gitlab.com/"+expectedPath+"\n✓ Added remote git@gitlab.com:"+expectedPath+".git\n")

		gitCmd := exec.Command("git", "remote", "get-url", "origin")
		gitCmd.Dir = repo
		gitCmd.Stdout = nil
		gitCmd.Stderr = nil
		remote, err := gitCmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		eq(t, string(remote), "git@gitlab.com:"+expectedPath+".git\n")
	})
	p, err := getProject(expectedPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to find project for cleanup"))
	}
	_, err = deleteProject(p.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to delete project during cleanup"))
	}
}

func Test_projectCreateCmdWithArgs(t *testing.T) {
	repo := copyTestRepo(t)
	expectedPath := "glab-cli/unittest"

	// remove the .git/config so no remotes exist
	err := os.Remove(filepath.Join(repo, ".git/config"))
	if err != nil {
		t.Errorf("could not remove .git/config: %v", err)
	}
	_, err = deleteProject(expectedPath)
	if err != nil {
		t.Logf("unable to delete project %s: %v", expectedPath, err)
	}
	t.Run("create_with_args", func(t *testing.T) {
		cmd := exec.Command(glabBinaryPath, "repo", "create", expectedPath, "--public")
		cmd.Dir = repo

		b, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(string(b))
			t.Fatal(err)
		}

		require.Contains(t, string(b),
			"✓ Created repository glab / unittest on GitLab: https://gitlab.com/"+expectedPath+"\n")
		err = initialiseRepo(expectedPath, "git@gitlab.com:"+expectedPath+".git")
		if err != nil {
			t.Fatal(err)
		}
	})
	p, err := getProject(expectedPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to find project for cleanup"))
	}
	_, err = deleteProject(p.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to delete project during cleanup"))
	}
}
