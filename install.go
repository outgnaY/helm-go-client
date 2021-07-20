package helmclient

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"io"
	"os"
	"time"
)

const (
	installDefaultCreateNamespace          = false
	installDefaultDryRun                   = false
	installDefaultDisableHooks             = false
	installDefaultReplace                  = false
	installDefaultWait                     = false
	installDefaultWaitForJobs              = false
	installDefaultDevel                    = false
	installDefaultDependencyUpdate         = false
	installDefaultTimeout                  = 300 * time.Second
	installDefaultGenerateName             = false
	installDefaultNameTemplate             = ""
	installDefaultDescription              = ""
	installDefaultAtomic                   = false
	installDefaultSkipCRDs                 = false
	installDefaultSubNotes                 = false
	installDefaultDisableOpenAPIValidation = false
)

type installClient interface {
	Install(args []string) (*release.Release, error)
}

type installClientImpl struct {
	cli       *action.Install
	env       *helmEnv
	valueOpts *valueOptions
}

type InstallOption struct {
	f func(o *installOptions)
}

type installOptions struct {
	createNamespace          bool
	dryRun                   bool
	disableHooks             bool
	replace                  bool
	wait                     bool
	waitForJobs              bool
	devel                    bool
	dependencyUpdate         bool
	timeout                  time.Duration
	generateName             bool
	nameTemplate             string
	description              string
	atomic                   bool
	skipCRDs                 bool
	subNotes                 bool
	disableOpenAPIValidation bool
}

func (o *installOptions) apply(opts []InstallOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newInstallOptions(opts []InstallOption) *installOptions {
	options := &installOptions{
		createNamespace:          installDefaultCreateNamespace,
		dryRun:                   installDefaultDryRun,
		disableHooks:             installDefaultDisableHooks,
		replace:                  installDefaultReplace,
		wait:                     installDefaultWait,
		waitForJobs:              installDefaultWaitForJobs,
		devel:                    installDefaultDevel,
		dependencyUpdate:         installDefaultDependencyUpdate,
		timeout:                  installDefaultTimeout,
		generateName:             installDefaultGenerateName,
		nameTemplate:             installDefaultNameTemplate,
		description:              installDefaultDescription,
		atomic:                   installDefaultAtomic,
		skipCRDs:                 installDefaultSkipCRDs,
		subNotes:                 installDefaultSubNotes,
		disableOpenAPIValidation: installDefaultDisableOpenAPIValidation,
	}
	options.apply(opts)
	return options
}

func InstallWithCreateNamespace(createNamespace bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.createNamespace = createNamespace
	}}
}

func InstallWithDryRun(dryRun bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.dryRun = dryRun
	}}
}

func InstallWithDisableHooks(disableHooks bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.disableHooks = disableHooks
	}}
}

func InstallWithReplace(replace bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.replace = replace
	}}
}

func InstallWithWait(wait bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.wait = wait
	}}
}

func InstallWithWaitForJobs(waitForJobs bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.waitForJobs = waitForJobs
	}}
}

func InstallWithDevel(devel bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.devel = devel
	}}
}

func InstallWithDependencyUpdate(dependencyUpdate bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.dependencyUpdate = dependencyUpdate
	}}
}

func InstallWithTimeout(timeout time.Duration) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.timeout = timeout
	}}
}

func InstallWithGenerateName(generateName bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.generateName = generateName
	}}
}

func InstallWithNameTemplate(nameTemplate string) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.nameTemplate = nameTemplate
	}}
}

func InstallWithDescription(description string) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.description = description
	}}
}

func InstallWithAtomic(atomic bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.atomic = atomic
	}}
}

func InstallWithSkipCRDs(skipCRDs bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.skipCRDs = skipCRDs
	}}
}

func InstallWithSubNotes(subNotes bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.subNotes = subNotes
	}}
}

func InstallWithDisableOpenAPIValidation(disableOpenAPIValidation bool) InstallOption {
	return InstallOption{f: func(o *installOptions) {
		o.disableOpenAPIValidation = disableOpenAPIValidation
	}}
}

func (c *installClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *installClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewInstall(cfg)
	// copy args
	copyInstallClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyInstallClientOptions(oldCli *action.Install, newCli *action.Install) {
	newCli.ChartPathOptions = oldCli.ChartPathOptions
	newCli.CreateNamespace = oldCli.CreateNamespace
	newCli.DryRun = oldCli.DryRun
	newCli.DisableHooks = oldCli.DisableHooks
	newCli.Replace = oldCli.Replace
	newCli.Wait = oldCli.Wait
	newCli.WaitForJobs = oldCli.WaitForJobs
	newCli.Devel = oldCli.Devel
	newCli.DependencyUpdate = oldCli.DependencyUpdate
	newCli.Timeout = oldCli.Timeout
	newCli.GenerateName = oldCli.GenerateName
	newCli.NameTemplate = oldCli.NameTemplate
	newCli.Description = oldCli.Description
	newCli.Atomic = oldCli.Atomic
	newCli.SkipCRDs = oldCli.SkipCRDs
	newCli.SubNotes = oldCli.SubNotes
	newCli.DisableOpenAPIValidation = oldCli.DisableOpenAPIValidation
}

func newInstallClient(opts []InstallOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption, env *helmEnv) (*installClientImpl, error) {
	o := newInstallOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(cfg)
	mergeInstallOptions(o, client)
	v := &valueOptions{
		ValueFiles:   []string{},
		StringValues: []string{},
		Values:       []string{},
		FileValues:   []string{},
	}
	addValueOptions(valueOpts, v)
	c := &chartPathOptions{
		CaFile:                "",
		CertFile:              "",
		KeyFile:               "",
		InsecureSkipTLSverify: false,
		Keyring:               defaultKeyring(),
		Password:              "",
		PassCredentialsAll:    false,
		RepoURL:               "",
		Username:              "",
		Verify:                false,
		Version:               "",
	}
	addChartPathOptions(chartPathOpts, c)
	mergeChartPathOptions(c, &client.ChartPathOptions)
	return &installClientImpl{
		cli:       client,
		env:       env,
		valueOpts: v,
	}, nil
}

func (c *installClientImpl) Install(args []string) (*release.Release, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("install requires at least 1 argument")
	}
	return runInstall(args, c.cli, (*values.Options)(c.valueOpts), os.Stdout, c.env)
}

func runInstall(args []string, client *action.Install, valueOpts *values.Options, out io.Writer, env *helmEnv) (*release.Release, error) {
	debug("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}
	name, chart, err := client.NameAndChart(args)
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name

	cp, err := client.ChartPathOptions.LocateChart(chart, env.settings)
	if err != nil {
		return nil, err
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(env.settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}
	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		warning("This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              out,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: env.settings.RepositoryConfig,
					RepositoryCache:  env.settings.RepositoryCache,
					Debug:            env.settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = env.settings.Namespace()
	return client.Run(chartRequested, vals)
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func mergeInstallOptions(o *installOptions, cli *action.Install) {
	cli.CreateNamespace = o.createNamespace
	cli.DryRun = o.dryRun
	cli.DisableHooks = o.disableHooks
	cli.Replace = o.replace
	cli.Wait = o.wait
	cli.WaitForJobs = o.waitForJobs
	cli.Devel = o.devel
	cli.DependencyUpdate = o.dependencyUpdate
	cli.Timeout = o.timeout
	cli.GenerateName = o.generateName
	cli.NameTemplate = o.nameTemplate
	cli.Description = o.description
	cli.Atomic = o.atomic
	cli.SkipCRDs = o.skipCRDs
	cli.SubNotes = o.subNotes
	cli.DisableOpenAPIValidation = o.disableOpenAPIValidation
}
