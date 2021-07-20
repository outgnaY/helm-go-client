package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestSearchRepo(t *testing.T) {
	t.Run("search repo", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		searchRepoCli, err := cli.SearchRepo([]helmclient.SearchRepoOption{})
		assert.Equal(t, err, nil)
		searchResults, err := searchRepoCli.SearchRepo([]string{})
		assert.Equal(t, err, nil)
		fmt.Println(searchResults)
	})
}
