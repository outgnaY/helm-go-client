package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRepoIndex(t *testing.T) {
	t.Run("repo index", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		repoIndexCli, err := cli.RepoIndex([]helmclient.RepoIndexOption{})
		assert.Equal(t, err, nil)
		err = repoIndexCli.RepoIndex("DIR")
		assert.Equal(t, err, nil)
	})
}
