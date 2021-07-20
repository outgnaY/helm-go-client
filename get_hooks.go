package helmclient

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

const (
	getHooksDefaultVersion = 0
)

type getHooksClient interface {
	GetHooks(name string) ([]*release.Hook, error)
}

type getHooksClientImpl struct {
	cli *action.Get
	env *helmEnv
}

type GetHooksOption struct {
	f func(o *getHooksOptions)
}

type getHooksOptions struct {
	version int
}

func (o *getHooksOptions) apply(opts []GetHooksOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newGetHooksOptions(opts []GetHooksOption) *getHooksOptions {
	options := &getHooksOptions{version: getHooksDefaultVersion}
	options.apply(opts)
	return options
}

func GetHooksWithVersion(version int) GetHooksOption {
	return GetHooksOption{f: func(o *getHooksOptions) {
		o.version = version
	}}
}

func (c *getHooksClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *getHooksClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewGet(cfg)
	// copy args
	copyGetHooksClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyGetHooksClientOptions(oldCli *action.Get, newCli *action.Get) {
	newCli.Version = oldCli.Version
}

func newGetHooksClient(opts []GetHooksOption, env *helmEnv) (*getHooksClientImpl, error) {
	o := newGetHooksOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewGet(cfg)
	mergeGetHooksOptions(o, client)
	return &getHooksClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *getHooksClientImpl) GetHooks(name string) ([]*release.Hook, error) {
	release, err := c.cli.Run(name)
	if err != nil {
		return nil, err
	}
	return release.Hooks, nil
}

func mergeGetHooksOptions(o *getHooksOptions, cli *action.Get) {
	cli.Version = o.version
}
