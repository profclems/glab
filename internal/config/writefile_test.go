package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WriteFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Skipf("unexpected error while creating temporary directory = %s", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	fpath := filepath.Join(dir, "test-file")

	t.Run("regular", func(t *testing.T) {
		require.Nilf(t,
			WriteFile(fpath, []byte("profclems/glab"), 0644),
			"unexpected error = %s", err,
		)

		result, err := ioutil.ReadFile(fpath)
		require.Nilf(t, err, "failed to read file %q due to %q", fpath, err)
		assert.Equal(t, "profclems/glab", string(result))

		permissions, err := os.Stat(fpath)
		require.Nilf(t, err, "failed to get stats for file %q due to %q", fpath, err)
		// TODO:
		assert.Equal(t, "0644", fmt.Sprintf("%04o", permissions.Mode()))
	})

	t.Run("symlink", func(t *testing.T) {
		symPath := filepath.Join(dir, "test-symlink")
		require.Nil(t, os.Symlink(fpath, symPath), "failed to create a symlink")
		require.Nilf(t,
			WriteFile(symPath, []byte("profclems/glab/symlink"), 0644),
			"unexpected error = %s", err,
		)

		result, err := ioutil.ReadFile(symPath)
		require.Nilf(t, err, "failed to read file %q due to %q", symPath, err)
		assert.Equal(t, "profclems/glab/symlink", string(result))

		permissions, err := os.Lstat(symPath)
		require.Nil(t, err, "failed to get info about the smylink", err)
		assert.Equal(t, os.ModeSymlink, permissions.Mode()&os.ModeSymlink, "this file should be a symlink")
	})
}
