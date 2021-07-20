package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestGetValues(t *testing.T) {
	t.Run("get values", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		getValues, err := cli.GetValues([]helmclient.GetValuesOption{})
		assert.Equal(t, err, nil)
		m, err := getValues.GetValues("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(m)
	})
}
