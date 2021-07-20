package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"os"
)

const (
	pullDefaultDevel       = false
	pullDefaultUntar       = false
	pullDefaultVerifyLater = false
	pullDefaultUntarDir    = "."
	pullDefaultDestDir     = "."
)

type pullClient interface {
	Pull(args []string) error
}

type pullClientImpl struct {
	cli *action.Pull
	env *helmEnv
}

type PullOption struct {
	f func(o *pullOptions)
}

type pullOptions struct {
	devel       bool
	untar       bool
	verifyLater bool
	untarDir    string
	destDir     string
}

func (o *pullOptions) apply(opts []PullOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newPullOptions(opts []PullOption) *pullOptions {
	options := &pullOptions{
		devel:       pullDefaultDevel,
		untar:       pullDefaultUntar,
		verifyLater: pullDefaultVerifyLater,
		untarDir:    pullDefaultUntarDir,
		destDir:     pullDefaultDestDir,
	}
	options.apply(opts)
	return options
}

func PullWithDevel(devel bool) PullOption {
	return PullOption{f: func(o *pullOptions) {
		o.devel = devel
	}}
}

func PullWithUntar(untar bool) PullOption {
	return PullOption{f: func(o *pullOptions) {
		o.untar = untar
	}}
}

func PullWithVerifyLater(verifyLater bool) PullOption {
	return PullOption{f: func(o *pullOptions) {
		o.verifyLater = verifyLater
	}}
}

func PullWithUntarDir(untarDir string) PullOption {
	return PullOption{f: func(o *pullOptions) {
		o.untarDir = untarDir
	}}
}

func PullWithDestDir(destDir string) PullOption {
	return PullOption{f: func(o *pullOptions) {
		o.destDir = destDir
	}}
}

func (c *pullClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *pullClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewPullWithOpts(action.WithConfig(cfg))
	// copy args
	copyPullClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyPullClientOptions(oldCli *action.Pull, newCli *action.Pull) {
	newCli.ChartPathOptions = oldCli.ChartPathOptions
	newCli.Devel = oldCli.Devel
	newCli.Untar = oldCli.Untar
	newCli.VerifyLater = oldCli.VerifyLater
	newCli.UntarDir = oldCli.UntarDir
	newCli.DestDir = oldCli.DestDir
}

func newPullClient(opts []PullOption, chartPathOpts []ChartPathOption, env *helmEnv) (*pullClientImpl, error) {
	o := newPullOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewPullWithOpts(action.WithConfig(cfg))
	mergePullOptions(o, client)
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
	return &pullClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *pullClientImpl) Pull(args []string) error {
	c.cli.Settings = c.env.settings
	if c.cli.Version == "" && c.cli.Devel {
		debug("setting version to >0.0.0-0")
		c.cli.Version = ">0.0.0-0"
	}

	for i := 0; i < len(args); i++ {
		output, err := c.cli.Run(args[i])
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, output)
	}
	return nil
}

func mergePullOptions(o *pullOptions, cli *action.Pull) {
	cli.Devel = o.devel
	cli.Untar = o.untar
	cli.VerifyLater = o.verifyLater
	cli.UntarDir = o.untarDir
	cli.DestDir = o.destDir
}
