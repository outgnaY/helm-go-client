package helmclient

import "helm.sh/helm/v3/pkg/action"

const (
	getNotesDefaultVersion = 0
)

type getNotesClient interface {
	GetNotes(name string) (string, error)
}

type getNotesClientImpl struct {
	cli *action.Get
	env *helmEnv
}

type GetNotesOption struct {
	f func(o *getNotesOptions)
}

type getNotesOptions struct {
	version int
}

func (o *getNotesOptions) apply(opts []GetNotesOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newGetNotesOptions(opts []GetNotesOption) *getNotesOptions {
	options := &getNotesOptions{version: getNotesDefaultVersion}
	options.apply(opts)
	return options
}

func GetNotesWithVersion(version int) GetNotesOption {
	return GetNotesOption{f: func(o *getNotesOptions) {
		o.version = version
	}}
}

func (c *getNotesClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *getNotesClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewGet(cfg)
	// copy args
	copyGetNotesClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyGetNotesClientOptions(oldCli *action.Get, newCli *action.Get) {
	newCli.Version = oldCli.Version
}

func newGetNotesClient(opts []GetNotesOption, env *helmEnv) (*getNotesClientImpl, error) {
	o := newGetNotesOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewGet(cfg)
	mergeGetNotesOptions(o, client)
	return &getNotesClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *getNotesClientImpl) GetNotes(name string) (string, error) {
	release, err := c.cli.Run(name)
	if err != nil {
		return "", err
	}
	return release.Info.Notes, nil
}

func mergeGetNotesOptions(o *getNotesOptions, cli *action.Get) {
	cli.Version = o.version
}
