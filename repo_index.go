package helmclient

import (
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
)

const (
	repoIndexDefaultUrl   = ""
	repoIndexDefaultMerge = ""
)

type repoIndexClient interface {
	RepoIndex(dir string) error
}

type repoIndexClientImpl struct {
	repoIndexOpts *repoIndexOptions
	env           *helmEnv
}

type RepoIndexOption struct {
	f func(o *repoIndexOptions)
}

type repoIndexOptions struct {
	dir   string
	url   string
	merge string
}

func (o *repoIndexOptions) apply(opts []RepoIndexOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newRepoIndexOptions(opts []RepoIndexOption) *repoIndexOptions {
	options := &repoIndexOptions{
		url:   repoIndexDefaultUrl,
		merge: repoIndexDefaultMerge,
	}
	options.apply(opts)
	return options
}

func RepoIndexWithUrl(url string) RepoIndexOption {
	return RepoIndexOption{f: func(o *repoIndexOptions) {
		o.url = url
	}}
}

func RepoIndexWithMerge(merge string) RepoIndexOption {
	return RepoIndexOption{f: func(o *repoIndexOptions) {
		o.merge = merge
	}}
}

func (c *repoIndexClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *repoIndexClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newRepoIndexClient(opts []RepoIndexOption, env *helmEnv) (*repoIndexClientImpl, error) {
	o := newRepoIndexOptions(opts)
	return &repoIndexClientImpl{
		repoIndexOpts: o,
		env:           env,
	}, nil
}

func (o *repoIndexOptions) run() error {
	path, err := filepath.Abs(o.dir)
	if err != nil {
		return err
	}

	return index(path, o.url, o.merge)
}

func index(dir, url, mergeTo string) error {
	out := filepath.Join(dir, "index.yaml")

	i, err := repo.IndexDirectory(dir, url)
	if err != nil {
		return err
	}
	if mergeTo != "" {
		// if index.yaml is missing then create an empty one to merge into
		var i2 *repo.IndexFile
		if _, err := os.Stat(mergeTo); os.IsNotExist(err) {
			i2 = repo.NewIndexFile()
			i2.WriteFile(mergeTo, 0644)
		} else {
			i2, err = repo.LoadIndexFile(mergeTo)
			if err != nil {
				return errors.Wrap(err, "merge failed")
			}
		}
		i.Merge(i2)
	}
	i.SortEntries()
	return i.WriteFile(out, 0644)
}

func (c *repoIndexClientImpl) RepoIndex(dir string) error {
	c.repoIndexOpts.dir = dir
	return c.repoIndexOpts.run()
}
