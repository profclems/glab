package glrepo

import (
	"errors"
	"sort"
	"strings"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/prompt"

	"github.com/xanzy/go-gitlab"
)

// cap the number of git remotes looked up, since the user might have an
// unusually large number of git remotes
const maxRemotesForLookup = 5

func ResolveRemotesToRepos(remotes Remotes, client *gitlab.Client, base string) (*ResolvedRemotes, error) {
	sort.Stable(remotes)

	result := &ResolvedRemotes{
		remotes:   remotes,
		apiClient: client,
	}

	var baseOverride Interface
	if base != "" {
		var err error
		baseOverride, err = FromFullName(base)
		if err != nil {
			return result, err
		}
		result.baseOverride = baseOverride
	}

	return result, nil
}

func resolveNetwork(result *ResolvedRemotes) error {
	// Loop over at most 5 (maxRemotesForLookup)
	for i := 0; i < len(result.remotes) && i < maxRemotesForLookup; i++ {
		networkResult, err := api.GetProject(result.apiClient, result.remotes[i].FullName())
		if err == nil {
			result.network = append(result.network, *networkResult)
		} else {
			return err
		}
	}
	return nil
}

type ResolvedRemotes struct {
	baseOverride Interface
	remotes      Remotes
	network      []gitlab.Project
	apiClient    *gitlab.Client
}

func (r *ResolvedRemotes) BaseRepo(interactive bool) (Interface, error) {
	if r.baseOverride != nil {
		return r.baseOverride, nil
	}

	// if any of the remotes already has a resolution, respect that
	for _, r := range r.remotes {
		if r.Resolved == "base" {
			return r, nil
		} else if strings.HasPrefix(r.Resolved, "base:") {
			repo, err := FromFullName(strings.TrimPrefix(r.Resolved, "base:"))
			if err != nil {
				return nil, err
			}
			return NewWithHost(repo.RepoOwner(), repo.RepoName(), r.RepoHost()), nil
		} else if r.Resolved != "" && !strings.HasPrefix(r.Resolved, "head") {
			// Backward compatibility kludge for remoteless resolutions created before
			// BaseRepo started creeating resolutions prefixed with `base:`
			repo, err := FromFullName(r.Resolved)
			if err != nil {
				return nil, err
			}
			// Rewrite resolution, ignore the error as this will keep working
			// in the future we might add a warning that we couldn't rewrite
			// it for compatiblity
			_ = git.SetRemoteResolution(r.Name, "base:"+r.Resolved)

			return NewWithHost(repo.RepoOwner(), repo.RepoName(), r.RepoHost()), nil
		}
	}

	if !interactive {
		// we cannot prompt, so just resort to the 1st remote
		return r.remotes[0], nil
	}

	// from here on, consult the API
	if r.network == nil {
		err := resolveNetwork(r)
		if err != nil {
			return nil, err
		}
		if len(r.network) == 0 {
			return nil, errors.New("no GitLab Projects found from remotes")
		}
	}

	var repoNames []string
	repoMap := map[string]*gitlab.Project{}
	add := func(r *gitlab.Project) {
		fn, _ := FullNameFromURL(r.HTTPURLToRepo)
		if _, ok := repoMap[fn]; !ok {
			repoMap[fn] = r
			repoNames = append(repoNames, fn)
		}
	}

	for i := range r.network {
		if r.network[i].ForkedFromProject != nil {
			fProject, _ := api.GetProject(r.apiClient, r.network[i].ForkedFromProject.PathWithNamespace)
			add(fProject)
		}
		add(&r.network[i])
	}

	baseName := repoNames[0]
	if len(repoNames) > 1 {
		err := prompt.Select(
			&baseName,
			"base",
			"Which should be the base repository (used for e.g. querying issues) for this directory?",
			repoNames,
		)
		if err != nil {
			return nil, err
		}
	}

	// determine corresponding git remote
	selectedRepo := repoMap[baseName]
	selectedRepoInfo, _ := FromFullName(selectedRepo.HTTPURLToRepo)
	resolution := "base"
	remote, _ := r.RemoteForRepo(selectedRepoInfo)
	if remote == nil {
		remote = r.remotes[0]
		resolution, _ = FullNameFromURL(selectedRepo.HTTPURLToRepo)
		resolution = "base:" + resolution
	}

	// cache the result to git config
	err := git.SetRemoteResolution(remote.Name, resolution)
	return selectedRepoInfo, err
}

func (r *ResolvedRemotes) HeadRepo(interactive bool) (Interface, error) {
	if r.baseOverride != nil {
		return r.baseOverride, nil
	}

	// if any of the remotes already has a resolution, respect that
	for _, r := range r.remotes {
		if r.Resolved == "head" {
			return r, nil
		} else if strings.HasPrefix(r.Resolved, "head:") {
			repo, err := FromFullName(strings.TrimPrefix(r.Resolved, "head:"))
			if err != nil {
				return nil, err
			}
			return NewWithHost(repo.RepoOwner(), repo.RepoName(), r.RepoHost()), nil
		}
	}

	// from here on, consult the API
	if r.network == nil {
		err := resolveNetwork(r)
		if err != nil {
			return nil, err
		}
		if len(r.network) == 0 {
			return nil, errors.New("no GitLab Projects found from remotes")
		}
	}

	var repoNames []string
	repoMap := map[string]*gitlab.Project{}
	add := func(r *gitlab.Project) {
		fn, _ := FullNameFromURL(r.HTTPURLToRepo)
		if _, ok := repoMap[fn]; !ok {
			repoMap[fn] = r
			repoNames = append(repoNames, fn)
		}
	}

	for i := range r.network {
		if r.network[i].ForkedFromProject != nil {
			fProject, _ := api.GetProject(r.apiClient, r.network[i].ForkedFromProject.PathWithNamespace)
			add(fProject)
		}
		add(&r.network[i])
	}

	headName := repoNames[0]
	if len(repoNames) > 1 {
		if !interactive {
			// We cannot prompt so get the first repo that is a fork
			for _, repo := range repoNames {
				if repoMap[repo].ForkedFromProject != nil {
					selectedRepoInfo, _ := FromFullName((repoMap[repo].HTTPURLToRepo))
					remote, _ := r.RemoteForRepo(selectedRepoInfo)
					return remote, nil
				}
			}
			// There are no forked repos so return the first repo
			return r.remotes[0], nil
		}

		err := prompt.Select(
			&headName,
			"head",
			"Which should be the head repository (where branches are pushed) for this directory?",
			repoNames,
		)
		if err != nil {
			return nil, err
		}
	}

	// determine corresponding git remote
	selectedRepo := repoMap[headName]
	selectedRepoInfo, _ := FromFullName(selectedRepo.HTTPURLToRepo)
	resolution := "head"
	remote, _ := r.RemoteForRepo(selectedRepoInfo)
	if remote == nil {
		remote = r.remotes[0]
		resolution, _ = FullNameFromURL(selectedRepo.HTTPURLToRepo)
		resolution = "head:" + resolution
	}

	// cache the result to git config
	err := git.SetRemoteResolution(remote.Name, resolution)
	return selectedRepoInfo, err
}

// RemoteForRepo finds the git remote that points to a repository
func (r *ResolvedRemotes) RemoteForRepo(repo Interface) (*Remote, error) {
	for _, remote := range r.remotes {
		if IsSame(remote, repo) {
			return remote, nil
		}
	}
	return nil, errors.New("not found")
}
