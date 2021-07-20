package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

type chartPushClient interface {
	ChartPush(ref string) error
}

type chartPushClientImpl struct {
	cli *chartPush
	env *helmEnv
}

type chartPush struct {
	registryClient *registry.Client
}

func (c *chartPushClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartPushClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartPush(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newChartPushClient(env *helmEnv) (*chartPushClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartPush(registryClient)
	return &chartPushClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartPush(registryClient *registry.Client) *chartPush {
	return &chartPush{
		registryClient: registryClient,
	}
}

func (c *chartPush) Run(ref string) error {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return err
	}
	return c.registryClient.PushChart(r)
}

func (c *chartPushClientImpl) ChartPush(ref string) error {
	return c.cli.Run(ref)
}
