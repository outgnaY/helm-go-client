package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRegistryLogout(t *testing.T) {
	t.Run("registry logout", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		registryLoginCli, err := cli.RegistryLogin([]helmclient.RegistryLoginOption{})
		assert.Equal(t, err, nil)
		err = registryLoginCli.RegistryLogin("host", "username", "password")
		assert.Equal(t, err, nil)

		registryLogoutCli, err := cli.RegistryLogout()
		assert.Equal(t, err, nil)
		err = registryLogoutCli.RegistryLogout("host")
		assert.Equal(t, err, nil)
	})
}
