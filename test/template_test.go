package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"gotest.tools/assert"
	"testing"
)

func TestTemplate(t *testing.T) {
	t.Run("template", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		templateCli, err := cli.Template([]helmclient.TemplateOption{}, []helmclient.ValueOption{})
		assert.Equal(t, err, nil)
		err = templateCli.Template([]string{"CHART"})
		assert.Equal(t, err, nil)
	})
}
