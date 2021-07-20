package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

type chartPullClient interface {
	ChartPull(ref string) error
}

type chartPullClientImpl struct {
	cli *chartPull
	env *helmEnv
}

type chartPull struct {
	registryClient *registry.Client
}

func (c *chartPullClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartPullClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartPull(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newChartPullClient(env *helmEnv) (*chartPullClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartPull(registryClient)
	return &chartPullClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartPull(registryClient *registry.Client) *chartPull {
	return &chartPull{
		registryClient: registryClient,
	}
}

func (c *chartPull) Run(ref string) error {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return err
	}
	return c.registryClient.PullChartToCache(r)
}

func (c *chartPullClientImpl) ChartPull(ref string) error {
	return c.cli.Run(ref)
}
