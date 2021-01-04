package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_HelperFunctions(t *testing.T) {
	// Base ios object that is modifiede as required
	ios := &IOStreams{
		In:     os.Stdin,
		StdOut: NewColorable(os.Stdout),
		StdErr: NewColorable(os.Stderr),

		IsaTTY:         IsTerminal(os.Stdout),
		IsErrTTY:       IsTerminal(os.Stderr),
		IsInTTY:        IsTerminal(os.Stdin),
		promptDisabled: false,

		pagerCommand: os.Getenv("PAGER"),
	}

	t.Run("InitIOStream()", func(t *testing.T) {
		got := InitIOStream()

		assert.Equal(t, ios.In, got.In)
		assert.Equal(t, ios.IsaTTY, got.IsaTTY)
		assert.Equal(t, ios.IsErrTTY, got.IsErrTTY)
		assert.Equal(t, ios.IsInTTY, got.IsInTTY)
		assert.Equal(t, ios.promptDisabled, got.promptDisabled)
		assert.Equal(t, ios.pagerCommand, got.pagerCommand)
	})

	t.Run("IsOutputTTY()", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			ios := *ios

			ios.IsaTTY = true
			ios.IsErrTTY = true

			got := ios.IsOutputTTY()
			assert.True(t, got)
		})
		t.Run("false", func(t *testing.T) {
			t.Run("IsaTTY=false", func(t *testing.T) {
				ios := *ios

				ios.IsaTTY = false
				ios.IsErrTTY = true

				got := ios.IsOutputTTY()
				assert.False(t, got)
			})
			t.Run("IsErrTTY=false", func(t *testing.T) {
				ios := *ios

				ios.IsaTTY = true
				ios.IsErrTTY = false

				got := ios.IsOutputTTY()
				assert.False(t, got)
			})
		})

	})

	t.Run("SetPager()", func(t *testing.T) {
		t.Run("more", func(t *testing.T) {
			ios := *ios
			ios.SetPager("more")
			assert.Equal(t, "more", ios.pagerCommand)
		})
	})

	t.Run("PromptEnabled()", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			var got bool
			ios := *ios

			ios.promptDisabled = false
			ios.IsaTTY = true
			ios.IsErrTTY = true

			got = ios.PromptEnabled()
			assert.True(t, got)
		})

		t.Run("false", func(t *testing.T) {
			t.Run("promptDisabled=true", func(t *testing.T) {
				var got bool
				ios := *ios

				ios.promptDisabled = true
				got = ios.PromptEnabled()
				assert.False(t, got)

			})

			t.Run("IsaTTY=false", func(t *testing.T) {
				var got bool
				ios := *ios

				ios.IsaTTY = false
				got = ios.PromptEnabled()
				assert.False(t, got)
			})

			t.Run("IsErrTTY=true", func(t *testing.T) {
				var got bool
				ios := *ios

				ios.IsErrTTY = false
				got = ios.PromptEnabled()
				assert.False(t, got)
			})
		})
	})

	t.Run("ColorEnabled()", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			ios := *ios

			ios.IsaTTY = true
			ios.IsErrTTY = true
			got := ios.ColorEnabled()
			assert.True(t, got)
		})
		t.Run("false", func(t *testing.T) {
			t.Run("IsaTTY=false", func(t *testing.T) {
				ios := *ios

				ios.IsaTTY = false
				ios.IsErrTTY = true
				got := ios.ColorEnabled()
				assert.False(t, got)
			})
			t.Run("IsErrTTY=false", func(t *testing.T) {
				ios := *ios

				ios.IsaTTY = true
				ios.IsErrTTY = false
				got := ios.ColorEnabled()
				assert.False(t, got)
			})
		})
	})

	t.Run("SetPrompt()", func(t *testing.T) {
		t.Run("disabled", func(t *testing.T) {
			t.Run("true", func(t *testing.T) {
				ios := *ios
				ios.SetPrompt("true")
				assert.True(t, ios.promptDisabled)
			})
			t.Run("1", func(t *testing.T) {
				ios := *ios
				ios.SetPrompt("1")
				assert.True(t, ios.promptDisabled)
			})
		})
		t.Run("enabled", func(t *testing.T) {
			t.Run("false", func(t *testing.T) {
				ios := *ios
				ios.SetPrompt("false")
				assert.False(t, ios.promptDisabled)
			})
			t.Run("0", func(t *testing.T) {
				ios := *ios
				ios.SetPrompt("0")
				assert.False(t, ios.promptDisabled)
			})
		})
	})

	t.Run("IOTest()", func(t *testing.T) {
		io, in, out, err := IOTest()

		assert.Equal(t, io.In, ioutil.NopCloser(in))
		assert.Equal(t, io.StdOut, out)
		assert.Equal(t, io.StdErr, err)

		assert.Equal(t, in, &bytes.Buffer{})
		assert.Equal(t, out, &bytes.Buffer{})
		assert.Equal(t, err, &bytes.Buffer{})
	})
}
