package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

const (
	registryLoginDefaultInsecure = false
)

type registryLoginClient interface {
	RegistryLogin(hostname string, username string, password string) error
}

type registryLoginClientImpl struct {
	cli      *registryLogin
	env      *helmEnv
	insecure bool
}

type registryLogin struct {
	registryClient *registry.Client
}

type RegistryLoginOption struct {
	f func(o *registryLoginOptions)
}

type registryLoginOptions struct {
	insecure bool
}

func (o *registryLoginOptions) apply(opts []RegistryLoginOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newRegistryLoginOptions(opts []RegistryLoginOption) *registryLoginOptions {
	options := &registryLoginOptions{insecure: registryLoginDefaultInsecure}
	options.apply(opts)
	return options
}

func RegistryLoginWithInsecure(insecure bool) RegistryLoginOption {
	return RegistryLoginOption{f: func(o *registryLoginOptions) {
		o.insecure = insecure
	}}
}

func (c *registryLoginClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *registryLoginClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newRegistryLogin(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newRegistryLoginClient(opts []RegistryLoginOption, env *helmEnv) (*registryLoginClientImpl, error) {
	o := newRegistryLoginOptions(opts)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newRegistryLogin(registryClient)
	return &registryLoginClientImpl{
		cli:      client,
		env:      env,
		insecure: o.insecure,
	}, nil
}

func newRegistryLogin(registryClient *registry.Client) *registryLogin {
	return &registryLogin{
		registryClient: registryClient,
	}
}

func (c *registryLogin) Run(hostname string, username string, password string, insecure bool) error {
	return c.registryClient.Login(hostname, username, password, insecure)
}

func (c *registryLoginClientImpl) RegistryLogin(hostname string, username string, password string) error {
	return c.cli.Run(hostname, username, password, c.insecure)
}
