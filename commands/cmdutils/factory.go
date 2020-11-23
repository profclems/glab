// Forked from https://github.com/cli/cli/blob/929e082c13909044e2585af292ae952c9ca6f25c/pkg/cmd/factory/default.go
package cmdutils

import (
	"fmt"

	"github.com/profclems/glab/internal/glinstance"

	"github.com/profclems/glab/internal/utils"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/glrepo"
	"github.com/profclems/glab/pkg/api"
	"github.com/xanzy/go-gitlab"
)

var (
	CachedConfig config.Config
	ConfigError  error
)

type Factory struct {
	HttpClient func() (*gitlab.Client, error)
	BaseRepo   func() (glrepo.Interface, error)
	Remotes    func() (glrepo.Remotes, error)
	Config     func() (config.Config, error)
	Branch     func() (string, error)
	IO         *utils.IOStreams
}

func (f *Factory) RepoOverride(repo string) error {
	f.BaseRepo = func() (glrepo.Interface, error) {
		return glrepo.FromFullName(repo)
	}
	newRepo, err := f.BaseRepo()
	if err != nil {
		return err
	}
	// Initialise new http client for new repo host
	cfg, err := f.Config()
	if err == nil {
		OverrideAPIProtocol(cfg, newRepo)
	}
	f.HttpClient = func() (*gitlab.Client, error) {
		return HttpClientFunc(newRepo.RepoHost(), cfg, false)
	}
	return nil
}

func HttpClientFunc(repoHost string, cfg config.Config, isGraphQL bool) (*gitlab.Client, error) {
	if repoHost == "" {
		repoHost = glinstance.OverridableDefault()
	}
	token, _ := cfg.Get(repoHost, "token")
	tlsVerify, _ := cfg.Get(repoHost, "skip_tls_verify")
	skipTlsVerify := tlsVerify == "true" || tlsVerify == "1"
	caCert, _ := cfg.Get(repoHost, "ca_cert")
	if caCert != "" {
		return api.InitWithCustomCA(repoHost, token, caCert, isGraphQL)
	}
	return api.Init(repoHost, token, skipTlsVerify, isGraphQL)
}

func remotesFunc() (glrepo.Remotes, error) {
	rr := &remoteResolver{
		readRemotes: git.Remotes,
		getConfig:   configFunc,
	}
	fn := rr.Resolver()
	return fn()
}

func configFunc() (config.Config, error) {
	if CachedConfig != nil || ConfigError != nil {
		return CachedConfig, ConfigError
	}
	CachedConfig, ConfigError = initConfig()
	return CachedConfig, ConfigError
}

func baseRepoFunc() (glrepo.Interface, error) {
	remotes, err := remotesFunc()
	if err != nil {
		return nil, err
	}
	return glrepo.FromURL(remotes[0].FetchURL)
}

// OverrideAPIProtocol sets api protocol for host to initialize http client
func OverrideAPIProtocol(cfg config.Config, repo glrepo.Interface) {
	api.Protocol, _ = cfg.Get(repo.RepoHost(), "api_protocol")
}

func HTTPClientFactory(f *Factory) {
	f.HttpClient = func() (*gitlab.Client, error) {
		cfg, err := configFunc()
		if err != nil {
			return nil, err
		}
		repo, err := baseRepoFunc()
		if err != nil {
			// use default hostname if remote resolver fails
			repo = glrepo.NewWithHost("", "", glinstance.OverridableDefault())
		}
		OverrideAPIProtocol(cfg, repo)
		return HttpClientFunc(repo.RepoHost(), cfg, false)
	}
}

func NewFactory() *Factory {

	return &Factory{
		Config:  configFunc,
		Remotes: remotesFunc,
		HttpClient: func() (*gitlab.Client, error) {
			// do not initialize httpclient since it may not be required by
			// some commands like version, help, etc...
			// It should be explicitly initialize with HTTPClientFactory()
			return nil, nil
		},
		BaseRepo: baseRepoFunc,
		Branch: func() (string, error) {
			currentBranch, err := git.CurrentBranch()
			if err != nil {
				return "", fmt.Errorf("could not determine current branch: %w", err)
			}
			return currentBranch, nil
		},
		IO: utils.InitIOStream(),
	}
}

func initConfig() (config.Config, error) {
	if err := config.MigrateOldConfig(); err != nil {
		return nil, err
	}
	return config.Init()
}
