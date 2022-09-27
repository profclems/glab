//go:build !windows
// +build !windows

package config

import (
	"os"
	"path/filepath"

	"github.com/google/renameio"
)

// WriteFile to the path
// If the path is smylink it will write to the symlink
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	pathToSymlink, err := filepath.EvalSymlinks(filename)
	if err == nil {
		filename = pathToSymlink
	}

	return renameio.WriteFile(filename, data, perm)
}
