package commands

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCompletion_bash(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "completion", "-s", "bash")
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "complete -o default -F __start_glab glab") {
		t.Errorf("problem in bash completion:\n%s", string(output))
	}
}

func TestCompletion_zsh(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "completion", "-s", "zsh")
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "#compdef _glab glab") {
		t.Errorf("problem in zsh completion:\n%s", string(output))
	}
}

func TestCompletion_fish(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "completion", "-s", "fish")
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "complete -c glab ") {
		t.Errorf("problem in fish completion:\n%s", string(output))
	}
}

func TestCompletion_powerShell(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "completion", "-s", "powershell")
	cmd.Dir = repo
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(output), "Register-ArgumentCompleter") {
		t.Errorf("problem in powershell completion:\n%s", string(output))
	}
}

func TestCompletion_unsupported(t *testing.T) {
	repo := copyTestRepo(t)
	cmd := exec.Command(glabBinaryPath, "completion", "-s", "csh")
	cmd.Dir = repo
	b, _ := cmd.CombinedOutput()
	out := firstLine(b)

	if out != `unsupported shell type "csh"` {
		t.Fatalf("error: %v", out)
	}
}
