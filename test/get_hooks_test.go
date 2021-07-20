package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestGetHooks(t *testing.T) {
	t.Run("get hooks", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		getHooks, err := cli.GetHooks([]helmclient.GetHooksOption{})
		assert.Equal(t, err, nil)
		hooks, err := getHooks.GetHooks("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(hooks)
	})
}
