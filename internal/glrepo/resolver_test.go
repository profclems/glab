package glrepo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/pkg/api"
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
				assert.Equal(t, got.Name, tC.output.Name)
				assert.Equal(t, got.Remote.Name, tC.output.Remote.Name)
				assert.Equal(t, got.Repo.FullName(), tC.output.Repo.FullName())
				assert.Equal(t, got.Repo.RepoHost(), tC.output.Repo.RepoHost())
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

		assert.Equal(t, r.apiClient, rem.apiClient)

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

		assert.Equal(t, r.baseOverride.FullName(), expectedBaseOverride.FullName())
		assert.Equal(t, r.baseOverride.RepoHost(), expectedBaseOverride.RepoHost())

		assert.Equal(t, r.apiClient, rem.apiClient)

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
		assert.EqualError(t, err, "expected the \"[HOST/]OWNER/[NAMESPACE/]REPO\" format, got \"badValue\"")

		assert.Equal(t, r.apiClient, rem.apiClient)

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
	api.GetProject = func(_ *gitlab.Client, ProjectID interface{}) (*gitlab.Project, error) {
		proj := &gitlab.Project{
			PathWithNamespace: fmt.Sprint(ProjectID),
		}
		return proj, nil
	}

	t.Run("simple", func(t *testing.T) {
		// Make our own copy of rem we can modify
		rem := *rem

		err := resolveNetwork(&rem)
		if err != nil {
			t.Errorf("resolveNetwork() unexpected error = %s", err)
		}

		assert.Len(t, rem.network, len(rem.remotes))
		for i := range rem.network {
			assert.Equal(t, rem.network[i].PathWithNamespace, rem.remotes[i].Repo.FullName())
		}
	})

	t.Run("API call failed", func(t *testing.T) {
		// Make our own copy of rem we can modify
		rem := *rem

		api.GetProject = func(_ *gitlab.Client, ProjectID interface{}) (*gitlab.Project, error) {
			return nil, errors.New("error")
		}

		err := resolveNetwork(&rem)
		if err != nil {
			t.Errorf("resolveNetwork() unexpected error = %s", err)
		}

		assert.Len(t, rem.network, 0)
	})
}
