package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

type chartRemoveClient interface {
	ChartRemove(ref string) error
}

type chartRemoveClientImpl struct {
	cli *chartRemove
	env *helmEnv
}

type chartRemove struct {
	registryClient *registry.Client
}

func (c *chartRemoveClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartRemoveClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartRemove(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newChartRemoveClient(env *helmEnv) (*chartRemoveClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartRemove(registryClient)
	return &chartRemoveClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartRemove(registryClient *registry.Client) *chartRemove {
	return &chartRemove{
		registryClient: registryClient,
	}
}

func (c *chartRemove) Run(ref string) error {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return err
	}
	return c.registryClient.RemoveChart(r)
}

func (c *chartRemoveClientImpl) ChartRemove(ref string) error {
	return c.cli.Run(ref)
}
