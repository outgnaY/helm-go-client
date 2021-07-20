package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestChartExport(t *testing.T) {
	t.Run("chart export", func(t *testing.T) {
		cli := helmclient.NewHelmClient("", "")
		chartExportCli, err := cli.ChartExport([]helmclient.ChartExportOption{})
		assert.Equal(t, err, nil)
		err = chartExportCli.ChartExport("ref")
	})
}
