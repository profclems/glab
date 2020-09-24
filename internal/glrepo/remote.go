package glrepo

import (
	"fmt"
	"github.com/profclems/glab/internal/git"
	"net/url"
	"strings"
)

// Remotes represents a set of git remotes
type Remotes []*Remote

// FindByName returns the first Remote whose name matches the list
func (r Remotes) FindByName(names ...string) (*Remote, error) {
	for _, name := range names {
		for _, rem := range r {
			if rem.Name == name || name == "*" {
				return rem, nil
			}
		}
	}
	return nil, fmt.Errorf("no GitLab remotes found")
}

// FindByRepo returns the first Remote that points to a specific GitHub repository
func (r Remotes) FindByRepo(owner, name string) (*Remote, error) {
	for _, rem := range r {
		if strings.EqualFold(rem.RepoOwner(), owner) && strings.EqualFold(rem.RepoName(), name) {
			return rem, nil
		}
	}
	return nil, fmt.Errorf("no matching remote found")
}

func remoteNameSortScore(name string) int {
	switch strings.ToLower(name) {
	case "upstream":
		return 3
	case "gitlab":
		return 2
	case "origin":
		return 1
	default:
		return 0
	}
}

// https://golang.org/pkg/sort/#Interface
func (r Remotes) Len() int      { return len(r) }
func (r Remotes) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r Remotes) Less(i, j int) bool {
	return remoteNameSortScore(r[i].Name) > remoteNameSortScore(r[j].Name)
}

// Remote represents a git remote mapped to a GitHub repository
type Remote struct {
	*git.Remote
	Repo Interface
}

func (r Remote) FullName() string {
	return r.Repo.FullName()
}

// RepoName is the name of the GitHub repository
func (r Remote) RepoName() string {
	return r.Repo.RepoName()
}

// RepoOwner is the name of the GitHub account that owns the repo
func (r Remote) RepoOwner() string {
	return r.Repo.RepoOwner()
}

// RepoHost is the GitHub hostname that the remote points to
func (r Remote) RepoHost() string {
	return r.Repo.RepoHost()
}

// TODO: accept an interface instead of git.RemoteSet
func TranslateRemotes(gitRemotes git.RemoteSet, urlTranslate func(*url.URL) *url.URL) (remotes Remotes) {
	for _, r := range gitRemotes {
		var repo Interface
		if r.FetchURL != nil {
			repo, _ = FromURL(urlTranslate(r.FetchURL))
		}
		if r.PushURL != nil && repo == nil {
			repo, _ = FromURL(urlTranslate(r.PushURL))
		}
		if repo == nil {
			continue
		}
		remotes = append(remotes, &Remote{
			Remote: r,
			Repo:   repo,
		})
	}
	return
}