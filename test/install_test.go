package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestInstall(t *testing.T) {
	t.Run("install chart", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		installCli, err := cli.Install([]helmclient.InstallOption{}, []helmclient.ValueOption{}, []helmclient.ChartPathOption{helmclient.WithVersion("1.1.1")})
		assert.Equal(t, err, nil)
		release, err := installCli.Install([]string{"NAME", "CHART"})
		fmt.Println(release)
		assert.Equal(t, err, nil)
	})
}
