package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/helmpath"
	"io"
	"os"
	"path/filepath"
)

const (
	createDefaultStarter = ""
)

type createClient interface {
	Create(name string) error
}

type createClientImpl struct {
	createOpts *createOptions
	env        *helmEnv
}

type CreateOption struct {
	f func(o *createOptions)
}

type createOptions struct {
	starter    string // --starter
	name       string
	starterDir string
}

func (o *createOptions) apply(opts []CreateOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newCreateOptions(opts []CreateOption) *createOptions {
	options := &createOptions{
		starter: createDefaultStarter,
	}
	options.apply(opts)
	return options
}

func CreateWithStarter(starter string) CreateOption {
	return CreateOption{f: func(o *createOptions) {
		o.starter = starter
	}}
}

func (c *createClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *createClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newCreateClient(opts []CreateOption, env *helmEnv) (*createClientImpl, error) {
	o := newCreateOptions(opts)
	return &createClientImpl{
		createOpts: o,
		env:        env,
	}, nil
}

func (o *createOptions) run(out io.Writer) error {
	fmt.Fprintf(out, "Creating %s\n", o.name)

	chartname := filepath.Base(o.name)
	cfile := &chart.Metadata{
		Name:        chartname,
		Description: "A Helm chart for Kubernetes",
		Type:        "application",
		Version:     "0.1.0",
		AppVersion:  "0.1.0",
		APIVersion:  chart.APIVersionV2,
	}

	if o.starter != "" {
		// Create from the starter
		lstarter := filepath.Join(o.starterDir, o.starter)
		// If path is absolute, we don't want to prefix it with helm starters folder
		if filepath.IsAbs(o.starter) {
			lstarter = o.starter
		}
		return chartutil.CreateFrom(cfile, filepath.Dir(o.name), lstarter)
	}

	chartutil.Stderr = out
	_, err := chartutil.Create(chartname, filepath.Dir(o.name))
	return err
}

func (c *createClientImpl) Create(name string) error {
	c.createOpts.name = name
	c.createOpts.starterDir = helmpath.DataPath("starters")
	return c.createOpts.run(os.Stdout)
}
