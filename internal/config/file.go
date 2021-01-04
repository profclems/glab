package config

import (
	"bufio"
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

// CheckFileHasLine : returns true if line exists in file, otherwise false (also for non-existant file).
func CheckFileHasLine(filePath, line string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	fs := bufio.NewScanner(f)
	fs.Split(bufio.ScanLines)

	for fs.Scan() {
		if fs.Text() == line {
			return true
		}
	}

	return false
}

// ReadAndAppend : appends string to file
func ReadAndAppend(file, text string) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(text)); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// BackupConfigFile creates a backup of the provided config file
var BackupConfigFile = func(filename string) error {
	return os.Rename(filename, filename+".bak")
}
