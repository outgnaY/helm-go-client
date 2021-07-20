package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"os"
	"time"
)

const (
	uninstallDefaultDisableHooks = false
	uninstallDefaultDryRun       = false
	uninstallDefaultKeepHistory  = false
	uninstallDefaultTimeout      = 300 * time.Second
	uninstallDefaultDescription  = ""
)

type uninstallClient interface {
	Uninstall(args []string) error
}

type uninstallClientImpl struct {
	cli *action.Uninstall
	env *helmEnv
}

type UninstallOption struct {
	f func(o *uninstallOptions)
}

type uninstallOptions struct {
	disableHooks bool
	dryRun       bool
	keepHistory  bool
	timeout      time.Duration
	description  string
}

func (o *uninstallOptions) apply(opts []UninstallOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newUninstallOptions(opts []UninstallOption) *uninstallOptions {
	options := &uninstallOptions{
		disableHooks: uninstallDefaultDisableHooks,
		dryRun:       uninstallDefaultDryRun,
		keepHistory:  uninstallDefaultKeepHistory,
		timeout:      uninstallDefaultTimeout,
		description:  uninstallDefaultDescription,
	}
	options.apply(opts)
	return options
}

func UninstallWithDisableHooks(disableHooks bool) UninstallOption {
	return UninstallOption{f: func(o *uninstallOptions) {
		o.disableHooks = disableHooks
	}}
}

func UninstallWithDryRun(dryRun bool) UninstallOption {
	return UninstallOption{f: func(o *uninstallOptions) {
		o.dryRun = dryRun
	}}
}

func UninstallWithKeepHistory(keepHistory bool) UninstallOption {
	return UninstallOption{f: func(o *uninstallOptions) {
		o.keepHistory = keepHistory
	}}
}

func UninstallWithTimeout(timeout time.Duration) UninstallOption {
	return UninstallOption{f: func(o *uninstallOptions) {
		o.timeout = timeout
	}}
}

func UninstallWithDescription(description string) UninstallOption {
	return UninstallOption{f: func(o *uninstallOptions) {
		o.description = description
	}}
}

func (c *uninstallClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *uninstallClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewUninstall(cfg)
	// copy args
	copyUninstallClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyUninstallClientOptions(oldCli *action.Uninstall, newCli *action.Uninstall) {
	newCli.DisableHooks = oldCli.DisableHooks
	newCli.DryRun = oldCli.DryRun
	newCli.KeepHistory = oldCli.KeepHistory
	newCli.Timeout = oldCli.Timeout
	newCli.Description = oldCli.Description
}

func newUninstallClient(opts []UninstallOption, env *helmEnv) (*uninstallClientImpl, error) {
	o := newUninstallOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewUninstall(cfg)
	mergeUninstallOptions(o, client)
	return &uninstallClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *uninstallClientImpl) Uninstall(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("uninstall requires at least 1 argument")
	}
	for i := 0; i < len(args); i++ {
		res, err := c.cli.Run(args[i])
		if err != nil {
			return err
		}
		if res != nil && res.Info != "" {
			fmt.Fprintln(os.Stdout, res.Info)
		}
		fmt.Fprintf(os.Stdout, "release \"%s\" uninstalled\n", args[i])
	}
	return nil
}

func mergeUninstallOptions(o *uninstallOptions, cli *action.Uninstall) {
	cli.DisableHooks = o.disableHooks
	cli.DryRun = o.dryRun
	cli.KeepHistory = o.keepHistory
	cli.Timeout = o.timeout
	cli.Description = o.description
}
