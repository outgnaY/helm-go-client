package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"fmt"
	"helm.sh/helm/v3/pkg/chartutil"
	"io"
	"os"
	"path/filepath"
)

const (
	chartExportDefaultDestination = "."
)

type chartExportClient interface {
	ChartExport(ref string) error
}

type chartExportClientImpl struct {
	cli *chartExport
	env *helmEnv
}

type ChartExportOption struct {
	f func(o *chartExportOptions)
}

type chartExportOptions struct {
	destination string
}

func (o *chartExportOptions) apply(opts []ChartExportOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newChartExportOptions(opts []ChartExportOption) *chartExportOptions {
	options := &chartExportOptions{destination: chartExportDefaultDestination}
	options.apply(opts)
	return options
}

func ChartExportWithDestination(destination string) ChartExportOption {
	return ChartExportOption{f: func(o *chartExportOptions) {
		o.destination = destination
	}}
}

type chartExport struct {
	registryClient *registry.Client
	destination    string
}

func (c *chartExportClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *chartExportClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newChartExport(registryClient)
	// copy args
	copyChartExportClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyChartExportClientOptions(oldCli *chartExport, newCli *chartExport) {
	newCli.destination = oldCli.destination
}

func newChartExportClient(opts []ChartExportOption, env *helmEnv) (*chartExportClientImpl, error) {
	o := newChartExportOptions(opts)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newChartExport(registryClient)
	mergeChartExportOptions(o, client)
	return &chartExportClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newChartExport(registryClient *registry.Client) *chartExport {
	return &chartExport{
		registryClient: registryClient,
	}
}

func (c *chartExport) Run(out io.Writer, ref string) error {
	r, err := registry.ParseReference(ref)
	if err != nil {
		return err
	}

	ch, err := c.registryClient.LoadChart(r)
	if err != nil {
		return err
	}

	// Save the chart to local destination directory
	err = chartutil.SaveDir(ch, c.destination)
	if err != nil {
		return err
	}

	d := filepath.Join(c.destination, ch.Metadata.Name)
	fmt.Fprintf(out, "Exported chart to %s/\n", d)
	return nil
}

func (c *chartExportClientImpl) ChartExport(ref string) error {
	return c.cli.Run(os.Stdout, ref)
}

func mergeChartExportOptions(o *chartExportOptions, cli *chartExport) {
	cli.destination = o.destination
}
