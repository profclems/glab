package create

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/profclems/glab/commands/cmdtest"
	"github.com/profclems/glab/pkg/api"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m)
}

func Test_projectCreateCmd(t *testing.T) {
	t.Parallel()
	repo := cmdtest.CopyTestRepo(t)
	expectedPath := fmt.Sprintf("glab-cli/%s", filepath.Base(repo))

	// remove the .git/config so no remotes exist
	err := os.Remove(filepath.Join(repo, ".git/config"))
	if err != nil {
		t.Errorf("could not remove .git/config: %v", err)
	}
	_, err = api.DeleteProject(nil, expectedPath)
	if err != nil {
		t.Logf("unable to delete project %s: %v", expectedPath, err)
	}
	t.Run("create", func(t *testing.T) {
		cmd := exec.Command(cmdtest.GlabBinaryPath, "repo", "create", "-g", "glab-cli", "--public")
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
		cmdtest.Eq(t, string(remote), "git@gitlab.com:"+expectedPath+".git\n")
	})
	p, err := api.GetProject(nil, expectedPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to find project for cleanup"))
	}
	_, err = api.DeleteProject(nil, p.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to delete project during cleanup"))
	}
}

func Test_projectCreateCmdWithArgs(t *testing.T) {
	repo := cmdtest.CopyTestRepo(t)
	expectedPath := "glab-cli/unittest"

	// remove the .git/config so no remotes exist
	err := os.Remove(filepath.Join(repo, ".git/config"))
	if err != nil {
		t.Errorf("could not remove .git/config: %v", err)
	}
	_, err = api.DeleteProject(nil, expectedPath)
	if err != nil {
		t.Logf("unable to delete project %s: %v", expectedPath, err)
	}
	t.Run("create_with_args", func(t *testing.T) {
		cmd := exec.Command(cmdtest.GlabBinaryPath, "repo", "create", expectedPath, "--public")
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
	p, err := api.GetProject(nil, expectedPath)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to find project for cleanup"))
	}
	_, err = api.DeleteProject(nil, p.ID)
	if err != nil {
		t.Fatal(errors.Wrap(err, "failed to delete project during cleanup"))
	}
}
