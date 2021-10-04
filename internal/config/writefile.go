//go:build !windows
// +build !windows

package config

import (
	"os"

	"github.com/google/renameio"
)

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return renameio.WriteFile(filename, data, perm)
}
