package cmdutils

import (
	"net/url"
	"testing"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/pkg/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MakeNowJust/heredoc"
)

func Test_remoteResolver(t *testing.T) {
	rr := &remoteResolver{
		readRemotes: func() (git.RemoteSet, error) {
			return git.RemoteSet{
				git.NewRemote("fork", "https://example.org/owner/fork.git"),
				git.NewRemote("origin", "https://gitlab.com/owner/repo.git"),
				git.NewRemote("upstream", "https://example.org/owner/repo.git"),
			}, nil
		},
		getConfig: func() (config.Config, error) {
			return config.NewFromString(heredoc.Doc(`
				hosts:
				  example.org:
				    oauth_token: OTOKEN
			`)), nil
		},
		urlTranslator: func(u *url.URL) *url.URL {
			return u
		},
	}

	resolver := rr.Resolver("")
	remotes, err := resolver()
	require.NoError(t, err)
	require.Equal(t, 2, len(remotes))

	assert.Equal(t, "upstream", remotes[0].Name)
	assert.Equal(t, "fork", remotes[1].Name)
}

func Test_remoteResolverOverride(t *testing.T) {
	rr := &remoteResolver{
		readRemotes: func() (git.RemoteSet, error) {
			return git.RemoteSet{
				git.NewRemote("fork", "https://example.org/ghe-owner/ghe-fork.git"),
				git.NewRemote("origin", "https://gitlab.com/owner/repo.git"),
				git.NewRemote("upstream", "https://example.org/ghe-owner/ghe-repo.git"),
			}, nil
		},
		getConfig: func() (config.Config, error) {
			return config.NewFromString(heredoc.Doc(`
				hosts:
				  example.org:
				    oauth_token: GHETOKEN
			`)), nil
		},
		urlTranslator: func(u *url.URL) *url.URL {
			return u
		},
	}

	resolver := rr.Resolver("gitlab.com")
	remotes, err := resolver()
	require.NoError(t, err)
	require.Equal(t, 1, len(remotes))

	assert.Equal(t, "origin", remotes[0].Name)
}
