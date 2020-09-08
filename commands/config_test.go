package commands

import (
	"os/exec"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/stretchr/testify/require"
)

func TestConfigGet(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "config", "get", "gitlab_uri")
	cmd.Dir = repo
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error running command `config get gitlab_uri`: %v", err)
	}
	out := firstLine(b)
	t.Log(out)

	eq(t, out, "https://gitlab.com")
}

func TestConfigGet_not_found(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "config", "get", "missing")
	cmd.Dir = repo
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error running command `config get missing`: %v", err)
	}
	out := firstLine(b)
	t.Log(out)

	eq(t, out, "")
}

func TestConfigSet(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "config", "set", "somekey", "someval")
	cmd.Dir = repo
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error running command `config set somekey someval`: %v", err)
	}
	out := firstLine(b)
	t.Log(out)

	eq(t, out, "")
}

func TestConfigSet_update(t *testing.T) {
	// get older value
	initial := config.GetEnv("TEST_CONFIG")

	// set new value
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "config", "set", "test_config", "changed")
	cmd.Dir = repo
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error running command `config set test_config changed`: %v", err)
	}
	out := firstLine(b)
	t.Log(out)
	if len(out) > 0 {
		t.Errorf("expected output to be blank: %q", out)
	}

	// get new value
	cmd = exec.Command(glabBinaryPath, "config", "get", "test_config")
	cmd.Dir = repo
	b, err = cmd.CombinedOutput()
	if err != nil {
		t.Log(string(b))
		t.Fatalf("error running command `config get gitlab_uri`: %v", err)
	}
	after := firstLine(b)
	t.Log(after)
	if after != "changed" {
		t.Errorf("expected output to be changed, got %q", after)
	}

	require.NotEqual(t, initial, after)
}