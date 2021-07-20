package helmclient

import (
	"github.com/outgnaY/helm-go-client/export"
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

type chartListClient interface {
	ChartList() ([]*export.ChartInfo, error)
}

type chartListClientImpl struct {
	cli *chartList
	env *helmEnv
}

type chartList struct {
	registryClient *registry.Client
}

func (c *chartListClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartListClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartList(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newChartListClient(env *helmEnv) (*chartListClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartList(registryClient)
	return &chartListClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartList(registryClient *registry.Client) *chartList {
	return &chartList{
		registryClient: registryClient,
	}
}

func (c *chartList) Run() ([]*export.ChartInfo, error) {
	return c.registryClient.GetChartInfos()
}

func (c *chartListClientImpl) ChartList() ([]*export.ChartInfo, error) {
	return c.cli.Run()
}
