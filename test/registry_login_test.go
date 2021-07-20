package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRegistryLogin(t *testing.T) {
	t.Run("registry login", func(t *testing.T) {
		cli := helmclient.NewHelmClient("", "")
		registryLoginCli, err := cli.RegistryLogin([]helmclient.RegistryLoginOption{})
		assert.Equal(t, err, nil)
		err = registryLoginCli.RegistryLogin("host", "username", "password")
		assert.Equal(t, err, nil)
	})
}
