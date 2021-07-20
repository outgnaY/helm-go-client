package helmclient

import (
	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"path/filepath"
	"strings"
)

const (
	searchMaxScore            = 25
	searchRepoDefaultVersions = false
	searchRepoDefaultRegexp   = false
	searchRepoDefaultDevel    = false
	searchRepoDefaultVersion  = ""
)

type searchRepoClient interface {
	SearchRepo(args []string) ([]*search.Result, error)
}

type searchRepoClientImpl struct {
	searchRepoOpts *searchRepoOptions
	env            *helmEnv
}

type SearchRepoOption struct {
	f func(o *searchRepoOptions)
}

type searchRepoOptions struct {
	versions     bool
	regexp       bool
	devel        bool
	version      string
	repoFile     string
	repoCacheDir string
}

func (o *searchRepoOptions) apply(opts []SearchRepoOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newSearchRepoOptions(opts []SearchRepoOption) *searchRepoOptions {
	options := &searchRepoOptions{
		versions: searchRepoDefaultVersions,
		regexp:   searchRepoDefaultRegexp,
		devel:    searchRepoDefaultDevel,
		version:  searchRepoDefaultVersion,
	}
	options.apply(opts)
	return options
}

func SearchRepoWithVersions(versions bool) SearchRepoOption {
	return SearchRepoOption{f: func(o *searchRepoOptions) {
		o.versions = versions
	}}
}

func SearchRepoWithRegexp(regexp bool) SearchRepoOption {
	return SearchRepoOption{f: func(o *searchRepoOptions) {
		o.regexp = regexp
	}}
}

func SearchRepoWithDevel(devel bool) SearchRepoOption {
	return SearchRepoOption{f: func(o *searchRepoOptions) {
		o.devel = devel
	}}
}

func SearchRepoWithVersion(version string) SearchRepoOption {
	return SearchRepoOption{f: func(o *searchRepoOptions) {
		o.version = version
	}}
}

func (c *searchRepoClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *searchRepoClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newSearchRepoClient(opts []SearchRepoOption, env *helmEnv) (*searchRepoClientImpl, error) {
	o := newSearchRepoOptions(opts)
	return &searchRepoClientImpl{
		searchRepoOpts: o,
		env:            env,
	}, nil
}

func (o *searchRepoOptions) run(args []string) ([]*search.Result, error) {
	o.setupSearchedVersion()
	index, err := o.buildIndex()
	if err != nil {
		return nil, err
	}
	var res []*search.Result
	if len(args) == 0 {
		res = index.All()
	} else {
		q := strings.Join(args, " ")
		res, err = index.Search(q, searchMaxScore, o.regexp)
		if err != nil {
			return nil, err
		}
	}

	search.SortScore(res)
	return o.applyConstraint(res)
}

func (c *searchRepoClientImpl) SearchRepo(args []string) ([]*search.Result, error) {
	c.searchRepoOpts.repoFile = c.env.settings.RepositoryConfig
	c.searchRepoOpts.repoCacheDir = c.env.settings.RepositoryCache
	return c.searchRepoOpts.run(args)
}

func (o *searchRepoOptions) setupSearchedVersion() {
	debug("Original chart version: %q", o.version)

	if o.version != "" {
		return
	}

	if o.devel { // search for releases and prereleases (alpha, beta, and release candidate releases).
		debug("setting version to >0.0.0-0")
		o.version = ">0.0.0-0"
	} else { // search only for stable releases, prerelease versions will be skip
		debug("setting version to >0.0.0")
		o.version = ">0.0.0"
	}
}

func (o *searchRepoOptions) applyConstraint(res []*search.Result) ([]*search.Result, error) {
	if o.version == "" {
		return res, nil
	}

	constraint, err := semver.NewConstraint(o.version)
	if err != nil {
		return res, errors.Wrap(err, "an invalid version/constraint format")
	}

	data := res[:0]
	foundNames := map[string]bool{}
	for _, r := range res {
		// if not returning all versions and already have found a result,
		// you're done!
		if !o.versions && foundNames[r.Name] {
			continue
		}
		v, err := semver.NewVersion(r.Chart.Version)
		if err != nil {
			continue
		}
		if constraint.Check(v) {
			data = append(data, r)
			foundNames[r.Name] = true
		}
	}

	return data, nil
}

func (o *searchRepoOptions) buildIndex() (*search.Index, error) {
	// Load the repositories.yaml
	rf, err := repo.LoadFile(o.repoFile)
	if isNotExist(err) || len(rf.Repositories) == 0 {
		return nil, errors.New("no repositories configured")
	}

	i := search.NewIndex()
	for _, re := range rf.Repositories {
		n := re.Name
		f := filepath.Join(o.repoCacheDir, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			warning("Repo %q is corrupt or missing. Try 'helm repo update'.", n)
			warning("%s", err)
			continue
		}

		i.AddRepo(n, ind, o.versions || len(o.version) > 0)
	}
	return i, nil
}
