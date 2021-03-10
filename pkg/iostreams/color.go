package iostreams

import (
	"io"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var (
	_isColorEnabled = true
	checkedNoColor  = false
)

type ColorPalette struct {
	// Magenta outputs ANSI color if stdout is a tty
	Magenta func(string) string
	// Cyan outputs ANSI color if stdout is a tty
	Cyan func(string) string
	// Red outputs ANSI color if stdout is a tty
	Red func(string) string
	// Yellow outputs ANSI color if stdout is a tty
	Yellow func(string) string
	// Blue outputs ANSI color if stdout is a tty
	Blue func(string) string
	// Green outputs ANSI color if stdout is a tty
	Green func(string) string
	// Gray outputs ANSI color if stdout is a tty
	Gray func(string) string
	// Bold outputs ANSI color if stdout is a tty
	Bold func(string) string
}

func (s *IOStreams) Color() *ColorPalette {
	return &ColorPalette{
		Magenta: makeColorFunc("magenta"),
		Cyan:    makeColorFunc("cyan"),
		Red:     makeColorFunc("red"),
		Yellow:  makeColorFunc("yellow"),
		Blue:    makeColorFunc("blue"),
		Green:   makeColorFunc("green"),
		Gray:    makeColorFunc("black+h"),
		Bold:    makeColorFunc("default+b"),
	}
}

// NewColorable returns an output stream that handles ANSI color sequences on Windows
func NewColorable(out io.Writer) io.Writer {
	if outFile, isFile := out.(*os.File); isFile {
		return colorable.NewColorable(outFile)
	}
	return out
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
		_, _isColorEnabled = os.LookupEnv("NO_COLOR")
		_isColorEnabled = !_isColorEnabled // Revert the value NO_COLOR disables color

		if !_isColorEnabled {
			_isColorEnabled = os.Getenv("COLOR_ENABLED") == "1" || os.Getenv("COLOR_ENABLED") == "true"
		}
		checkedNoColor = true
	}
	return _isColorEnabled
}

func Is256ColorSupported() bool {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(term, "256") ||
		strings.Contains(term, "24bit") ||
		strings.Contains(term, "truecolor") ||
		strings.Contains(colorterm, "256") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}
