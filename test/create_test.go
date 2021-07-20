package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Run("create chart directory", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		createCli, err := cli.Create([]helmclient.CreateOption{})
		assert.Equal(t, err, nil)
		err = createCli.Create("name")
		assert.Equal(t, err, nil)
	})
}
