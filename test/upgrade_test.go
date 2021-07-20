package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestUpTrade(t *testing.T) {
	t.Run("upgrade release", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		upgradeCli, err := cli.Upgrade([]helmclient.UpgradeOption{}, []helmclient.ValueOption{}, []helmclient.ChartPathOption{})
		assert.Equal(t, err, nil)
		release, err := upgradeCli.Upgrade([]string{"RELEASE", "CHART"})
		fmt.Println(release)
		assert.Equal(t, err, nil)
	})
}
