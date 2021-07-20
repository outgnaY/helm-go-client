package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestHistory(t *testing.T) {
	t.Run("history", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		history, err := cli.History([]helmclient.HistoryOption{})
		assert.Equal(t, err, nil)
		releaseHistory, err := history.History("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(releaseHistory)
	})
}
