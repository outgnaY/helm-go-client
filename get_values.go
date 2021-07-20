package helmclient

import (
	"helm.sh/helm/v3/pkg/action"
)

const (
	getValuesDefaultVersion   = 0
	getValuesDefaultAllValues = false
)

type getValuesClient interface {
	GetValues(name string) (map[string]interface{}, error)
}

type getValuesClientImpl struct {
	cli *action.GetValues
	env *helmEnv
}

type GetValuesOption struct {
	f func(o *getValuesOptions)
}

type getValuesOptions struct {
	version   int
	allValues bool
}

func (o *getValuesOptions) apply(opts []GetValuesOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newGetValuesOptions(opts []GetValuesOption) *getValuesOptions {
	options := &getValuesOptions{
		version:   getValuesDefaultVersion,
		allValues: getValuesDefaultAllValues,
	}
	options.apply(opts)
	return options
}

func GetValuesWithVersion(version int) GetValuesOption {
	return GetValuesOption{f: func(o *getValuesOptions) {
		o.version = version
	}}
}

func GetValuesWithAllValues(allValues bool) GetValuesOption {
	return GetValuesOption{f: func(o *getValuesOptions) {
		o.allValues = allValues
	}}
}

func (c *getValuesClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *getValuesClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewGetValues(cfg)
	// copy args
	copyGetValuesClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyGetValuesClientOptions(oldCli *action.GetValues, newCli *action.GetValues) {
	newCli.Version = oldCli.Version
	newCli.AllValues = oldCli.AllValues
}

func newGetValuesClient(opts []GetValuesOption, env *helmEnv) (*getValuesClientImpl, error) {
	o := newGetValuesOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewGetValues(cfg)
	mergeGetValuesOptions(o, client)
	return &getValuesClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *getValuesClientImpl) GetValues(name string) (map[string]interface{}, error) {
	return c.cli.Run(name)
}

func mergeGetValuesOptions(o *getValuesOptions, cli *action.GetValues) {
	cli.Version = o.version
	cli.AllValues = o.allValues
}
