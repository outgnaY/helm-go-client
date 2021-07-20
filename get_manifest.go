package helmclient

import "helm.sh/helm/v3/pkg/action"

const (
	getManifestDefaultVersion = 0
)

type getManifestClient interface {
	GetManifest(name string) (string, error)
}

type getManifestClientImpl struct {
	cli *action.Get
	env *helmEnv
}

type GetManifestOption struct {
	f func(o *getManifestOptions)
}

type getManifestOptions struct {
	version int
}

func (o *getManifestOptions) apply(opts []GetManifestOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newGetManifestOptions(opts []GetManifestOption) *getManifestOptions {
	options := &getManifestOptions{version: getManifestDefaultVersion}
	options.apply(opts)
	return options
}

func GetManifestWithVersion(version int) GetManifestOption {
	return GetManifestOption{f: func(o *getManifestOptions) {
		o.version = version
	}}
}

func (c *getManifestClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *getManifestClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewGet(cfg)
	// copy args
	copyGetManifestClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyGetManifestClientOptions(oldCli *action.Get, newCli *action.Get) {
	newCli.Version = oldCli.Version
}

func newGetManifestClient(opts []GetManifestOption, env *helmEnv) (*getManifestClientImpl, error) {
	o := newGetManifestOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewGet(cfg)
	mergeGetManifestOptions(o, client)
	return &getManifestClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *getManifestClientImpl) GetManifest(name string) (string, error) {
	release, err := c.cli.Run(name)
	if err != nil {
		return "", err
	}
	return release.Manifest, nil
}

func mergeGetManifestOptions(o *getManifestOptions, cli *action.Get) {
	cli.Version = o.version
}
