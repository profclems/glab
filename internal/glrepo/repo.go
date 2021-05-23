package glrepo

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/profclems/glab/api"

	"github.com/profclems/glab/internal/config"

	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/utils"
	"github.com/xanzy/go-gitlab"
)

// RemoteURL returns correct git clone URL of a repo
// based on the user's git_protocol preference
func RemoteURL(project *gitlab.Project, protocol string) string {
	if protocol == "ssh" {
		return project.SSHURLToRepo
	}
	return project.HTTPURLToRepo
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
// Contains methods for these methods representing these placeholders for a
// project path with :host/:group/:namespace/:repo
// RepoHost = :host, RepoOwner = :group/:namespace, RepoNamespace = :namespace,
// FullName = :group/:namespace/:repo, RepoGroup = :group, RepoName = :repo
type Interface interface {
	RepoName() string
	RepoOwner() string
	RepoNamespace() string
	RepoGroup() string
	RepoHost() string
	FullName() string
	Project(*gitlab.Client) (*gitlab.Project, error)
}

// New instantiates a GitLab repository from owner and repo name arguments
func New(owner, repo string) Interface {
	return NewWithHost(owner, repo, glinstance.OverridableDefault())
}

// NewWithGroup instantiates a GitLab repository from group, namespace and repo name arguments
func NewWithGroup(group, namespace, repo, hostname string) Interface {
	owner := fmt.Sprintf("%s/%s", group, namespace)
	if hostname == "" {
		return New(owner, repo)
	}
	return NewWithHost(owner, repo, hostname)
}

// NewWithHost is like New with an explicit host name
func NewWithHost(owner, repo, hostname string) Interface {
	rp := &glRepo{
		owner:    owner,
		name:     repo,
		fullname: fmt.Sprintf("%s/%s", owner, repo),
		hostname: normalizeHostname(hostname),
	}
	if ri := strings.SplitN(owner, "/", 2); len(ri) == 2 {
		rp.group = ri[0]
		rp.namespace = ri[1]
	} else {
		rp.namespace = owner
	}
	return rp
}

// FromFullName extracts the GitLab repository information from the following
// formats: "OWNER/REPO", "HOST/OWNER/REPO", "HOST/GROUP/NAMESPACE/REPO", and a full URL.
func FromFullName(nwo string) (Interface, error) {
	nwo = strings.TrimSpace(nwo)
	// check if it's a valid git URL and parse it
	if git.IsValidURL(nwo) {
		u, err := git.ParseURL(nwo)
		if err != nil {
			return nil, err
		}
		return FromURL(u)
	}
	// check if it is valid URL and parse it
	if utils.IsValidURL(nwo) {
		u, _ := url.Parse(nwo)
		return FromURL(u)
	}

	repo := nwo[strings.LastIndex(nwo, "/")+1:]
	nwoWithoutRepo := strings.TrimSuffix(nwo[:strings.LastIndex(nwo, "/")+1], "/")
	parts := strings.SplitN(nwoWithoutRepo, "/", 2)

	if repo == "" {
		return nil, fmt.Errorf(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got %q`, nwo)
	}
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got %q`, nwo)
		}
	}
	switch len(parts) {
	case 2: // GROUP/NAMESPACE/REPO or HOST/OWNER/REPO or //HOST/GROUP/NAMESPACE/REPO
		// First, checks if the first part matches the default instance host (i.e. gitlab.com) or the
		// overridden default host (mostly from the GITLAB_HOST env variable)
		if parts[0] == glinstance.Default() || parts[0] == glinstance.OverridableDefault() {
			return NewWithHost(parts[1], repo, normalizeHostname(parts[0])), nil
		}
		// Dots (.) are allowed in group names by GitLab.
		// So we check if if the first part contains a dot.
		// However, it could be that the user is specifying a hostname but we can't be sure of that
		// So we check in the list of authenticated hosts and see if it matches any
		// if not, we assume it is a group name that contains a dot
		if strings.ContainsRune(parts[0], '.') {
			var rI Interface
			cfg, err := config.Init()
			if err == nil {
				hosts, _ := cfg.Hosts()
				for _, host := range hosts {
					if host == parts[0] {
						rI = NewWithHost(parts[1], repo, normalizeHostname(parts[0]))
						break
					}
				}
				if rI != nil {
					return rI, nil
				}
			}
		}
		// if the first part is not a valid URL, and does not match an
		// authenticated hostname then we assume it is in
		// the format GROUP/NAMESPACE/REPO
		return NewWithGroup(parts[0], parts[1], repo, ""), nil
	case 1: // OWNER/REPO
		return New(parts[0], repo), nil
	default:
		return nil, fmt.Errorf(`expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got %q`, nwo)
	}
}

// FromURL extracts the GitLab repository information from a git remote URL
func FromURL(u *url.URL) (Interface, error) {
	if u.Hostname() == "" {
		return nil, fmt.Errorf("no hostname detected")
	}
	path := strings.Trim(strings.TrimSuffix(u.Path, ".git"), "/")
	repo := path[strings.LastIndex(path, "/")+1:]
	pathWithoutRepo := strings.TrimSuffix(path[:strings.LastIndex(path, "/")+1], "/")
	if repo != "" && pathWithoutRepo != "" {
		parts := strings.SplitN(pathWithoutRepo, "/", 2)
		if len(parts) == 1 {
			return NewWithHost(parts[0], repo, u.Hostname()), nil
		}
		if len(parts) == 2 {
			return NewWithGroup(parts[0], parts[1], repo, u.Hostname()), nil
		}
	}
	return nil, fmt.Errorf("invalid path: %s", u.Path)
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
	group     string
	owner     string
	name      string
	fullname  string
	hostname  string
	namespace string

	project *Project
}

type Project struct {
	*gitlab.Project
	// for cache invalidation
	fullname string
	// for cache invalidation
	hostname string
}

func (r glRepo) Project(client *gitlab.Client) (*gitlab.Project, error) {
	if r.project != nil && r.project.fullname == r.fullname && r.project.hostname == r.hostname {
		return r.project.Project, nil
	}
	p, err := api.GetProject(client, r.fullname)
	if err != nil {
		return nil, err
	}
	r.project = &Project{
		Project:  p,
		fullname: r.fullname,
		hostname: r.hostname,
	}
	return r.project.Project, err
}

// RepoNamespace returns the namespace of the project. Eg. if project path is :group/:namespace:/repo
// RepoNamespace returns the :namespace
func (r glRepo) RepoNamespace() string {
	return r.namespace
}

// RepoGroup returns the group namespace of the project. Eg. if project path is :group/:namespace:/repo
// RepoGroup returns the :group
func (r glRepo) RepoGroup() string {
	return r.group
}

// RepoOwner returns the group and namespace in the form "group/namespace". Returns "namespace" if group is not present
func (r glRepo) RepoOwner() string {
	if r.group != "" {
		return r.group + "/" + r.namespace
	}
	return r.owner
}

// RepoName returns the repo name without the path or namespace.
func (r glRepo) RepoName() string {
	return r.name
}

// RepoHost returns the hostname
func (r glRepo) RepoHost() string {
	return r.hostname
}

// FullName returns the full project path :group/:namespace/:repo or :namespace/:repo if group is not present
func (r glRepo) FullName() string {
	return r.fullname
}
