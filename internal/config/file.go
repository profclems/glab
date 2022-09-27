package config

import (
	"os"
)

// CheckPathExists checks if a folder exists and is a directory
func CheckPathExists(path string) bool {
	if info, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		return info.IsDir()
	}
	return false
}

// CheckFileExists : checks if a file exists and is not a directory.
func CheckFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// BackupConfigFile creates a backup of the provided config file
var BackupConfigFile = func(filename string) error {
	return os.Rename(filename, filename+".bak")
}
