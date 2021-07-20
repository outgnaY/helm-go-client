package helmclient

import (
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/repo"
)

type repoListClient interface {
	RepoList() ([]*repo.Entry, error)
}

type repoListClientImpl struct {
	env *helmEnv
}

func (c *repoListClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *repoListClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newRepoListClient(env *helmEnv) (*repoListClientImpl, error) {
	return &repoListClientImpl{
		env: env,
	}, nil
}

func (c *repoListClientImpl) RepoList() ([]*repo.Entry, error) {
	f, err := repo.LoadFile(c.env.settings.RepositoryConfig)
	if isNotExist(err) {
		return nil, errors.New("no repositories to show")
	}
	return f.Repositories, nil
}
