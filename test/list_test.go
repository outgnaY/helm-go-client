package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestList(t *testing.T) {
	t.Run("list releases", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		listCli, err := cli.List([]helmclient.ListOption{})
		assert.Equal(t, err, nil)
		releases, err := listCli.List()
		assert.Equal(t, err, nil)
		fmt.Println(releases)
	})
}
