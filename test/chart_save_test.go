package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestChartSave(t *testing.T) {
	t.Run("chart save", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		chartSaveCli, err := cli.ChartSave()
		assert.Equal(t, err, nil)
		err = chartSaveCli.ChartSave([]string{"path", "ref"})
		assert.Equal(t, err, nil)
	})
}
