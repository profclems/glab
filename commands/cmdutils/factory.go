// Forked from https://github.com/cli/cli/blob/929e082c13909044e2585af292ae952c9ca6f25c/pkg/cmd/factory/default.go
package cmdutils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	gLab "github.com/profclems/glab/internal/gitlab"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/xanzy/go-gitlab"
)

type Factory struct {
	HttpClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Remotes    func() (glrepo.Remotes, error)
	Config     func() (config.Config, error)
	Branch     func() (string, error)
}

func (f *Factory) NewClient(repo string) (*Factory, error) {
	f.BaseRepo = func() (glrepo.Interface, error) {
		return glrepo.FromFullName(repo)
	}
	newRepo, err := f.BaseRepo()
	if err != nil {
		return nil, err
	}
	cfg, _ := f.Config()
	f.HttpClient = func() (*gitlab.Client, error) {
		return httpClientFunc(cfg, newRepo)
	}
	return f, nil
}

func httpClientFunc(cfg config.Config, repo glrepo.Interface) (*gitlab.Client, error) {
	token, _ := cfg.Get(repo.RepoHost(), "token")
	tlsVerify, _ := cfg.Get(repo.RepoHost(), "skip_tls_verify")
	skipTlsVerify, _ := strconv.ParseBool(tlsVerify)
	caCert, _ := cfg.Get(repo.RepoHost(), "ca_cert")
	if caCert != "" {
		return gLab.InitWithCustomCA(repo.RepoHost(), token, caCert)
	}
	return gLab.Init(repo.RepoHost(), token, skipTlsVerify)
}

func New(cachedConfig config.Config, configError error) *Factory {

	configFunc := func() (config.Config, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.ParseDefaultConfig()
		if errors.Is(configError, os.ErrNotExist) {
			cachedConfig = config.NewBlankConfig()
			configError = nil
		}
		return cachedConfig, configError
	}

	rr := &remoteResolver{
		readRemotes: git.Remotes,
		getConfig:   configFunc,
	}

	remotesFunc := rr.Resolver()

	baseRepoFunc := func() (glrepo.Interface, error) {
		remotes, err := remotesFunc()
		if err != nil {
			return nil, err
		}
		return glrepo.FromURL(remotes[0].FetchURL)
	}
	return &Factory{
		Config:  configFunc,
		Remotes: remotesFunc,
		HttpClient: func() (*gitlab.Client, error) {
			cfg, err := configFunc()
			if err != nil {
				return nil, err
			}
			repo, err := baseRepoFunc()
			if err != nil {
				return nil, err
			}
			return httpClientFunc(cfg, repo)
		},
		BaseRepo: func() (glrepo.Interface, error) {
			remotes, err := remotesFunc()
			if err != nil {
				return nil, err
			}
			return glrepo.FromURL(remotes[0].FetchURL)
		},
		Branch: func() (string, error) {
			currentBranch, err := git.CurrentBranch()
			if err != nil {
				return "", fmt.Errorf("could not determine current branch: %w", err)
			}
			return currentBranch, nil
		},
	}
}
