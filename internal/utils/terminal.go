package utils

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// IsTerminal reports whether the file descriptor is connected to a terminal
var IsTerminal = func(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || IsCygwinTerminal(f)
}

func IsCygwinTerminal(f *os.File) bool {
	return isatty.IsCygwinTerminal(f.Fd())
}

var TerminalSize = func(w interface{}) (int, int, error) {
	if f, isFile := w.(*os.File); isFile {
		return terminal.GetSize(int(f.Fd()))
	}

	return 0, 0, fmt.Errorf("%v is not a file", w)
}