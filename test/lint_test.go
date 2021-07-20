package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestLint(t *testing.T) {
	t.Run("lint", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		lintCli, err := cli.Lint([]helmclient.LintOption{}, []helmclient.ValueOption{})
		assert.Equal(t, err, nil)
		err = lintCli.Lint([]string{"PATH"})
		fmt.Println(err)
	})
}
