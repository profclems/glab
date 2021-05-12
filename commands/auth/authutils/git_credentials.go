package authutils

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/google/shlex"
)

type GitCredentialFlow struct {
	Executable string

	shouldSetup bool
	helper      string
}

func (gc *GitCredentialFlow) Prompt(hostname, protocol string) error {
	gc.helper, _ = gitCredentialHelper(hostname, protocol)
	if isOurCredentialHelper(gc.helper) {
		return nil
	}

	err := prompt.AskOne(&survey.Confirm{
		Message: "Authenticate Git with your GitLab credentials?",
		Default: true,
	}, &gc.shouldSetup)
	if err != nil {
		return fmt.Errorf("could not prompt: %w", err)
	}

	return nil
}

func (gc *GitCredentialFlow) ShouldSetup() bool {
	return gc.shouldSetup
}

func (gc *GitCredentialFlow) Setup(hostname, protocol, username, authToken string) error {
	return gc.gitCredentialSetup(hostname, protocol, username, authToken)
}

func (gc *GitCredentialFlow) gitCredentialSetup(hostname, protocol, username, password string) error {
	if gc.helper == "" {
		// first use a blank value to indicate to git we want to sever the chain of credential helpers
		preConfigureCmd := git.GitCommand("config", "--global", gitCredentialHelperKey(hostname, protocol), "")
		if err := run.PrepareCmd(preConfigureCmd).Run(); err != nil {
			return err
		}

		// use glab as a credential helper (for this host only)
		configureCmd := git.GitCommand(
			"config", "--global", "--add",
			gitCredentialHelperKey(hostname, protocol),
			fmt.Sprintf("!%s auth git-credential", shellQuote(gc.Executable)),
		)
		return run.PrepareCmd(configureCmd).Run()
	}

	// clear previous cached credentials
	rejectCmd := git.GitCommand("credential", "reject")

	rejectCmd.Stdin = bytes.NewBufferString(heredoc.Docf(`
		protocol=%s
		host=%s
	`, protocol, hostname))

	err := run.PrepareCmd(rejectCmd).Run()
	if err != nil {
		return err
	}

	approveCmd := git.GitCommand("credential", "approve")

	approveCmd.Stdin = bytes.NewBufferString(heredoc.Docf(`
		protocol=https
		host=%s
		username=%s
		password=%s
	`, hostname, username, password))

	err = run.PrepareCmd(approveCmd).Run()
	if err != nil {
		return err
	}

	return nil
}

func gitCredentialHelperKey(hostname, protocol string) string {
	return fmt.Sprintf("credential.%s://%s.helper", protocol, hostname)
}

func gitCredentialHelper(hostname, protocol string) (helper string, err error) {
	helper, err = git.Config(gitCredentialHelperKey(hostname, protocol))
	if helper != "" {
		return
	}
	helper, err = git.Config("credential.helper")
	return
}

func isOurCredentialHelper(cmd string) bool {
	if !strings.HasPrefix(cmd, "!") {
		return false
	}

	args, err := shlex.Split(cmd[1:])
	if err != nil || len(args) == 0 {
		return false
	}

	return strings.TrimSuffix(filepath.Base(args[0]), ".exe") == "glab"
}

func shellQuote(s string) string {
	if strings.ContainsAny(s, " $") {
		return "'" + s + "'"
	}
	return s
}
