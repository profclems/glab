package git

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/profclems/glab/internal/run"
)

func GetRemoteURL(remoteAlias string) (string, error) {
	return Config("remote." + remoteAlias + ".url")
}

// GetDefaultBranch finds and returns the remote's default branch
func GetDefaultBranch(remote string) (string, error) {
	getDefBranch := exec.Command("git", "remote", "show", remote)
	output, err := run.PrepareCmd(getDefBranch).Output()
	if err != nil {
		return "master", err
	}

	var headBranch string

	for _, o := range strings.Split(string(output), "\n") {
		o = strings.TrimSpace(o)
		r, err := regexp.Compile(`(HEAD branch:)\s+`)
		if err != nil {
			return "master", err
		}
		if r.MatchString(o) {
			headBranch = strings.TrimPrefix(o, "HEAD branch: ")
			break
		}
	}
	return headBranch, err
}

// ErrNotOnAnyBranch indicates that the users is in detached HEAD state
var ErrNotOnAnyBranch = errors.New("git: not on any branch")

// Ref represents a git commit reference
type Ref struct {
	Hash string
	Name string
}

// TrackingRef represents a ref for a remote tracking branch
type TrackingRef struct {
	RemoteName string
	BranchName string
}

func (r TrackingRef) String() string {
	return "refs/remotes/" + r.RemoteName + "/" + r.BranchName
}

// ShowRefs resolves fully-qualified refs to commit hashes
func ShowRefs(ref ...string) ([]Ref, error) {
	args := append([]string{"show-ref", "--verify", "--"}, ref...)
	showRef := exec.Command("git", args...)
	output, err := run.PrepareCmd(showRef).Output()

	var refs []Ref
	for _, line := range outputLines(output) {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		refs = append(refs, Ref{
			Hash: parts[0],
			Name: parts[1],
		})
	}

	return refs, err
}

// CurrentBranch reads the checked-out branch for the git repository
func CurrentBranch() (string, error) {
	refCmd := GitCommand("symbolic-ref", "--quiet", "--short", "HEAD")

	output, err := run.PrepareCmd(refCmd).Output()
	if err == nil {
		// Found the branch name
		return firstLine(output), nil
	}

	var cmdErr *run.CmdError
	if errors.As(err, &cmdErr) {
		if cmdErr.Stderr.Len() == 0 {
			// Detached head
			return "", ErrNotOnAnyBranch
		}
	}

	// Unknown error
	return "", err
}

func listRemotes() ([]string, error) {
	remoteCmd := exec.Command("git", "remote", "-v")
	output, err := run.PrepareCmd(remoteCmd).Output()
	return outputLines(output), err
}

func Config(name string) (string, error) {
	configCmd := exec.Command("git", "config", name)
	output, err := run.PrepareCmd(configCmd).Output()
	if err != nil {
		return "", fmt.Errorf("unknown config key: %s", name)
	}

	return firstLine(output), nil

}

var GitCommand = func(args ...string) *exec.Cmd {
	return exec.Command("git", args...)
}

func UncommittedChangeCount() (int, error) {
	statusCmd := GitCommand("status", "--porcelain")
	output, err := run.PrepareCmd(statusCmd).Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(output), "\n")

	count := 0

	for _, l := range lines {
		if l != "" {
			count++
		}
	}

	return count, nil
}

type Commit struct {
	Sha   string
	Title string
}

func LatestCommit(ref string) (*Commit, error) {
	logCmd := GitCommand("show", "-s", "--format=%h %s", ref)
	output, err := run.PrepareCmd(logCmd).Output()
	if err != nil {
		return &Commit{}, err
	}
	commit := &Commit{}
	split := strings.SplitN(string(output), " ", 2)
	if len(split) != 2 {
		return commit, fmt.Errorf("could not find commit for %s", ref)
	}
	commit = &Commit{
		Sha:   split[0],
		Title: split[1],
	}
	return commit, nil
}

func Commits(baseRef, headRef string) ([]*Commit, error) {
	logCmd := GitCommand(
		"-c", "log.ShowSignature=false",
		"log", "--pretty=format:%H,%s",
		"--cherry", fmt.Sprintf("%s...%s", baseRef, headRef))
	output, err := run.PrepareCmd(logCmd).Output()
	if err != nil {
		return []*Commit{}, err
	}

	var commits []*Commit
	sha := 0
	title := 1
	for _, line := range outputLines(output) {
		split := strings.SplitN(line, ",", 2)
		if len(split) != 2 {
			continue
		}
		commits = append(commits, &Commit{
			Sha:   split[sha],
			Title: split[title],
		})
	}

	if len(commits) == 0 {
		return commits, fmt.Errorf("could not find any commits between %s and %s", baseRef, headRef)
	}

	return commits, nil
}

func CommitBody(sha string) (string, error) {
	showCmd := GitCommand("-c", "log.ShowSignature=false", "show", "-s", "--pretty=format:%b", sha)
	output, err := run.PrepareCmd(showCmd).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Push publishes a git ref to a remote
func Push(remote string, ref string, cmdOut, cmdErr io.Writer) error {
	pushCmd := GitCommand("push", remote, ref)
	pushCmd.Stdout = cmdOut
	pushCmd.Stderr = cmdErr
	return run.PrepareCmd(pushCmd).Run()
}

// SetUpstream sets the upstream (tracking) of a branch
func SetUpstream(remote string, branch string, cmdOut, cmdErr io.Writer) error {
	setCmd := GitCommand("branch", "--set-upstream-to", fmt.Sprintf("%s/%s", remote, branch))
	setCmd.Stdout = cmdOut
	setCmd.Stderr = cmdErr
	return run.PrepareCmd(setCmd).Run()
}

// Rebase rebases the current branch on top of its tracking branch
func Rebase(cmdOut, cmdErr io.Writer) error {
	rebaseCmd := GitCommand("rebase")
	rebaseCmd.Stdout = cmdOut
	rebaseCmd.Stderr = cmdErr
	return run.PrepareCmd(rebaseCmd).Run()
}

type BranchConfig struct {
	RemoteName string
	RemoteURL  *url.URL
	MergeRef   string
}

// ReadBranchConfig parses the `branch.BRANCH.(remote|merge)` part of git config
func ReadBranchConfig(branch string) (cfg BranchConfig) {
	prefix := regexp.QuoteMeta(fmt.Sprintf("branch.%s.", branch))
	configCmd := GitCommand("config", "--get-regexp", fmt.Sprintf("^%s(remote|merge)$", prefix))
	output, err := run.PrepareCmd(configCmd).Output()
	if err != nil {
		return
	}
	for _, line := range outputLines(output) {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		keys := strings.Split(parts[0], ".")
		switch keys[len(keys)-1] {
		case "remote":
			if strings.Contains(parts[1], ":") {
				u, err := ParseURL(parts[1])
				if err != nil {
					continue
				}
				cfg.RemoteURL = u
			} else if !isFilesystemPath(parts[1]) {
				cfg.RemoteName = parts[1]
			}
		case "merge":
			cfg.MergeRef = parts[1]
		}
	}
	return
}

func DeleteLocalBranch(branch string) error {
	branchCmd := GitCommand("branch", "-D", branch)
	err := run.PrepareCmd(branchCmd).Run()
	return err
}

func HasLocalBranch(branch string) bool {
	configCmd := GitCommand("rev-parse", "--verify", "refs/heads/"+branch)
	_, err := run.PrepareCmd(configCmd).Output()
	return err == nil
}

func CheckoutBranch(branch string) error {
	configCmd := GitCommand("checkout", branch)
	err := run.PrepareCmd(configCmd).Run()
	return err
}

func parseCloneArgs(extraArgs []string) (args []string, target string) {
	args, target = parseArgs(extraArgs)
	return
}

func parseArgs(cmdWithArgs []string) (args []string, command string) {
	args = cmdWithArgs

	if len(args) > 0 {
		if !strings.HasPrefix(args[0], "-") {
			command, args = args[0], args[1:]
		}
	}
	return
}

func RunClone(cloneURL string, args []string) (target string, err error) {
	cloneArgs, target := parseCloneArgs(args)

	cloneArgs = append(cloneArgs, cloneURL)

	// If the args contain an explicit target, pass it to clone
	//    otherwise, parse the URL to determine where git cloned it to so we can return it
	if target != "" {
		cloneArgs = append(cloneArgs, target)
	} else {
		target = path.Base(strings.TrimSuffix(cloneURL, ".git"))
	}

	cloneArgs = append([]string{"clone"}, cloneArgs...)

	cloneCmd := GitCommand(cloneArgs...)
	cloneCmd.Stdin = os.Stdin
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	err = run.PrepareCmd(cloneCmd).Run()
	return
}

func AddUpstreamRemote(upstreamURL, cloneDir string) error {
	cloneCmd := GitCommand("-C", cloneDir, "remote", "add", "-f", "upstream", upstreamURL)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	return run.PrepareCmd(cloneCmd).Run()
}

func isFilesystemPath(p string) bool {
	return p == "." || strings.HasPrefix(p, "./") || strings.HasPrefix(p, "/")
}

// ToplevelDir returns the top-level directory path of the current repository
func ToplevelDir() (string, error) {
	showCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := run.PrepareCmd(showCmd).Output()
	return firstLine(output), err

}

func outputLines(output []byte) []string {
	lines := strings.TrimSuffix(string(output), "\n")
	return strings.Split(lines, "\n")

}

func firstLine(output []byte) string {
	if i := bytes.IndexAny(output, "\n"); i >= 0 {
		return string(output)[0:i]
	}
	return string(output)
}

var remoteRE = regexp.MustCompile(`(.+)\s+(.+)\s+\((push|fetch)\)`)

// RemoteSet is a slice of git remotes
type RemoteSet []*Remote

func NewRemote(name string, u string) *Remote {
	pu, _ := url.Parse(u)
	return &Remote{
		Name:     name,
		FetchURL: pu,
		PushURL:  pu,
	}
}

// Remote is a parsed git remote
type Remote struct {
	Name     string
	Resolved string
	FetchURL *url.URL
	PushURL  *url.URL
}

func (r *Remote) String() string {
	return r.Name
}

// Remotes gets the git remotes set for the current repo
func Remotes() (RemoteSet, error) {
	list, err := listRemotes()
	if err != nil {
		return nil, err
	}
	remotes := parseRemotes(list)

	// this is affected by SetRemoteResolution
	remoteCmd := exec.Command("git", "config", "--get-regexp", `^remote\..*\.glab-resolved$`)
	output, _ := run.PrepareCmd(remoteCmd).Output()
	for _, l := range outputLines(output) {
		parts := strings.SplitN(l, " ", 2)
		if len(parts) < 2 {
			continue
		}
		rp := strings.SplitN(parts[0], ".", 3)
		if len(rp) < 2 {
			continue
		}
		name := rp[1]
		for _, r := range remotes {
			if r.Name == name {
				r.Resolved = parts[1]
				break
			}
		}
	}

	return remotes, nil
}

func parseRemotes(gitRemotes []string) (remotes RemoteSet) {
	for _, r := range gitRemotes {
		match := remoteRE.FindStringSubmatch(r)
		if match == nil {
			continue
		}
		name := strings.TrimSpace(match[1])
		urlStr := strings.TrimSpace(match[2])
		urlType := strings.TrimSpace(match[3])

		var rem *Remote
		if len(remotes) > 0 {
			rem = remotes[len(remotes)-1]
			if name != rem.Name {
				rem = nil
			}
		}
		if rem == nil {
			rem = &Remote{Name: name}
			remotes = append(remotes, rem)
		}

		u, err := ParseURL(urlStr)
		if err != nil {
			continue
		}

		switch urlType {
		case "fetch":
			rem.FetchURL = u
		case "push":
			rem.PushURL = u
		}
	}
	return
}

// AddRemote adds a new git remote and auto-fetches objects from it
func AddRemote(name, u string) (*Remote, error) {
	addCmd := exec.Command("git", "remote", "add", "-f", name, u)
	err := run.PrepareCmd(addCmd).Run()
	if err != nil {
		return nil, err
	}

	var urlParsed *url.URL
	if strings.HasPrefix(u, "https") {
		urlParsed, err = url.Parse(u)
		if err != nil {
			return nil, err
		}

	} else {
		urlParsed, err = ParseURL(u)
		if err != nil {
			return nil, err
		}

	}

	return &Remote{
		Name:     name,
		FetchURL: urlParsed,
		PushURL:  urlParsed,
	}, nil
}

var SetRemoteResolution = func(name, resolution string) error {
	return SetConfig(name, "glab-resolved", resolution)
}

func SetConfig(remote, key, value string) error {
	addCmd := exec.Command("git", "config", "--add", fmt.Sprintf("remote.%s.%s", remote, key), value)
	return run.PrepareCmd(addCmd).Run()
}

func RunCmd(args []string) (err error) {
	gitCmd := GitCommand(args...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	err = run.PrepareCmd(gitCmd).Run()
	return
}
