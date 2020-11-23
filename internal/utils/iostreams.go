package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
)

type IOStreams struct {
	In     io.ReadCloser
	StdOut io.Writer
	StdErr io.Writer

	IsaTTY         bool //stdout is a tty
	IsErrTTY       bool //stderr is a tty
	IsInTTY        bool //stdin is a tty
	promptDisabled bool //disable prompting for input

	pagerCommand string
	pagerProcess *os.Process
}

func InitIOStream() *IOStreams {
	stdoutIsTTY := IsTerminal(os.Stdout)
	stderrIsTTY := IsTerminal(os.Stderr)

	var pagerCommand string
	if glabPager, glabPagerExists := os.LookupEnv("GLAB_PAGER"); glabPagerExists {
		pagerCommand = glabPager
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

	if stdin, ok := ioStream.In.(*os.File); ok {
		ioStream.IsInTTY = IsTerminal(stdin)
	}

	_isColorEnabled = isColorEnabled() && stdoutIsTTY && stderrIsTTY

	return ioStream
}

func (s *IOStreams) PromptEnabled() bool {
	if s.promptDisabled {
		return false
	}
	return s.IsInTTY && s.IsaTTY
}

func (s *IOStreams) ColorEnabled() bool {
	return isColorEnabled() && s.IsaTTY && s.IsErrTTY
}

func (s *IOStreams) SetPrompt(promptDisabled string) {
	if promptDisabled == "true" || promptDisabled == "1" {
		s.promptDisabled = true
	} else if promptDisabled == "false" || promptDisabled == "0" {
		s.promptDisabled = false
	}
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

func (s *IOStreams) TerminalWidth() int {
	return TerminalWidth(s.StdOut)
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
