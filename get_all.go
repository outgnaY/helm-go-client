package helmclient

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

const (
	getAllDefaultVersion = 0
)

type getAllClient interface {
	GetAll(name string) (*release.Release, error)
}

type getAllClientImpl struct {
	cli *action.Get
	env *helmEnv
}

type GetAllOption struct {
	f func(o *getAllOptions)
}

type getAllOptions struct {
	version int
}

func (o *getAllOptions) apply(opts []GetAllOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newGetAllOptions(opts []GetAllOption) *getAllOptions {
	options := &getAllOptions{
		version: getAllDefaultVersion,
	}
	options.apply(opts)
	return options
}

func GetAllWithVersion(version int) GetAllOption {
	return GetAllOption{f: func(o *getAllOptions) {
		o.version = version
	}}
}

func (c *getAllClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *getAllClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewGet(cfg)
	// copy args
	copyGetAllClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyGetAllClientOptions(oldCli *action.Get, newCli *action.Get) {
	newCli.Version = oldCli.Version
}

func newGetAllClient(opts []GetAllOption, env *helmEnv) (*getAllClientImpl, error) {
	o := newGetAllOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewGet(cfg)
	mergeGetAllOptions(o, client)
	return &getAllClientImpl{cli: client}, nil
}

func (c *getAllClientImpl) GetAll(name string) (*release.Release, error) {
	return c.cli.Run(name)
}

func mergeGetAllOptions(o *getAllOptions, cli *action.Get) {
	cli.Version = o.version
}
