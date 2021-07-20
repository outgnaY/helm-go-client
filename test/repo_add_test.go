package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRepoAdd(t *testing.T) {
	t.Run("repo add", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		repoAddCli, err := cli.RepoAdd([]helmclient.RepoAddOption{})
		assert.Equal(t, err, nil)
		err = repoAddCli.RepoAdd([]string{"NAME", "URL"})
		assert.Equal(t, err, nil)
	})
}
