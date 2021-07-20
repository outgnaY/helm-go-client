package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"fmt"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"os"
	"path/filepath"
)

type chartSaveClient interface {
	ChartSave(args []string) error
}

type chartSaveClientImpl struct {
	cli *chartSave
	env *helmEnv
}

type chartSave struct {
	registryClient *registry.Client
}

func (c *chartSaveClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartSaveClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartSave(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newChartSaveClient(env *helmEnv) (*chartSaveClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartSave(registryClient)
	return &chartSaveClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartSave(registryClient *registry.Client) *chartSave {
	return &chartSave{
		registryClient: registryClient,
	}
}

func (c *chartSave) Run(ch *chart.Chart, ref string) error {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return err
	}

	// If no tag is present, use the chart version
	if r.Tag == "" {
		r.Tag = ch.Metadata.Version
	}

	return c.registryClient.SaveChart(ch, r)
}

func (c *chartSaveClientImpl) ChartSave(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("chart save requires at least 2 arguments")
	}
	path := args[0]
	ref := args[1]

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	ch, err := loader.Load(path)
	if err != nil {
		return err
	}
	return c.cli.Run(ch, ref)
}
