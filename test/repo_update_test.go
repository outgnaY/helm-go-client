package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRepoUpdate(t *testing.T) {
	t.Run("repo update", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		repoUpdateCli, err := cli.RepoUpdate()
		assert.Equal(t, err, nil)
		err = repoUpdateCli.RepoUpdate()
		assert.Equal(t, err, nil)
	})
}
