package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"os"
	"time"
)

const (
	upgradeDefaultInstall                  = false
	upgradeDefaultDevel                    = false
	upgradeDefaultSkipCRDs                 = false
	upgradeDefaultTimeout                  = 300 * time.Second
	upgradeDefaultWait                     = false
	upgradeDefaultWaitForJobs              = false
	upgradeDefaultDisableHooks             = false
	upgradeDefaultDryRun                   = false
	upgradeDefaultForce                    = false
	upgradeDefaultResetValues              = false
	upgradeDefaultReuseValues              = false
	upgradeDefaultRecreate                 = false
	upgradeDefaultMaxHistory               = 10
	upgradeDefaultAtomic                   = false
	upgradeDefaultCleanupOnFail            = false
	upgradeDefaultSubNotes                 = false
	upgradeDefaultDescription              = ""
	upgradeDefaultDisableOpenAPIValidation = false
	upgradeDefaultCreateNamespace          = false
)

type upgradeClient interface {
	Upgrade(args []string) (*release.Release, error)
}

type upgradeClientImpl struct {
	cli             *action.Upgrade
	env             *helmEnv
	cfg             *action.Configuration
	valueOpts       *valueOptions
	createNamespace bool
}

type UpgradeOption struct {
	f func(o *upgradeOptions)
}

type upgradeOptions struct {
	install                  bool
	devel                    bool
	skipCRDs                 bool
	timeout                  time.Duration
	wait                     bool
	waitForJobs              bool
	disableHooks             bool
	dryRun                   bool
	force                    bool
	resetValues              bool
	reuseValues              bool
	recreate                 bool
	maxHistory               int
	atomic                   bool
	cleanupOnFail            bool
	subNotes                 bool
	description              string
	disableOpenAPIValidation bool
	createNamespace          bool
}

func (o *upgradeOptions) apply(opts []UpgradeOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newUpgradeOptions(opts []UpgradeOption) *upgradeOptions {
	options := &upgradeOptions{
		install:                  upgradeDefaultInstall,
		devel:                    upgradeDefaultDevel,
		skipCRDs:                 upgradeDefaultSkipCRDs,
		timeout:                  upgradeDefaultTimeout,
		wait:                     upgradeDefaultWait,
		waitForJobs:              upgradeDefaultWaitForJobs,
		disableHooks:             upgradeDefaultDisableHooks,
		dryRun:                   upgradeDefaultDryRun,
		force:                    upgradeDefaultForce,
		resetValues:              upgradeDefaultResetValues,
		reuseValues:              upgradeDefaultReuseValues,
		recreate:                 upgradeDefaultRecreate,
		maxHistory:               upgradeDefaultMaxHistory,
		atomic:                   upgradeDefaultAtomic,
		cleanupOnFail:            upgradeDefaultCleanupOnFail,
		subNotes:                 upgradeDefaultSubNotes,
		description:              upgradeDefaultDescription,
		disableOpenAPIValidation: upgradeDefaultDisableOpenAPIValidation,
		createNamespace:          upgradeDefaultCreateNamespace,
	}
	options.apply(opts)
	return options
}

func UpgradeWithInstall(install bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.install = install
	}}
}

func UpgradeWithDevel(devel bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.devel = devel
	}}
}

func UpgradeWithSkipCRDs(skipCRDs bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.skipCRDs = skipCRDs
	}}
}

func UpgradeWithTimeout(timeout time.Duration) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.timeout = timeout
	}}
}

func UpgradeWithWait(wait bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.wait = wait
	}}
}

func UpgradeWithWaitForJobs(waitForJobs bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.waitForJobs = waitForJobs
	}}
}

func UpgradeWithDisableHooks(disableHooks bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.disableHooks = disableHooks
	}}
}

func UpgradeWithDryRun(dryRun bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.dryRun = dryRun
	}}
}

func UpgradeWithForce(force bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.force = force
	}}
}

func UpgradeWithResetValues(resetValues bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.resetValues = resetValues
	}}
}

func UpgradeWithReuseValues(reuseValues bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.reuseValues = reuseValues
	}}
}

func UpgradeWithRecreate(recreate bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.recreate = recreate
	}}
}

func UpgradeWithMaxHistory(maxHistory int) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.maxHistory = maxHistory
	}}
}

func UpgradeWithAtomic(atomic bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.atomic = atomic
	}}
}

func UpgradeWithCleanupOnFail(cleanupOnFail bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.cleanupOnFail = cleanupOnFail
	}}
}

func UpgradeWithSubNotes(subNotes bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.subNotes = subNotes
	}}
}

func UpgradeWithDescription(description string) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.description = description
	}}
}

func UpgradeWithDisableOpenAPIValidation(disableOpenAPIValidation bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.disableOpenAPIValidation = disableOpenAPIValidation
	}}
}

func UpgradeWithCreateNamespace(createNamespace bool) UpgradeOption {
	return UpgradeOption{f: func(o *upgradeOptions) {
		o.createNamespace = createNamespace
	}}
}

func (c *upgradeClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *upgradeClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewUpgrade(cfg)
	// copy args
	copyUpgradeClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	c.cfg = cfg
	return nil
}

func copyUpgradeClientOptions(oldCli *action.Upgrade, newCli *action.Upgrade) {
	newCli.ChartPathOptions = oldCli.ChartPathOptions
	newCli.Install = oldCli.Install
	newCli.Devel = oldCli.Devel
	newCli.SkipCRDs = oldCli.SkipCRDs
	newCli.Timeout = oldCli.Timeout
	newCli.Wait = oldCli.Wait
	newCli.WaitForJobs = oldCli.WaitForJobs
	newCli.DisableHooks = oldCli.DisableHooks
	newCli.DryRun = oldCli.DryRun
	newCli.Force = oldCli.Force
	newCli.ResetValues = oldCli.ResetValues
	newCli.ReuseValues = oldCli.ReuseValues
	newCli.Recreate = oldCli.Recreate
	newCli.MaxHistory = oldCli.MaxHistory
	newCli.Atomic = oldCli.Atomic
	newCli.CleanupOnFail = oldCli.CleanupOnFail
	newCli.SubNotes = oldCli.SubNotes
	newCli.Description = oldCli.Description
	newCli.DisableOpenAPIValidation = oldCli.DisableOpenAPIValidation
}

func newUpgradeClient(opts []UpgradeOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption, env *helmEnv) (*upgradeClientImpl, error) {
	o := newUpgradeOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewUpgrade(cfg)
	mergeUpgradeOptions(o, client)
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
	return &upgradeClientImpl{
		cli:             client,
		env:             env,
		cfg:             cfg,
		valueOpts:       v,
		createNamespace: o.createNamespace,
	}, nil
}

func (c *upgradeClientImpl) Upgrade(args []string) (*release.Release, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("upgrade requires 2 arguments exactly")
	}
	c.cli.Namespace = c.env.Namespace()

	// Fixes #7002 - Support reading values from STDIN for `upgrade` command
	// Must load values AFTER determining if we have to call install so that values loaded from stdin are are not read twice
	if c.cli.Install {
		// If a release does not exist, install it.
		histClient := action.NewHistory(c.cfg)
		histClient.Max = 1
		if _, err := histClient.Run(args[0]); err == driver.ErrReleaseNotFound {
			instClient := action.NewInstall(c.cfg)
			instClient.CreateNamespace = c.createNamespace
			instClient.ChartPathOptions = c.cli.ChartPathOptions
			instClient.DryRun = c.cli.DryRun
			instClient.DisableHooks = c.cli.DisableHooks
			instClient.SkipCRDs = c.cli.SkipCRDs
			instClient.Timeout = c.cli.Timeout
			instClient.Wait = c.cli.Wait
			instClient.WaitForJobs = c.cli.WaitForJobs
			instClient.Devel = c.cli.Devel
			instClient.Namespace = c.cli.Namespace
			instClient.Atomic = c.cli.Atomic
			instClient.PostRenderer = c.cli.PostRenderer
			instClient.DisableOpenAPIValidation = c.cli.DisableOpenAPIValidation
			instClient.SubNotes = c.cli.SubNotes
			instClient.Description = c.cli.Description

			return runInstall(args, instClient, (*values.Options)(c.valueOpts), os.Stdout, c.env)

		} else if err != nil {
			return nil, err
		}
	}

	if c.cli.Version == "" && c.cli.Devel {
		debug("setting version to >0.0.0-0")
		c.cli.Version = ">0.0.0-0"
	}

	chartPath, err := c.cli.ChartPathOptions.LocateChart(args[1], c.env.settings)
	if err != nil {
		return nil, err
	}

	vals, err := (*values.Options)(c.valueOpts).MergeValues(getter.All(c.env.settings))
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}
	if ch.Metadata.Deprecated {
		warning("This chart is deprecated")
	}
	return c.cli.Run(args[0], ch, vals)
}

func mergeUpgradeOptions(o *upgradeOptions, cli *action.Upgrade) {
	cli.Install = o.install
	cli.Devel = o.devel
	cli.SkipCRDs = o.skipCRDs
	cli.Timeout = o.timeout
	cli.Wait = o.wait
	cli.WaitForJobs = o.waitForJobs
	cli.DisableHooks = o.disableHooks
	cli.DryRun = o.dryRun
	cli.Force = o.force
	cli.ResetValues = o.resetValues
	cli.ReuseValues = o.reuseValues
	cli.Recreate = o.recreate
	cli.MaxHistory = o.maxHistory
	cli.Atomic = o.atomic
	cli.CleanupOnFail = o.cleanupOnFail
	cli.SubNotes = o.subNotes
	cli.Description = o.description
	cli.DisableOpenAPIValidation = o.disableOpenAPIValidation
}
