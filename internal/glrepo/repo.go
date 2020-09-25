package glrepo

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glinstance"

	"github.com/xanzy/go-gitlab"
)

type RemoteArgs struct {
	protocol string
	token    string
	url      string
	username string
}

// RemoteURL returns correct git clone URL of a repo
// based on the user's git_protocol preference
func RemoteURL(project *gitlab.Project, a *RemoteArgs) (string, error) {
	if a.protocol == "https" {

		if a.username == "" {
			a.username = "oauth2"
		}

		a.protocol = "https://"
		if strings.Contains(a.url, "https://") {
			a.url = strings.TrimPrefix(a.url, "https://")
		} else if strings.HasPrefix(a.url, "http://") {
			a.url = strings.TrimPrefix(a.url, "http://")
			a.protocol = "http://"
		}
		return fmt.Sprintf("%s%s:%s@%s/%s.git",
			a.protocol, a.username, a.token, a.url, project.PathWithNamespace), nil
	}
	return project.SSHURLToRepo, nil
}

// FullName returns the the repo with its namespace (like profclems/glab). Respects group and subgroups names
func FullNameFromURL(remoteURL string) (string, error) {
	parts := strings.Split(remoteURL, "//")

	if len(parts) == 1 {
		// scp-like short syntax (e.g. git@gitlab.com...)
		part := parts[0]
		parts = strings.Split(part, ":")
	} else if len(parts) == 2 {
		// other protocols (e.g. ssh://, http://, git://)
		part := parts[1]
		parts = strings.SplitN(part, "/", 2)
	} else {
		return "", errors.New("cannot parse remote: " + remoteURL)
	}

	if len(parts) != 2 {
		return "", errors.New("cannot parse remote: " + remoteURL)
	}
	repo := parts[1]
	repo = strings.TrimSuffix(repo, ".git")
	return repo, nil
}

// Interface describes an object that represents a GitLab repository
type Interface interface {
	RepoName() string
	RepoOwner() string
	RepoHost() string
	FullName() string
}

// New instantiates a GitLab repository from owner and name arguments
func New(owner, repo string) Interface {
	return NewWithHost(owner, repo, glinstance.OverridableDefault())
}

// NewWithHost is like New with an explicit host name
func NewWithHost(owner, repo, hostname string) Interface {
	return &glRepo{
		owner:    owner,
		name:     repo,
		fullname: fmt.Sprintf("%s/%s", owner, repo),
		hostname: normalizeHostname(hostname),
	}
}

// FromFullName extracts the GitLab repository information from the following
// formats: "OWNER/REPO", "HOST/OWNER/REPO", "HOST/GROUP/NAMESPACE/REPO", and a full URL.
func FromFullName(nwo string) (Interface, error) {
	if git.IsValidURL(nwo) {
		u, err := git.ParseURL(nwo)
		if err != nil {
			return nil, err
		}
		return FromURL(u)
	}

	parts := strings.SplitN(nwo, "/", 3)
	for _, p := range parts {
		if len(p) == 0 {
			return nil, fmt.Errorf(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got %q`, nwo)
		}
	}
	switch len(parts) {
	case 3:
		return NewWithHost(parts[1], parts[2], normalizeHostname(parts[0])), nil
	case 2:
		return New(parts[0], parts[1]), nil
	default:
		return nil, fmt.Errorf(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got %q`, nwo)
	}
}

// FromURL extracts the GitLab repository information from a git remote URL
func FromURL(u *url.URL) (Interface, error) {
	if u.Hostname() == "" {
		return nil, fmt.Errorf("no hostname detected")
	}

	parts := strings.SplitN(strings.Trim(u.Path, "/"), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid path: %s", u.Path)
	}
	return NewWithHost(parts[0], strings.TrimSuffix(parts[1], ".git"), u.Hostname()), nil
}

func normalizeHostname(h string) string {
	return strings.ToLower(strings.TrimPrefix(h, "www."))
}

// IsSame compares two GitLab repositories
func IsSame(a, b Interface) bool {
	return strings.EqualFold(a.RepoOwner(), b.RepoOwner()) &&
		strings.EqualFold(a.RepoName(), b.RepoName()) &&
		normalizeHostname(a.RepoHost()) == normalizeHostname(b.RepoHost())
}

type glRepo struct {
	owner    string
	name     string
	fullname string
	hostname string
}

func (r glRepo) RepoOwner() string {
	return r.owner
}

func (r glRepo) RepoName() string {
	return r.name
}

func (r glRepo) RepoHost() string {
	return r.hostname
}

func (r glRepo) FullName() string {
	return r.fullname
}
