package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"strconv"
	"time"
)

const (
	rollbackDefaultTimeout       = 300 * time.Second
	rollbackDefaultWait          = false
	rollbackDefaultWaitForJobs   = false
	rollbackDefaultDisableHooks  = false
	rollbackDefaultDryRun        = false
	rollbackDefaultRecreate      = false
	rollbackDefaultForce         = false
	rollbackDefaultCleanupOnFail = false
	rollbackDefaultMaxHistory    = 10
)

type rollbackClient interface {
	Rollback(args []string) error
}

type rollbackClientImpl struct {
	cli *action.Rollback
	env *helmEnv
}

type RollbackOption struct {
	f func(o *rollbackOptions)
}

type rollbackOptions struct {
	// version int
	timeout       time.Duration
	wait          bool
	waitForJobs   bool
	disableHooks  bool
	dryRun        bool
	recreate      bool
	force         bool
	cleanupOnFail bool
	maxHistory    int
}

func (o *rollbackOptions) apply(opts []RollbackOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newRollbackOptions(opts []RollbackOption) *rollbackOptions {
	options := &rollbackOptions{
		timeout:       rollbackDefaultTimeout,
		wait:          rollbackDefaultWait,
		waitForJobs:   rollbackDefaultWaitForJobs,
		disableHooks:  rollbackDefaultDisableHooks,
		dryRun:        rollbackDefaultDryRun,
		recreate:      rollbackDefaultRecreate,
		force:         rollbackDefaultForce,
		cleanupOnFail: rollbackDefaultCleanupOnFail,
		maxHistory:    rollbackDefaultMaxHistory,
	}
	options.apply(opts)
	return options
}

func RollbackWithTimeout(timeout time.Duration) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.timeout = timeout
	}}
}

func RollbackWithWait(wait bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.wait = wait
	}}
}

func RollbackWithWaitForJobs(waitForJobs bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.waitForJobs = waitForJobs
	}}
}

func RollbackWithDisableHooks(disableHooks bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.disableHooks = disableHooks
	}}
}

func RollbackWithDryRun(dryRun bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.dryRun = dryRun
	}}
}

func RollbackWithRecreate(recreate bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.recreate = recreate
	}}
}

func RollbackWithForce(force bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.force = force
	}}
}

func RollbackWithCleanupOnFail(cleanUpOnFail bool) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.cleanupOnFail = cleanUpOnFail
	}}
}

func RollbackWithMaxHistory(maxHistory int) RollbackOption {
	return RollbackOption{f: func(o *rollbackOptions) {
		o.maxHistory = maxHistory
	}}
}

func (c *rollbackClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *rollbackClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewRollback(cfg)
	// copy args
	copyRollbackClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyRollbackClientOptions(oldCli *action.Rollback, newCli *action.Rollback) {
	newCli.Timeout = oldCli.Timeout
	newCli.Wait = oldCli.Wait
	newCli.WaitForJobs = oldCli.WaitForJobs
	newCli.DisableHooks = oldCli.DisableHooks
	newCli.DryRun = oldCli.DryRun
	newCli.Recreate = oldCli.Recreate
	newCli.Force = oldCli.Force
	newCli.CleanupOnFail = oldCli.CleanupOnFail
	newCli.MaxHistory = oldCli.MaxHistory
}

func newRollbackClient(opts []RollbackOption, env *helmEnv) (*rollbackClientImpl, error) {
	o := newRollbackOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewRollback(cfg)
	mergeRollbackOptions(o, client)
	return &rollbackClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *rollbackClientImpl) Rollback(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("rollback requires at least 1 argument")
	}
	if len(args) > 1 {
		ver, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("could not convert revision to a number: %v", err)
		}
		c.cli.Version = ver
	}
	return c.cli.Run(args[0])
}

func mergeRollbackOptions(o *rollbackOptions, cli *action.Rollback) {
	// cli.Version = o.version
	cli.Timeout = o.timeout
	cli.Wait = o.wait
	cli.WaitForJobs = o.waitForJobs
	cli.DisableHooks = o.disableHooks
	cli.DryRun = o.dryRun
	cli.Recreate = o.recreate
	cli.Force = o.force
	cli.CleanupOnFail = o.cleanupOnFail
	cli.MaxHistory = o.maxHistory
}
