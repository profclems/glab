package iostreams

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/profclems/glab/pkg/execext"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	_isStdoutTerminal = false
	checkedTerminal   = false
)

func isStdoutTerminal() bool {
	if !checkedTerminal {
		_isStdoutTerminal = IsTerminal(os.Stdout)
		checkedTerminal = true
	}
	return _isStdoutTerminal
}

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

func isCygwinTerminal(w io.Writer) bool {
	if f, isFile := w.(*os.File); isFile {
		return IsCygwinTerminal(f)
	}
	return false
}

func TerminalWidth(out io.Writer) int {
	defaultWidth := 80

	if w, _, err := TerminalSize(out); err == nil {
		return w
	}

	if isCygwinTerminal(out) {
		tputExe, err := execext.LookPath("tput")
		if err != nil {
			return defaultWidth
		}
		tputCmd := exec.Command(tputExe, "cols")
		tputCmd.Stdin = os.Stdin
		if out, err := tputCmd.Output(); err == nil {
			if w, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
				return w
			}
		}
	}

	return defaultWidth
}
