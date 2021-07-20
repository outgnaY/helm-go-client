package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/experimental/registry"
	"os"
)

type registryLogoutClient interface {
	RegistryLogout(hostname string) error
}

type registryLogoutClientImpl struct {
	cli *registryLogout
	env *helmEnv
}

type registryLogout struct {
	registryClient *registry.Client
}

func (c *registryLogoutClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *registryLogoutClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return err
	}
	client := newRegistryLogout(registryClient)
	c.cli = client
	c.env = env
	return nil
}

func newRegistryLogoutClient(env *helmEnv) (*registryLogoutClientImpl, error) {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(env.settings.Debug),
		registry.ClientOptWriter(os.Stdout),
		registry.ClientOptCredentialsFile(env.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	client := newRegistryLogout(registryClient)
	return &registryLogoutClientImpl{
		cli: client,
		env: env,
	}, nil
}

func newRegistryLogout(registryClient *registry.Client) *registryLogout {
	return &registryLogout{
		registryClient: registryClient,
	}
}

func (c *registryLogout) Run(hostname string) error {
	return c.registryClient.Logout(hostname)
}

func (c *registryLogoutClientImpl) RegistryLogout(hostname string) error {
	return c.cli.Run(hostname)
}
