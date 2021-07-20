package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestRollback(t *testing.T) {
	t.Run("rollback release", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		rollbackCli, err := cli.Rollback([]helmclient.RollbackOption{})
		assert.Equal(t, err, nil)
		err = rollbackCli.Rollback([]string{"RELEASE"})
		assert.Equal(t, err, nil)
	})
}
