package glrepo

import (
	"testing"

	"github.com/profclems/glab/internal/git"
	"github.com/stretchr/testify/assert"
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
