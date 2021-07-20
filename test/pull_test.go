package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestPull(t *testing.T) {
	t.Run("pull", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		pullCli, err := cli.Pull([]helmclient.PullOption{}, []helmclient.ChartPathOption{})
		assert.Equal(t, err, nil)
		err = pullCli.Pull([]string{"chart URL | repo/chartname"})
		assert.Equal(t, err, nil)
	})
}
