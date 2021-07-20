package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestGetManifest(t *testing.T) {
	t.Run("get manifest", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		getManifest, err := cli.GetManifest([]helmclient.GetManifestOption{})
		assert.Equal(t, err, nil)
		manifest, err := getManifest.GetManifest("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(manifest)
	})
}
