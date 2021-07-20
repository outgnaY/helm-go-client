package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	packageDefaultSign             = false
	packageDefaultKey              = ""
	packageDefaultPassphraseFile   = ""
	packageDefaultVersion          = ""
	packageDefaultAppVersion       = ""
	packageDefaultDestination      = "."
	packageDefaultDependencyUpdate = false
)

var (
	packageDefaultKeyring = defaultKeyring()
)

type packageClient interface {
	Package(args []string) error
}

type packageClientImpl struct {
	cli *action.Package
	env *helmEnv
}

type PackageOption struct {
	f func(o *packageOptions)
}

type packageOptions struct {
	sign             bool
	key              string
	keyring          string
	passphraseFile   string
	version          string
	appVersion       string
	destination      string
	dependencyUpdate bool
}

func (o *packageOptions) apply(opts []PackageOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newPackageOptions(opts []PackageOption) *packageOptions {
	options := &packageOptions{
		sign:             packageDefaultSign,
		key:              packageDefaultKey,
		keyring:          packageDefaultKeyring,
		passphraseFile:   packageDefaultPassphraseFile,
		version:          packageDefaultVersion,
		appVersion:       packageDefaultAppVersion,
		destination:      packageDefaultDestination,
		dependencyUpdate: packageDefaultDependencyUpdate,
	}
	options.apply(opts)
	return options
}

func PackageWithSign(sign bool) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.sign = sign
	}}
}

func PackageWithKey(key string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.key = key
	}}
}

func PackageWithKeyring(keyring string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.keyring = keyring
	}}
}

func PackageWithPassphraseFile(passphraseFile string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.passphraseFile = passphraseFile
	}}
}

func PackageWithVersion(version string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.version = version
	}}
}

func PackageWithAppVersion(appVersion string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.appVersion = appVersion
	}}
}

func PackageWithDestination(destination string) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.destination = destination
	}}
}

func PackageWithDependencyUpdate(dependencyUpdate bool) PackageOption {
	return PackageOption{f: func(o *packageOptions) {
		o.dependencyUpdate = dependencyUpdate
	}}
}

func (c *packageClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *packageClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	client := action.NewPackage()
	// copy args
	copyPackageClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyPackageClientOptions(oldCli *action.Package, newCli *action.Package) {
	newCli.Sign = oldCli.Sign
	newCli.Key = oldCli.Key
	newCli.Keyring = oldCli.Keyring
	newCli.PassphraseFile = oldCli.PassphraseFile
	newCli.Version = oldCli.Version
	newCli.AppVersion = oldCli.AppVersion
	newCli.Destination = oldCli.Destination
	newCli.DependencyUpdate = oldCli.DependencyUpdate
}

func newPackageClient(opts []PackageOption, env *helmEnv) (*packageClientImpl, error) {
	o := newPackageOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewPackage()
	mergePackageOptions(o, client)
	return &packageClientImpl{
		cli: client,
		env: env,
	}, nil
}

func (c *packageClientImpl) Package(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("need at least one argument, the path to the chart")
	}
	if c.cli.Sign {
		if c.cli.Key == "" {
			return fmt.Errorf("--key is required for signing a package")
		}
		if c.cli.Keyring == "" {
			return fmt.Errorf("--keyring is required for signing a package")
		}
	}

	valueOpts := &values.Options{}
	c.cli.RepositoryConfig = c.env.settings.RepositoryConfig
	c.cli.RepositoryCache = c.env.settings.RepositoryCache
	p := getter.All(c.env.settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}

	for i := 0; i < len(args); i++ {
		path, err := filepath.Abs(args[i])
		if err != nil {
			return err
		}
		if _, err := os.Stat(args[i]); err != nil {
			return err
		}

		if c.cli.DependencyUpdate {
			downloadManager := &downloader.Manager{
				Out:              ioutil.Discard,
				ChartPath:        path,
				Keyring:          c.cli.Keyring,
				Getters:          p,
				Debug:            c.env.settings.Debug,
				RepositoryConfig: c.env.settings.RepositoryConfig,
				RepositoryCache:  c.env.settings.RepositoryCache,
			}

			if err := downloadManager.Update(); err != nil {
				return err
			}
		}
		p, err := c.cli.Run(path, vals)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Successfully packaged chart and saved it to: %s\n", p)
	}
	return nil
}

func mergePackageOptions(o *packageOptions, cli *action.Package) {
	cli.Sign = o.sign
	cli.Key = o.key
	cli.Keyring = o.keyring
	cli.PassphraseFile = o.passphraseFile
	cli.Version = o.version
	cli.AppVersion = o.appVersion
	cli.Destination = o.destination
	cli.DependencyUpdate = o.dependencyUpdate
}
