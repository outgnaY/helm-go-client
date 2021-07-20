package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestChartList(t *testing.T) {
	t.Run("chart list", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		chartListCli, err := cli.ChartList()
		assert.Equal(t, err, nil)
		chartInfos, err := chartListCli.ChartList()
		assert.Equal(t, err, nil)
		fmt.Println(chartInfos[0])
	})
}
