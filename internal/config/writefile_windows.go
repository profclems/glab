package config

import (
	"io/ioutil"
	"os"
)

// Note: this is not atomic, but apparently there's no way to atomically
//       replace a file on windows which is why renameio doesn't support
//       windows.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
