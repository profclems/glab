package config

import (
	"errors"
	"testing"

	"github.com/profclems/glab/pkg/prompt"
	"github.com/stretchr/testify/assert"
)

func Test_Prompt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ask, teardown := prompt.InitAskStubber()
		defer teardown()

		ask.Stub([]*prompt.QuestionStub{
			{
				Name:  "config",
				Value: "profclems/glab",
			},
		})

		got, err := Prompt("pick a repo", "defaultValue")
		assert.NoError(t, err)
		assert.Equal(t, "profclems/glab", got)
	})
	t.Run("failed", func(t *testing.T) {
		ask, teardown := prompt.InitAskStubber()
		defer teardown()

		ask.Stub([]*prompt.QuestionStub{
			{
				Name:  "config",
				Value: errors.New("failed"),
			},
		})

		got, err := Prompt("pick a repo", "defaultValue")
		assert.EqualError(t, err, "failed")
		assert.Equal(t, got, "")
	})
}
