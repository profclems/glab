package iostreams

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mgutz/ansi"
)

var (
	_isColorEnabled = true
	checkedNoColor  = false

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
