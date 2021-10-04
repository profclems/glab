package iostreams

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/shlex"
	"github.com/muesli/termenv"
)

type IOStreams struct {
	In     io.ReadCloser
	StdOut io.Writer
	StdErr io.Writer

	IsaTTY         bool //stdout is a tty
	IsErrTTY       bool //stderr is a tty
	IsInTTY        bool //stdin is a tty
	promptDisabled bool //disable prompting for input

	is256ColorEnabled bool

	pagerCommand string
	pagerProcess *os.Process

	spinner *spinner.Spinner

	backgroundColor string

	displayHyperlinks string
}

func Init() *IOStreams {
	stdoutIsTTY := IsTerminal(os.Stdout)
	stderrIsTTY := IsTerminal(os.Stderr)

	var pagerCommand string
	if glabPager, glabPagerExists := os.LookupEnv("GLAB_PAGER"); glabPagerExists {
		pagerCommand = glabPager
	} else {
		pagerCommand = os.Getenv("PAGER")
	}

	ioStream := &IOStreams{
		In:                os.Stdin,
		StdOut:            NewColorable(os.Stdout),
		StdErr:            NewColorable(os.Stderr),
		pagerCommand:      pagerCommand,
		IsaTTY:            stdoutIsTTY,
		IsErrTTY:          stderrIsTTY,
		is256ColorEnabled: Is256ColorSupported(),
		displayHyperlinks: "auto",
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
	return s.IsOutputTTY()
}

func (s *IOStreams) ColorEnabled() bool {
	return isColorEnabled() && s.IsaTTY && s.IsErrTTY
}

func (s *IOStreams) Is256ColorSupported() bool {
	return s.is256ColorEnabled
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

	_ = s.StdOut.(io.ReadCloser).Close()
	_, _ = s.pagerProcess.Wait()
	s.pagerProcess = nil
}

func (s *IOStreams) StartSpinner(format string, a ...interface{}) {
	if s.IsOutputTTY() {
		s.spinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(s.StdErr))
		if format != "" {
			s.spinner.Suffix = fmt.Sprintf(" "+format, a...)
		}
		s.spinner.Start()
	}
}

func (s *IOStreams) StopSpinner(format string, a ...interface{}) {
	if s.spinner != nil {
		s.spinner.Suffix = ""
		s.spinner.FinalMSG = fmt.Sprintf(format, a...)
		s.spinner.Stop()
		s.spinner = nil
	}
}

func (s *IOStreams) TerminalWidth() int {
	return TerminalWidth(s.StdOut)
}

//IsOutputTTY returns true if both stdout and stderr is TTY
func (s *IOStreams) IsOutputTTY() bool {
	return s.IsErrTTY && s.IsaTTY
}

func (s *IOStreams) ResolveBackgroundColor(style string) string {
	if style == "" {
		style = os.Getenv("GLAMOUR_STYLE")
	}
	if style != "" && style != "auto" {
		s.backgroundColor = style
		return style
	}
	if (!s.ColorEnabled()) ||
		(s.pagerProcess != nil) {
		s.backgroundColor = "none"
		return "none"
	}

	if termenv.HasDarkBackground() {
		s.backgroundColor = "dark"
		return "dark"
	}

	s.backgroundColor = "light"
	return "light"
}

func (s *IOStreams) BackgroundColor() string {
	if s.backgroundColor == "" {
		return "none"
	}
	return s.backgroundColor
}

func (s *IOStreams) SetDisplayHyperlinks(displayHyperlinks string) {
	s.displayHyperlinks = displayHyperlinks
}

func (s *IOStreams) DisplayHyperlinks() bool {
	switch s.displayHyperlinks {
	case "always":
		return true
	case "never":
		return false
	default:
		return s.IsaTTY
	}
}

func (s *IOStreams) MakeHyperlink(displayText, targetURL string) string {
	openSequence := fmt.Sprintf("\x1b]8;;%s\x1b\\", targetURL)
	closeSequence := "\x1b]8;;\x1b\\"

	return openSequence + displayText + closeSequence
}

func Test() (streams *IOStreams, in *bytes.Buffer, out *bytes.Buffer, errOut *bytes.Buffer) {
	in = &bytes.Buffer{}
	out = &bytes.Buffer{}
	errOut = &bytes.Buffer{}
	streams = &IOStreams{
		In:     ioutil.NopCloser(in),
		StdOut: out,
		StdErr: errOut,
	}
	return
}
