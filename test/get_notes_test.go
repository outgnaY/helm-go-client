package test

import (
	helmclient "github.com/outgnaY/helm-go-client"
	"fmt"
	"gotest.tools/assert"
	"testing"
)

func TestGetNotes(t *testing.T) {
	t.Run("get notes", func(t *testing.T) {
		cli := helmclient.NewHelmClient(kubeConfigForTest, "default")
		getNotes, err := cli.GetNotes([]helmclient.GetNotesOption{})
		assert.Equal(t, err, nil)
		notes, err := getNotes.GetNotes("RELEASE_NAME")
		assert.Equal(t, err, nil)
		fmt.Println(notes)
	})
}
