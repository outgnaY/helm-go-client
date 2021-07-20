package helmclient

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"os"
	"sync"
)

var errNoRepositories = errors.New("no repositories found. You must add one before updating")

type repoUpdateClient interface {
	RepoUpdate() error
}

type repoUpdateClientImpl struct {
	repoUpdateOpts *repoUpdateOptions
	env            *helmEnv
}

type repoUpdateOptions struct {
	update    func([]*repo.ChartRepository, io.Writer)
	repoFile  string
	repoCache string
}

func newRepoUpdateOptions() *repoUpdateOptions {
	options := &repoUpdateOptions{
		update: updateCharts,
	}
	return options
}

func updateCharts(repos []*repo.ChartRepository, out io.Writer) {
	fmt.Fprintln(out, "Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Fprintf(out, "...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Fprintf(out, "...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Fprintln(out, "Update Complete. ⎈Happy Helming!⎈")
}

func (c *repoUpdateClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *repoUpdateClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newRepoUpdateClient(env *helmEnv) (*repoUpdateClientImpl, error) {
	o := newRepoUpdateOptions()
	return &repoUpdateClientImpl{
		repoUpdateOpts: o,
		env:            env,
	}, nil
}

func (o *repoUpdateOptions) run(out io.Writer, settings *cli.EnvSettings) error {
	f, err := repo.LoadFile(o.repoFile)
	switch {
	case isNotExist(err):
		return errNoRepositories
	case err != nil:
		return errors.Wrapf(err, "failed loading file: %s", o.repoFile)
	case len(f.Repositories) == 0:
		return errNoRepositories
	}

	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		if o.repoCache != "" {
			r.CachePath = o.repoCache
		}
		repos = append(repos, r)
	}

	o.update(repos, out)
	return nil
}

func (c *repoUpdateClientImpl) RepoUpdate() error {
	c.repoUpdateOpts.repoFile = c.env.settings.RepositoryConfig
	c.repoUpdateOpts.repoCache = c.env.settings.RepositoryCache
	return c.repoUpdateOpts.run(os.Stdout, c.env.settings)
}
