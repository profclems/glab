package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	_isStdoutTerminal = false
	checkedTerminal   = false
)

type IOStreams struct {
	In     io.ReadCloser
	StdOut io.Writer
	StdErr io.Writer

	IsaTTY   bool
	IsErrTTY bool

	pagerCommand string
	pagerProcess *os.Process
}

func InitIOStream() *IOStreams {
	stdoutIsTTY := IsTerminal(os.Stdout)
	stderrIsTTY := IsTerminal(os.Stderr)

	var pagerCommand string
	if ghPager, ghPagerExists := os.LookupEnv("GH_PAGER"); ghPagerExists {
		pagerCommand = ghPager
	} else {
		pagerCommand = os.Getenv("PAGER")
	}

	ioStream := &IOStreams{
		In:           os.Stdin,
		StdOut:       NewColorable(os.Stdout),
		StdErr:       NewColorable(os.Stderr),
		pagerCommand: pagerCommand,
		IsaTTY:       stdoutIsTTY,
		IsErrTTY:     stderrIsTTY,
	}
	_isColorEnabled = isColorEnabled() && stdoutIsTTY

	return ioStream
}

func (s *IOStreams) SetPager(cmd string) {
	s.pagerCommand = cmd
}

func (s *IOStreams) StartPager() error {
	if s.pagerCommand == "" || s.pagerCommand == "cat" || !isStdoutTerminal() {
		return nil
	}

	pagerArgs, err := shlex.Split(s.pagerCommand)
	if err != nil {
		return err
	}

	pagerEnv := os.Environ()
	for i := len(pagerEnv) - 1; i >= 0; i-- {
		if strings.HasPrefix(pagerEnv[i], "PAGER=") {
			pagerEnv = append(pagerEnv[0:i], pagerEnv[i+1:]...)
		}
	}
	if _, ok := os.LookupEnv("LESS"); !ok {
		pagerEnv = append(pagerEnv, "LESS=FRX")
	}
	if _, ok := os.LookupEnv("LV"); !ok {
		pagerEnv = append(pagerEnv, "LV=-c")
	}

	pagerCmd := exec.Command(pagerArgs[0], pagerArgs[1:]...)
	pagerCmd.Env = pagerEnv
	pagerCmd.Stdout = s.StdOut
	pagerCmd.Stderr = s.StdErr
	pagedOut, err := pagerCmd.StdinPipe()
	if err != nil {
		return err
	}
	s.StdOut = pagedOut
	err = pagerCmd.Start()
	if err != nil {
		return err
	}
	s.pagerProcess = pagerCmd.Process
	return nil
}

func (s *IOStreams) StopPager() {
	if s.pagerProcess == nil {
		return
	}

	s.StdOut.(io.ReadCloser).Close()
	_, _ = s.pagerProcess.Wait()
	s.pagerProcess = nil
}

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

func IOTest() (*IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	return &IOStreams{
		In:     ioutil.NopCloser(in),
		StdOut: out,
		StdErr: errOut,
	}, in, out, errOut
}

// TODO: remove after making sure no function uses it
func ColorableOut(cmd *cobra.Command) io.Writer {
	out := cmd.OutOrStdout()
	if outFile, isFile := out.(*os.File); isFile {
		return NewColorable(outFile)
	}
	return out
}

// TODO: remove
func ColorableErr(cmd *cobra.Command) io.Writer {
	err := cmd.ErrOrStderr()
	if outFile, isFile := err.(*os.File); isFile {
		return NewColorable(outFile)
	}
	return err
}
