package utils

import (
	"testing"

	"github.com/alecthomas/assert"
)

func Test_HelperFuncs(t *testing.T) {
	t.Run("StringToInt()", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			got := StringToInt("200")
			assert.Equal(t, 200, got)
		})
		t.Run("failed-return-0", func(t *testing.T) {
			got := StringToInt("NotAnInt")
			assert.Equal(t, 0, got)
		})
	})
	t.Run("ReplaceNonAlphaNumericChars()", func(t *testing.T) {
		got := ReplaceNonAlphaNumericChars("profclems-glab", "/")
		assert.Equal(t, "profclems/glab", got)
	})
}
