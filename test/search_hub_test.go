package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestSearchHub(t *testing.T) {
	t.Run("search hub", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		searchHubCli, err := cli.SearchHub([]helmclient.SearchHubOption{})
		assert.Equal(t, err, nil)
		searchResults, err := searchHubCli.SearchHub([]string{})
		assert.Equal(t, err, nil)
		fmt.Println(searchResults)
	})
}
