package helmclient

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"os"
)

const (
	listDefaultAll           = false
	listDefaultAllNamespaces = false
	listDefaultByDate        = false
	listDefaultSortReverse   = false
	listDefaultLimit         = 256
	listDefaultOffset        = 0
	listDefaultFilter        = ""
	listDefaultShort         = false
	listDefaultTimeFormat    = ""
	listDefaultUninstalled   = false
	listDefaultSuperseded    = false
	listDefaultUninstalling  = false
	listDefaultDeployed      = false
	listDefaultFailed        = false
	listDefaultPending       = false
	listDefaultSelector      = ""
)

type listClient interface {
	List() ([]*release.Release, error)
}

type listClientImpl struct {
	cli *action.List
	env *helmEnv
}

type ListOption struct {
	f func(o *listOptions)
}

type listOptions struct {
	all           bool
	allNamespaces bool
	byDate        bool
	sortReverse   bool
	limit         int
	offset        int
	filter        string
	short         bool
	timeFormat    string
	uninstalled   bool
	superseded    bool
	uninstalling  bool
	deployed      bool
	failed        bool
	pending       bool
	selector      string
}

func (o *listOptions) apply(opts []ListOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newListOptions(opts []ListOption) *listOptions {
	options := &listOptions{
		all:           listDefaultAll,
		allNamespaces: listDefaultAllNamespaces,
		byDate:        listDefaultByDate,
		sortReverse:   listDefaultSortReverse,
		limit:         listDefaultLimit,
		offset:        listDefaultOffset,
		filter:        listDefaultFilter,
		short:         listDefaultShort,
		timeFormat:    listDefaultTimeFormat,
		uninstalled:   listDefaultUninstalled,
		superseded:    listDefaultSuperseded,
		uninstalling:  listDefaultUninstalling,
		deployed:      listDefaultDeployed,
		failed:        listDefaultFailed,
		pending:       listDefaultPending,
		selector:      listDefaultSelector,
	}
	options.apply(opts)
	return options
}

func ListWithAll(all bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.all = all
	}}
}

func ListWithAllNamespaces(allNamespaces bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.allNamespaces = allNamespaces
	}}
}

func ListWithByDate(byDate bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.byDate = byDate
	}}
}

func ListWithSortReverse(sortReverse bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.sortReverse = sortReverse
	}}
}

func ListWithLimit(limit int) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.limit = limit
	}}
}

func ListWithOffset(offset int) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.offset = offset
	}}
}

func ListWithFilter(filter string) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.filter = filter
	}}
}

func ListWithShort(short bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.short = short
	}}
}

func ListWithTimeFormat(timeFormat string) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.timeFormat = timeFormat
	}}
}

func ListWithUninstalled(uninstalled bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.uninstalled = uninstalled
	}}
}

func ListWithSuperseded(superseded bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.superseded = superseded
	}}
}

func ListWithUninstalling(uninstalling bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.uninstalling = uninstalling
	}}
}

func ListWithDeployed(deployed bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.deployed = deployed
	}}
}

func ListWithFailed(failed bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.failed = failed
	}}
}

func ListWithPending(pending bool) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.pending = pending
	}}
}

func ListWithSelector(selector string) ListOption {
	return ListOption{f: func(o *listOptions) {
		o.selector = selector
	}}
}

func (c *listClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *listClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewList(cfg)
	// copy args
	copyListClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyListClientOptions(oldCli *action.List, newCli *action.List) {
	newCli.All = oldCli.All
	newCli.AllNamespaces = oldCli.AllNamespaces
	newCli.ByDate = oldCli.ByDate
	newCli.SortReverse = oldCli.SortReverse
	newCli.Limit = oldCli.Limit
	newCli.Offset = oldCli.Offset
	newCli.Filter = oldCli.Filter
	newCli.Short = oldCli.Short
	newCli.TimeFormat = oldCli.TimeFormat
	newCli.Uninstalled = oldCli.Uninstalled
	newCli.Superseded = oldCli.Superseded
	newCli.Uninstalling = oldCli.Uninstalling
	newCli.Deployed = oldCli.Deployed
	newCli.Failed = oldCli.Failed
	newCli.Pending = oldCli.Pending
	newCli.Selector = oldCli.Selector
}

func newListClient(opts []ListOption, env *helmEnv) (*listClientImpl, error) {
	o := newListOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewList(cfg)
	mergeListOptions(o, client)
	if client.AllNamespaces {
		if err := cfg.Init(env.clientGetter, "", os.Getenv("HELM_DRIVER"), debug); err != nil {
			return nil, err
		}
	}
	client.SetStateMask()
	return &listClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *listClientImpl) List() ([]*release.Release, error) {
	// concurrent access is not allowed!!!
	return c.cli.Run()
}

func mergeListOptions(o *listOptions, cli *action.List) {
	cli.All = o.all
	cli.AllNamespaces = o.allNamespaces
	cli.ByDate = o.byDate
	cli.SortReverse = o.sortReverse
	cli.Limit = o.limit
	cli.Offset = o.offset
	cli.Filter = o.filter
	cli.Short = o.short
	cli.TimeFormat = o.timeFormat
	cli.Uninstalled = o.uninstalled
	cli.Superseded = o.superseded
	cli.Uninstalling = o.uninstalling
	cli.Deployed = o.deployed
	cli.Failed = o.failed
	cli.Pending = o.pending
	cli.Selector = o.selector
}
