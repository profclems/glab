package utils

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

var (
	_isColorEnabled = true
	_isStdoutTerminal = false
	checkedTerminal = false
	checkedNoColor = false

	// Magenta outputs ANSI color if stdout is a tty
	Magenta = makeColorFunc("magenta")

	// Cyan outputs ANSI color if stdout is a tty
	Cyan = makeColorFunc("cyan")

	// Red outputs ANSI color if stdout is a tty
	Red = makeColorFunc("red")

	// Yellow outputs ANSI color if stdout is a tty
	Yellow = makeColorFunc("yellow")

	// Blue outputs ANSI color if stdout is a tty
	Blue = makeColorFunc("blue")

	// Green outputs ANSI color if stdout is a tty
	Green = makeColorFunc("green")

	// Gray outputs ANSI color if stdout is a tty
	Gray = makeColorFunc("black+h")

	// Bold outputs ANSI color if stdout is a tty
	Bold = makeColorFunc("default+b")
)

func isStdoutTerminal() bool {
	if !checkedTerminal {
		_isStdoutTerminal = IsTerminal(os.Stdout)
		checkedTerminal = true
	}
	return _isStdoutTerminal
}

// IsTerminal reports whether the file descriptor is connected to a terminal
func IsTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

// NewColorable returns an output stream that handles ANSI color sequences on Windows
func NewColorable(f *os.File) io.Writer {
	return colorable.NewColorable(f)
}

func makeColorFunc(color string) func(string) string {
	cf := ansi.ColorFunc(color)
	return func(arg string) string {
		if isColorEnabled() && isStdoutTerminal() {
			return cf(arg)
		}
		return arg
	}
}

func isColorEnabled() bool {
	if !checkedNoColor {
		_isColorEnabled = os.Getenv("NO_COLOR") == "" ||
			os.Getenv("COLOR_ENABLED") == "1" ||
			os.Getenv("COLOR_ENABLED") == "true"
		checkedNoColor = true
	}
	return _isColorEnabled
}
