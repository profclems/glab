package glrepo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/pkg/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_RemoteForRepo(t *testing.T) {
	r := &ResolvedRemotes{
		remotes: Remotes{
			&Remote{
				Remote: &git.Remote{
					Name: "upstream",
				},
				Repo: NewWithHost("profclems", "glab", "gitlab.com"),
			},
			&Remote{
				Remote: &git.Remote{
					Name: "origin",
				},
				Repo: NewWithHost("maxice8", "glab", "gitlab.com"),
			},
		},
	}
	testCases := []struct {
		name    string
		input   Interface
		output  *Remote // Expected remote if there is a match
		wantErr string  // Expected error
	}{
		{
			name:  "match upstream",
			input: NewWithHost("profclems", "glab", "gitlab.com"),
			output: &Remote{
				Remote: &git.Remote{
					Name: "upstream",
				},
				Repo: NewWithHost("profclems", "glab", "gitlab.com"),
			},
		},
		{
			name:  "match origin",
			input: NewWithHost("maxice8", "glab", "gitlab.com"),
			output: &Remote{
				Remote: &git.Remote{
					Name: "origin",
				},
				Repo: NewWithHost("maxice8", "glab", "gitlab.com"),
			},
		},
		{
			name:    "no match via Hostname",
			input:   NewWithHost("profclems", "glab", "gitlab.extradomain.com"),
			wantErr: "not found",
		},
		{
			name:    "no match via Username",
			input:   NewWithHost("noexist", "glab", "gitlab.com"),
			wantErr: "not found",
		},
		{
			name:    "no match via Name",
			input:   NewWithHost("profclems", "maxice8", "gitlab.com"),
			wantErr: "not found",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got, err := r.RemoteForRepo(tC.input)
			if tC.wantErr == "" && err != nil {
				t.Errorf("RemoteForRepo() unexpected error = %s", err)
			}
			if tC.wantErr != "" {
				if tC.wantErr != err.Error() {
					t.Errorf("RemoteForRepo() expected error = %s, got = %s", tC.wantErr, err)
				}
			} else {
				// Make sure both return all the exact same thing
				assert.Equal(t, tC.output.Name, got.Name)
				assert.Equal(t, tC.output.Remote.Name, got.Remote.Name)
				assert.Equal(t, tC.output.Repo.FullName(), got.Repo.FullName())
				assert.Equal(t, tC.output.Repo.RepoHost(), got.Repo.RepoHost())
			}
		})
	}
}

func Test_ResolveRemotesToRepos(t *testing.T) {
	rem := &ResolvedRemotes{
		remotes: Remotes{
			&Remote{
				Remote: &git.Remote{
					Name: "origin",
				},
				Repo: NewWithHost("profclems", "glab", "gitlab.com"),
			},
		},
		apiClient: &gitlab.Client{},
	}

	// Test the normal and most expected usage
	t.Run("simple", func(t *testing.T) {
		r, err := ResolveRemotesToRepos(rem.remotes, rem.apiClient, "")
		assert.Nil(t, err)

		assert.Equal(t, rem.apiClient, r.apiClient)

		assert.Len(t, r.remotes, 1)

		for i := range r.remotes {
			assert.Equal(t, r.remotes[i].Name, rem.remotes[i].Name)
			assert.Equal(t, r.remotes[i].Repo.FullName(), rem.remotes[i].Repo.FullName())
			assert.Equal(t, r.remotes[i].Repo.RepoHost(), rem.remotes[i].Repo.RepoHost())
		}
	})

	// Test the usage of baseOverride
	t.Run("baseOverride", func(t *testing.T) {
		expectedBaseOverride := NewWithHost("profclems", "glab", "gitlab.com")

		r, err := ResolveRemotesToRepos(rem.remotes, rem.apiClient, "gitlab.com/profclems/glab")
		assert.Nil(t, err)

		assert.Equal(t, expectedBaseOverride.FullName(), r.baseOverride.FullName())
		assert.Equal(t, expectedBaseOverride.RepoHost(), r.baseOverride.RepoHost())

		assert.Equal(t, rem.apiClient, r.apiClient)

		assert.Len(t, r.remotes, 1)

		for i := range r.remotes {
			assert.Equal(t, r.remotes[i].Name, rem.remotes[i].Name)
			assert.Equal(t, r.remotes[i].Repo.FullName(), rem.remotes[i].Repo.FullName())
			assert.Equal(t, r.remotes[i].Repo.RepoHost(), rem.remotes[i].Repo.RepoHost())
		}
	})

	// Test the usage of baseOverride when it is passed an invalid value
	t.Run("baseOverrideFail", func(t *testing.T) {
		r, err := ResolveRemotesToRepos(rem.remotes, rem.apiClient, "badValue")
		assert.EqualError(t, err, `expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "badValue"`)

		assert.Equal(t, rem.apiClient, r.apiClient)

		assert.Len(t, r.remotes, 1)

		for i := range r.remotes {
			assert.Equal(t, r.remotes[i].Name, rem.remotes[i].Name)
			assert.Equal(t, r.remotes[i].Repo.FullName(), rem.remotes[i].Repo.FullName())
			assert.Equal(t, r.remotes[i].Repo.RepoHost(), rem.remotes[i].Repo.RepoHost())
		}
	})
}

func Test_resolveNetwork(t *testing.T) {
	rem := &ResolvedRemotes{
		remotes: Remotes{
			&Remote{
				Remote: &git.Remote{
					Name: "origin",
				},
				Repo: NewWithHost("profclems", "glab", "gitlab.com"),
			},
		},
		apiClient: &gitlab.Client{},
	}

	// Override api.GetProejct to not use the network
	mockAPIGetProject := func(_ *gitlab.Client, ProjectID interface{}) (*gitlab.Project, error) {
		proj := &gitlab.Project{
			PathWithNamespace: fmt.Sprint(ProjectID),
		}
		return proj, nil
	}

	t.Run("simple", func(t *testing.T) {
		// Make our own copy of rem we can modify
		rem := *rem

		api.GetProject = mockAPIGetProject

		resolveNetwork(&rem)

		assert.Len(t, rem.network, len(rem.remotes))
		for i := range rem.network {
			assert.Equal(t, rem.remotes[i].Repo.FullName(), rem.network[i].PathWithNamespace)
		}
	})

	t.Run("API call failed", func(t *testing.T) {
		// Make our own copy of rem we can modify
		rem := *rem

		api.GetProject = func(_ *gitlab.Client, ProjectID interface{}) (*gitlab.Project, error) {
			return nil, errors.New("error")
		}

		resolveNetwork(&rem)

		assert.Len(t, rem.network, 0)
	})

	t.Run("MaxRemotesForLookup limit", func(t *testing.T) {
		// Make our own copy of rem we can modify
		rem := *rem

		api.GetProject = mockAPIGetProject

		for i := 0; i < maxRemotesForLookup; i++ {
			rem.remotes = append(rem.remotes, rem.remotes[i])

		}
		// Make sure we have at least one more remote than the limit set from maxRemotesForLookup
		assert.Len(t, rem.remotes, maxRemotesForLookup+1)

		resolveNetwork(&rem)

		assert.Len(t, rem.network, maxRemotesForLookup)
		for i := range rem.network {
			assert.Equal(t, rem.remotes[i].Repo.FullName(), rem.network[i].PathWithNamespace)
		}
	})
}

func Test_BaseRepo(t *testing.T) {
	// Make it a function that must be called by each test so none of them overlap
	rem := func() ResolvedRemotes {
		rem := &ResolvedRemotes{
			remotes: Remotes{
				&Remote{
					Remote: &git.Remote{
						Name: "upstream",
					},
					Repo: NewWithHost("profclems", "glab", "gitlab.com"),
				},
			},
			apiClient: &gitlab.Client{},
			network: []gitlab.Project{
				{
					ID:                1,
					PathWithNamespace: "profclems/glab",
					HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				},
			},
		}
		return *rem
	}

	mockGitlabProject := func(i interface{}) gitlab.Project {
		p := &gitlab.Project{
			PathWithNamespace: fmt.Sprint(i),
			HTTPURLToRepo:     fmt.Sprintf("https://gitlab.com/%s", i),
		}
		return *p
	}

	// Override git.SetRemoteResolution so it doesn't mess with the user configs
	git.SetRemoteResolution = func(_, _ string) error {
		return nil
	}

	api.GetProject = func(_ *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
		p := mockGitlabProject(projectID)
		return &p, nil
	}

	t.Run("baseOverride", func(t *testing.T) {
		localRem := rem()
		localRem.baseOverride = NewWithHost("profclems", "glab", "gitlab.com")

		got, err := localRem.BaseRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.baseOverride.FullName(), got.FullName())
		assert.Equal(t, localRem.baseOverride.RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->base", func(t *testing.T) {
		localRem := rem()

		// Set a base resolution
		localRem.remotes[0].Resolved = "base"

		got, err := localRem.BaseRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->base:", func(t *testing.T) {
		localRem := rem()

		expectedResolution := NewWithHost("maxice8", "glab", "gitlab.com")

		// Set a base resolution
		localRem.remotes[0].Resolved = "base: gitlab.com/maxice8/glab"

		got, err := localRem.BaseRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, expectedResolution.FullName(), got.FullName())
		assert.Equal(t, expectedResolution.RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->base: (invalid)", func(t *testing.T) {
		localRem := rem()

		// Set a base resolution
		localRem.remotes[0].Resolved = "base:NotAnActualValidValue"

		got, err := localRem.BaseRepo(false)
		assert.Nil(t, got)
		assert.EqualError(t, err, `expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "NotAnActualValidValue"`)
	})

	t.Run("Resolved->backwards-compatibility", func(t *testing.T) {
		localRem := rem()

		expectedResolution := NewWithHost("maxice8", "glab", "gitlab.com")

		// Set a base resolution
		localRem.remotes[0].Resolved = "gitlab.com/maxice8/glab"

		got, err := localRem.BaseRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, expectedResolution.FullName(), got.FullName())
		assert.Equal(t, expectedResolution.RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->backwards-compatibility: (invalid)", func(t *testing.T) {
		localRem := rem()

		// Set a base resolution
		localRem.remotes[0].Resolved = "NotAnActualValidValue"

		got, err := localRem.BaseRepo(false)
		assert.Nil(t, got)
		assert.EqualError(t, err, `expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "NotAnActualValidValue"`)
	})

	t.Run("Prompt==false", func(t *testing.T) {
		localRem := rem()

		got, err := localRem.BaseRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network 1 repo", func(t *testing.T) {
		localRem := rem()

		// Prompt must be true otherwise we won't reach the code we want to test
		got, err := localRem.BaseRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, no remotes", func(t *testing.T) {
		localRem := rem()

		// Wipe out all remotes
		localRem.remotes = Remotes{}
		localRem.network = nil

		_, err := localRem.BaseRepo(true)
		assert.EqualError(t, err, "no GitLab Projects found from remotes")
	})

	t.Run("Consult the network, multiple projects, pick origin", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "base",
				Value: "maxice8/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.BaseRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, originRemote.Repo.FullName(), got.FullName())
		assert.Equal(t, originRemote.Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, multiple projects, pick upstream", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "base",
				Value: "profclems/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.BaseRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].Repo.FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, one forked project, get fork", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
			ForkedFromProject: &gitlab.ForkParent{
				ID:                1,
				HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				PathWithNamespace: "profclems/glab",
			},
		}

		localRem.remotes = Remotes{originRemote}
		localRem.network = []gitlab.Project{originNetwork}

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "base",
				Value: "maxice8/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.BaseRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, originRemote.Repo.FullName(), got.FullName())
		assert.Equal(t, originRemote.Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, one forked project, get upstream", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
			ForkedFromProject: &gitlab.ForkParent{
				ID:                1,
				HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				PathWithNamespace: "profclems/glab",
			},
		}

		localRem.remotes = Remotes{originRemote}
		localRem.network = []gitlab.Project{originNetwork}

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "base",
				Value: "profclems/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.BaseRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, "profclems/glab", got.FullName())
		assert.Equal(t, "gitlab.com", got.RepoHost())
	})

	t.Run("Consult the network, multiple projects, prompt fails", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "base",
				Value: errors.New("could not prompt"),
			},
		})

		got, err := localRem.BaseRepo(true)
		assert.Nil(t, got)
		assert.EqualError(t, err, "could not prompt")
	})
}

func Test_HeadRepo(t *testing.T) {
	// Make it a function that must be called by each test so none of them overlap
	rem := func() ResolvedRemotes {
		rem := &ResolvedRemotes{
			remotes: Remotes{
				&Remote{
					Remote: &git.Remote{
						Name: "upstream",
					},
					Repo: NewWithHost("profclems", "glab", "gitlab.com"),
				},
			},
			apiClient: &gitlab.Client{},
			network: []gitlab.Project{
				{
					ID:                1,
					PathWithNamespace: "profclems/glab",
					HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				},
			},
		}
		return *rem
	}

	mockGitlabProject := func(i interface{}) gitlab.Project {
		p := &gitlab.Project{
			PathWithNamespace: fmt.Sprint(i),
			HTTPURLToRepo:     fmt.Sprintf("https://gitlab.com/%s", i),
		}
		return *p
	}

	// Override git.SetRemoteResolution so it doesn't mess with the user configs
	git.SetRemoteResolution = func(_, _ string) error {
		return nil
	}

	api.GetProject = func(_ *gitlab.Client, projectID interface{}) (*gitlab.Project, error) {
		p := mockGitlabProject(projectID)
		return &p, nil
	}

	t.Run("baseOverride", func(t *testing.T) {
		localRem := rem()
		localRem.baseOverride = NewWithHost("profclems", "glab", "gitlab.com")

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.baseOverride.FullName(), got.FullName())
		assert.Equal(t, localRem.baseOverride.RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->head", func(t *testing.T) {
		localRem := rem()

		// Set a head resolution
		localRem.remotes[0].Resolved = "head"

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->head:", func(t *testing.T) {
		localRem := rem()

		expectedResolution := NewWithHost("maxice8", "glab", "gitlab.com")

		// Set a base resolution
		localRem.remotes[0].Resolved = "head: gitlab.com/maxice8/glab"

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, expectedResolution.FullName(), got.FullName())
		assert.Equal(t, expectedResolution.RepoHost(), got.RepoHost())
	})

	t.Run("Resolved->head: (invalid)", func(t *testing.T) {
		localRem := rem()

		// Set a base resolution
		localRem.remotes[0].Resolved = "head:NotAnActualValidValue"

		got, err := localRem.HeadRepo(false)
		assert.Nil(t, got)
		assert.EqualError(t, err, `expected the "[HOST/]OWNER/[NAMESPACE/]REPO" format, got "NotAnActualValidValue"`)
	})

	t.Run("Prompt==false", func(t *testing.T) {
		localRem := rem()

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].FullName(), got.FullName())
	})

	t.Run("Consult the network 1 repo", func(t *testing.T) {
		localRem := rem()

		// Prompt must be true otherwise we won't reach the code we want to test
		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, got.FullName(), localRem.remotes[0].FullName())
	})

	t.Run("Consult the network, more than 1 forked repo, pick the fork", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
			ForkedFromProject: &gitlab.ForkParent{
				ID:                1,
				HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				PathWithNamespace: "profclems/glab",
			},
		}

		localRem.remotes = Remotes{originRemote}
		localRem.network = []gitlab.Project{originNetwork}

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, "maxice8/glab", got.FullName())
	})

	t.Run("Consult the network, more than 1 repo, pick the first", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		got, err := localRem.HeadRepo(false)
		assert.NoError(t, err)

		assert.Equal(t, "profclems/glab", got.FullName())
	})

	t.Run("Consult the network, no remotes", func(t *testing.T) {
		localRem := rem()

		// Wipe out all remotes
		localRem.remotes = Remotes{}
		localRem.network = nil

		_, err := localRem.HeadRepo(true)
		assert.EqualError(t, err, "no GitLab Projects found from remotes")
	})

	t.Run("Consult the network, multiple projects, pick origin", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "head",
				Value: "maxice8/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.HeadRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, originRemote.Repo.FullName(), got.FullName())
		assert.Equal(t, originRemote.Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, multiple projects, pick upstream", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "head",
				Value: "profclems/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.HeadRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, localRem.remotes[0].Repo.FullName(), got.FullName())
		assert.Equal(t, localRem.remotes[0].Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, one forked project, get fork", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
			ForkedFromProject: &gitlab.ForkParent{
				ID:                1,
				HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				PathWithNamespace: "profclems/glab",
			},
		}

		localRem.remotes = Remotes{originRemote}
		localRem.network = []gitlab.Project{originNetwork}

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "head",
				Value: "maxice8/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.HeadRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, originRemote.Repo.FullName(), got.FullName())
		assert.Equal(t, originRemote.Repo.RepoHost(), got.RepoHost())
	})

	t.Run("Consult the network, one forked project, get upstream", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
			ForkedFromProject: &gitlab.ForkParent{
				ID:                1,
				HTTPURLToRepo:     "https://gitlab.com/profclems/glab",
				PathWithNamespace: "profclems/glab",
			},
		}

		localRem.remotes = Remotes{originRemote}
		localRem.network = []gitlab.Project{originNetwork}

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "head",
				Value: "profclems/glab", // We expect to get `origin`
			},
		})

		got, err := localRem.HeadRepo(true)
		assert.NoError(t, err)

		assert.Equal(t, "profclems/glab", got.FullName())
		assert.Equal(t, "gitlab.com", got.RepoHost())
	})

	t.Run("Consult the network, multiple projects, prompt fails", func(t *testing.T) {
		localRem := rem()

		originRemote := &Remote{
			Remote: &git.Remote{Name: "origin"},
			Repo:   NewWithHost("maxice8", "glab", "gitlab.com"),
		}

		originNetwork := gitlab.Project{
			ID:                2,
			PathWithNamespace: "maxice8/glab",
			HTTPURLToRepo:     "https://gitlab.com/maxice8/glab",
		}

		localRem.remotes = append(localRem.remotes, originRemote)
		localRem.network = append(localRem.network, originNetwork)

		// Mock the prompt
		as, restoreAsk := prompt.InitAskStubber()
		defer restoreAsk()

		as.Stub([]*prompt.QuestionStub{
			{
				Name:  "head",
				Value: errors.New("could not prompt"),
			},
		})

		got, err := localRem.HeadRepo(true)
		assert.Nil(t, got)
		assert.EqualError(t, err, "could not prompt")
	})
}
