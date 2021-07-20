package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestChartPull(t *testing.T) {
	t.Run("chart pull", func(t *testing.T) {
		cli := helmclient.NewHelmClient("", "")
		chartPullCli, err := cli.ChartPull()
		assert.Equal(t, err, nil)
		err = chartPullCli.ChartPull("ref")
		assert.Equal(t, err, nil)
	})
}
