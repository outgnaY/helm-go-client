package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestRepoList(t *testing.T) {
	t.Run("repo list", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		repoListCli, err := cli.RepoList()
		assert.Equal(t, err, nil)
		repos, err := repoListCli.RepoList()
		assert.Equal(t, err, nil)
		fmt.Println(repos)
	})
}
