package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocalConfigDir(t *testing.T) {
	got := LocalConfigDir()
	assert.ElementsMatch(t, []string{".git", "glab-cli"}, got)
}

func Test_LocalConfigFile(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		expectedPath := filepath.Join(".git", "glab-cli", "config.yml")
		got := LocalConfigFile()
		assert.Equal(t, expectedPath, got)
	})

	t.Run("old config file", func(t *testing.T) {
		expectedPath := filepath.Join(".glab-cli", "config", "config.yml")
		got := OldLocalConfigFile()
		assert.Equal(t, expectedPath, got)
	})

	t.Run("modified-LocalConfigDir()", func(t *testing.T) {
		expectedPath := filepath.Join(".config", "glab-cli", "config.yml")

		LocalConfigDir = func() []string {
			return []string{".config", "glab-cli"}
		}

		got := LocalConfigFile()
		assert.Equal(t, expectedPath, got)
	})
}
