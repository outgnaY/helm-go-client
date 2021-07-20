package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Run("show helm version", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		versionCli, err := cli.Version([]helmclient.VersionOption{ /*helmclient.VersionWithShort(true)*/ })
		assert.Equal(t, err, nil)
		version, err := versionCli.Version()
		assert.Equal(t, err, nil)
		fmt.Println(version)
	})
}
