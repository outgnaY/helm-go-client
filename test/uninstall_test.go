package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestUninstall(t *testing.T) {
	t.Run("uninstall release", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		uninstallCli, err := cli.Uninstall([]helmclient.UninstallOption{})
		assert.Equal(t, err, nil)
		// hello-app preinstalled
		err = uninstallCli.Uninstall([]string{"RELEASE_NAME"})
		assert.Equal(t, err, nil)
	})
}
