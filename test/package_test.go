package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestPackage(t *testing.T) {
	t.Run("package", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		packageCli, err := cli.Package([]helmclient.PackageOption{})
		assert.Equal(t, err, nil)
		err = packageCli.Package([]string{"CHART_PATH"})
		assert.Equal(t, err, nil)
	})
}
