package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRepoRemove(t *testing.T) {
	t.Run("repo remove", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		repoRemoveCli, err := cli.RepoRemove()
		assert.Equal(t, err, nil)
		err = repoRemoveCli.RepoRemove([]string{"REPO"})
	})
}
