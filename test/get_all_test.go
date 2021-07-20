package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestGetAll(t *testing.T) {
	t.Run("get all", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		getAll, err := cli.GetAll([]helmclient.GetAllOption{})
		assert.Equal(t, err, nil)
		release, err := getAll.GetAll("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(release)
	})
}
