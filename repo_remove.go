package helmclient

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"os"
	"path/filepath"
)

type repoRemoveClient interface {
	RepoRemove(args []string) error
}

type repoRemoveClientImpl struct {
	repoRemoveOpts *repoRemoveOptions
	env            *helmEnv
}

type repoRemoveOptions struct {
	names     []string
	repoFile  string
	repoCache string
}

func newRepoRemoveOptions() *repoRemoveOptions {
	options := &repoRemoveOptions{}
	return options
}

func (c *repoRemoveClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *repoRemoveClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newRepoRemoveClient(env *helmEnv) (*repoRemoveClientImpl, error) {
	o := newRepoRemoveOptions()
	return &repoRemoveClientImpl{
		repoRemoveOpts: o,
		env:            env,
	}, nil
}

func (o *repoRemoveOptions) run(out io.Writer) error {
	r, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(r.Repositories) == 0 {
		return errors.New("no repositories configured")
	}

	for _, name := range o.names {
		if !r.Remove(name) {
			return errors.Errorf("no repo named %q found", name)
		}
		if err := r.WriteFile(o.repoFile, 0644); err != nil {
			return err
		}

		if err := removeRepoCache(o.repoCache, name); err != nil {
			return err
		}
		fmt.Fprintf(out, "%q has been removed from your repositories\n", name)
	}

	return nil
}

func removeRepoCache(root, name string) error {
	idx := filepath.Join(root, helmpath.CacheChartsFile(name))
	if _, err := os.Stat(idx); err == nil {
		os.Remove(idx)
	}

	idx = filepath.Join(root, helmpath.CacheIndexFile(name))
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "can't remove index file %s", idx)
	}
	return os.Remove(idx)
}

func (c *repoRemoveClientImpl) RepoRemove(args []string) error {
	c.repoRemoveOpts.repoFile = c.env.settings.RepositoryConfig
	c.repoRemoveOpts.repoCache = c.env.settings.RepositoryCache
	c.repoRemoveOpts.names = args
	return c.repoRemoveOpts.run(os.Stdout)
}
