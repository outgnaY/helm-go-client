package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestChartPush(t *testing.T) {
	t.Run("chart push", func(t *testing.T) {
		cli := helmclient.NewHelmClient("", "")
		chartPushCli, err := cli.ChartPush()
		assert.Equal(t, err, nil)
		err = chartPushCli.ChartPush("ref")
		assert.Equal(t, err, nil)
	})
}
