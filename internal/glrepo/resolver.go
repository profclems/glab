package glrepo

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/profclems/glab/internal/git"
	gLab "github.com/profclems/glab/internal/gitlab"
	"github.com/xanzy/go-gitlab"
	"sort"
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
	for _, r := range result.remotes {
		networkResult, err := gLab.GetProject(result.apiClient, r)
		if err == nil {
			result.network = append(result.network, *networkResult)
		}
		if len(result.remotes) == maxRemotesForLookup {
			break
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

func (r *ResolvedRemotes) BaseRepo(prompt bool) (Interface, error) {
	if r.baseOverride != nil {
		return r.baseOverride, nil
	}

	// if any of the remotes already has a resolution, respect that
	for _, r := range r.remotes {
		if r.Resolved == "base" {
			return r, nil
		} else if r.Resolved != "" {
			repo, err := FromFullName(r.Resolved)
			if err != nil {
				return nil, err
			}
			return NewWithHost(repo.RepoOwner(), repo.RepoName(), r.RepoHost()), nil
		}
	}

	if !prompt {
		// we cannot prompt, so just resort to the 1st remote
		return r.remotes[0], nil
	}

	// from here on, consult the API
	if r.network == nil {
		err := resolveNetwork(r)
		if err != nil {
			return nil, err
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

	for _, repo := range r.network {
		if repo.ForkedFromProject != nil {
			fProject, _ := gLab.GetProject(r.apiClient, repo.ForkedFromProject.PathWithNamespace)
			add(fProject)
		}
		add(&repo)
	}

	if len(repoNames) == 0 {
		return r.remotes[0], nil
	}

	baseName := repoNames[0]
	if len(repoNames) > 1 {
		err := survey.AskOne(&survey.Select{
			Message: "Which should be the base repository (used for e.g. querying issues) for this directory?",
			Options: repoNames,
		}, &baseName)
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
	}

	// cache the result to git config
	err := git.SetRemoteResolution(remote.Name, resolution)
	return selectedRepoInfo, err
}

func (r *ResolvedRemotes) HeadRepos() ([]*gitlab.Project, error) {
	if r.network == nil {
		err := resolveNetwork(r)
		if err != nil {
			return nil, err
		}
	}

	var results []*gitlab.Project
	for _, repo := range r.network {
		results = append(results, &repo)
	}
	return results, nil
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
