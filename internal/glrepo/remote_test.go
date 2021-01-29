package glrepo

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/profclems/glab/pkg/git"
	"github.com/stretchr/testify/assert"
)

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindByName(t *testing.T) {
	list := Remotes{
		&Remote{Remote: &git.Remote{Name: "mona"}, Repo: New("monalisa", "myfork")},
		&Remote{Remote: &git.Remote{Name: "origin"}, Repo: New("monalisa", "octo-cat")},
		&Remote{Remote: &git.Remote{Name: "upstream"}, Repo: New("hubot", "tools")},
	}

	r, err := list.FindByName("upstream", "origin")
	eq(t, err, nil)
	eq(t, r.Name, "upstream")

	r, err = list.FindByName("nonexist", "*")
	eq(t, err, nil)
	eq(t, r.Name, "mona")

	_, err = list.FindByName("nonexist")
	eq(t, err, errors.New(`no GitLab remotes found`))
}

func TestTranslateRemotes(t *testing.T) {
	publicURL, _ := url.Parse("https://gitlab.com/monalisa/hello")
	originURL, _ := url.Parse("http://example.com/repo")
	upstreamURL, _ := url.Parse("https://gitlab.com/profclems/glab")

	gitRemotes := git.RemoteSet{
		&git.Remote{
			Name:     "origin",
			FetchURL: originURL,
		},
		&git.Remote{
			Name:     "public",
			FetchURL: publicURL,
		},
		&git.Remote{
			Name:    "upstream",
			PushURL: upstreamURL,
		},
	}

	identityURL := func(u *url.URL) *url.URL {
		return u
	}
	result := TranslateRemotes(gitRemotes, identityURL)

	if len(result) != 2 {
		t.Errorf("got %d results", len(result))
	}
	if result[0].Name != "public" {
		t.Errorf("got %q", result[0].Name)
	}
	if result[0].RepoName() != "hello" {
		t.Errorf("got %q", result[0].RepoName())
	}
	if result[1].Name != "upstream" {
		t.Errorf("got %q", result[1].Name)
	}
	if result[1].RepoName() != "glab" {
		t.Errorf("got %q", result[1].Name)
	}
}

func Test_remoteNameSortingScore(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output int
	}{
		{
			name:   "upstream",
			input:  "upstream",
			output: 3,
		},
		{
			name:   "gitlab",
			input:  "gitlab",
			output: 2,
		},
		{
			name:   "origin",
			input:  "origin",
			output: 1,
		},
		{
			name:   "else",
			input:  "anyOtherName",
			output: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := remoteNameSortScore(tC.input)
			assert.Equal(t, tC.output, got)
		})
	}
}

func Test_FindByRepo(t *testing.T) {
	r := Remotes{
		&Remote{
			Remote: &git.Remote{
				Name: "origin",
			},
			Repo: NewWithHost("profclems", "glab", "gitlab.com"),
		},
	}

	t.Run("success", func(t *testing.T) {
		got, err := r.FindByRepo("profclems", "glab")
		assert.NoError(t, err)

		assert.Equal(t, r[0].FullName(), got.FullName())
	})

	t.Run("fail/owner", func(t *testing.T) {
		got, err := r.FindByRepo("maxice8", "glab")
		assert.Nil(t, got)
		assert.EqualError(t, err, "no matching remote found")
	})

	t.Run("fail/project", func(t *testing.T) {
		got, err := r.FindByRepo("profclems", "balg")
		assert.Nil(t, got)
		assert.EqualError(t, err, "no matching remote found")
	})

	t.Run("fail/owner and project", func(t *testing.T) {
		got, err := r.FindByRepo("maxice8", "balg")
		assert.Nil(t, got)
		assert.EqualError(t, err, "no matching remote found")
	})
}

func Test_RepoFuncs(t *testing.T) {
	testCases := []struct {
		name          string
		input         []string
		wantHostname  string
		wantOwner     string
		wantGroup     string
		wantNamespace string
		wantName      string
		wantFullname  string
	}{
		{
			name:          "Simple",
			input:         []string{"profclems", "glab", "gitlab.com"},
			wantHostname:  "gitlab.com",
			wantNamespace: "profclems",
			wantOwner:     "profclems",
			wantName:      "glab",
			wantFullname:  "profclems/glab",
		},
		{
			name:          "group",
			input:         []string{"company/profclems", "glab", "gitlab.com"},
			wantHostname:  "gitlab.com",
			wantNamespace: "profclems",
			wantOwner:     "company/profclems",
			wantGroup:     "company",
			wantName:      "glab",
			wantFullname:  "company/profclems/glab",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := Remote{
				Repo: NewWithHost(tC.input[0], tC.input[1], tC.input[2]),
			}
			if tC.wantHostname != "" {
				assert.Equal(t, tC.wantHostname, got.RepoHost())
			}
			if tC.wantOwner != "" {
				assert.Equal(t, tC.wantOwner, got.RepoOwner())
			}
			if tC.wantGroup != "" {
				assert.Equal(t, tC.wantGroup, got.RepoGroup())
			}
			if tC.wantNamespace != "" {
				assert.Equal(t, tC.wantNamespace, got.RepoNamespace())
			}
			if tC.wantName != "" {
				assert.Equal(t, tC.wantName, got.RepoName())
			}
			if tC.wantFullname != "" {
				assert.Equal(t, tC.wantFullname, got.FullName())
			}
		})
	}
}

func Test_Swap(t *testing.T) {
	r := Remotes{
		&Remote{
			Remote: &git.Remote{
				Name: "origin",
			},
			Repo: NewWithHost("maxice8", "glab", "gitlab.com"),
		},
		&Remote{
			Remote: &git.Remote{
				Name: "upstream",
			},
			Repo: NewWithHost("profclems", "glab", "gitlab.com"),
		},
	}

	assert.Equal(t, "origin", r[0].Remote.Name)
	assert.Equal(t, "upstream", r[1].Remote.Name)

	assert.Equal(t, "maxice8/glab", r[0].Repo.FullName())
	assert.Equal(t, "profclems/glab", r[1].Repo.FullName())

	r.Swap(0, 1)

	assert.Equal(t, "upstream", r[0].Remote.Name)
	assert.Equal(t, "origin", r[1].Remote.Name)

	assert.Equal(t, "profclems/glab", r[0].Repo.FullName())
	assert.Equal(t, "maxice8/glab", r[1].Repo.FullName())
}

func Test_Less(t *testing.T) {
	r := Remotes{
		&Remote{
			Remote: &git.Remote{
				Name: "else",
			},
			Repo: NewWithHost("somebody", "glab", "gitlab.com"),
		},
		&Remote{
			Remote: &git.Remote{
				Name: "origin",
			},
			Repo: NewWithHost("maxice8", "glab", "gitlab.com"),
		},
		&Remote{
			Remote: &git.Remote{
				Name: "gitlab",
			},
			Repo: NewWithHost("profclems", "glab", "gitlab.com"),
		},
		&Remote{
			Remote: &git.Remote{
				Name: "upstream",
			},
			Repo: NewWithHost("profclems", "glab", "gitlab.com"),
		},
	}

	assert.True(t, r.Less(3, 2))
	assert.True(t, r.Less(3, 1))
	assert.True(t, r.Less(3, 0))
	assert.True(t, r.Less(2, 1))
	assert.True(t, r.Less(2, 0))
	assert.True(t, r.Less(1, 0))

	assert.False(t, r.Less(0, 1))
	assert.False(t, r.Less(0, 2))
	assert.False(t, r.Less(0, 3))
	assert.False(t, r.Less(1, 2))
	assert.False(t, r.Less(1, 3))
	assert.False(t, r.Less(2, 3))
}
