package iostreams

import (
	"os"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_isColorEnabled(t *testing.T) {
	preRun := func() {
		os.Unsetenv("NO_COLOR")
		os.Unsetenv("COLOR_ENALBED")
		checkedNoColor = false // Reset it before each run
	}

	t.Run("default", func(t *testing.T) {
		preRun()

		got := isColorEnabled()
		assert.True(t, got)
	})

	t.Run("NO_COLOR", func(t *testing.T) {
		preRun()

		_ = os.Setenv("NO_COLOR", "")

		got := isColorEnabled()
		assert.False(t, got)
	})

	t.Run("COLOR_ENABLED == 1", func(t *testing.T) {
		preRun()

		_ = os.Setenv("NO_COLOR", "")
		_ = os.Setenv("COLOR_ENABLED", "1")

		got := isColorEnabled()
		assert.True(t, got)
	})

	t.Run("COLOR_ENABLED == true", func(t *testing.T) {
		preRun()

		_ = os.Setenv("NO_COLOR", "")
		_ = os.Setenv("COLOR_ENABLED", "true")

		got := isColorEnabled()
		assert.True(t, got)
	})

}
