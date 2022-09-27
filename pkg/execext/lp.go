//go:build !windows
// +build !windows

package execext

import "os/exec"

// LookPath searches for an executable named file in the directories named by
// the PATH environment variable.
// If file contains a slash, it is tried directly and the PATH is not consulted.
// The result may be an absolute path or a path relative to the current directory.
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
