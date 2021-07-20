package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestChartRemove(t *testing.T) {
	t.Run("chart remove", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		chartRemoveCli, err := cli.ChartRemove()
		assert.Equal(t, err, nil)
		err = chartRemoveCli.ChartRemove("ref")
		assert.Equal(t, err, nil)
	})
}
