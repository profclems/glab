package test

import (
	"bytes"
	"fmt"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/profclems/glab/internal/run"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

var CachedTestFactory *cmdutils.Factory

// TODO copypasta from command package
type CmdOut struct {
	OutBuf, ErrBuf *bytes.Buffer
}

func (c CmdOut) String() string {
	return c.OutBuf.String()
}

func (c CmdOut) Stderr() string {
	return c.ErrBuf.String()
}

// OutputStub implements a simple utils.Runnable
type OutputStub struct {
	Out   []byte
	Error error
}

func (s OutputStub) Output() ([]byte, error) {
	if s.Error != nil {
		return s.Out, s.Error
	}
	return s.Out, nil
}

func (s OutputStub) Run() error {
	if s.Error != nil {
		return s.Error
	}
	return nil
}

type CmdStubber struct {
	Stubs []*OutputStub
	Count int
	Calls []*exec.Cmd
}

func InitCmdStubber() (*CmdStubber, func()) {
	cs := CmdStubber{}
	teardown := run.SetPrepareCmd(createStubbedPrepareCmd(&cs))
	return &cs, teardown
}

func (cs *CmdStubber) Stub(desiredOutput string) {
	// TODO maybe have some kind of command mapping but going simple for now
	cs.Stubs = append(cs.Stubs, &OutputStub{[]byte(desiredOutput), nil})
}

func (cs *CmdStubber) StubError(errText string) {
	// TODO support error types beyond CmdError
	stderrBuff := bytes.NewBufferString(errText)
	args := []string{"stub"} // TODO make more real?
	err := errors.New(errText)
	cs.Stubs = append(cs.Stubs, &OutputStub{Error: &run.CmdError{
		Stderr: stderrBuff,
		Args:   args,
		Err:    err,
	}})
}

func createStubbedPrepareCmd(cs *CmdStubber) func(*exec.Cmd) run.Runnable {
	return func(cmd *exec.Cmd) run.Runnable {
		cs.Calls = append(cs.Calls, cmd)
		call := cs.Count
		cs.Count += 1
		if call >= len(cs.Stubs) {
			panic(fmt.Sprintf("more execs than stubs. most recent call: %v", cmd))
		}
		// fmt.Printf("Called stub for `%v`\n", cmd) // Helpful for debugging
		return cs.Stubs[call]
	}
}

type T interface {
	Helper()
	Errorf(string, ...interface{})
}

func ExpectLines(t T, output string, lines ...string) {
	t.Helper()
	var r *regexp.Regexp
	for _, l := range lines {
		r = regexp.MustCompile(l)
		if !r.MatchString(output) {
			t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
			return
		}
	}
}

type fatalLogger interface {
	Fatal(...interface{})
}

func CopyTestRepo(log fatalLogger) string {
	rand.Seed(time.Now().UnixNano())
	dest, err := filepath.Abs(os.ExpandEnv("$GOPATH/src/github.com/profclems/glab/test/testdata-" + strconv.Itoa(int(rand.Uint64()))))
	if err != nil {
		log.Fatal(err)
	}
	src, err := filepath.Abs(os.ExpandEnv("$GOPATH/src/github.com/profclems/glab/test/testdata"))
	if err != nil {
		log.Fatal(err)
	}
	if err := copy.Copy(src, dest); err != nil {
		log.Fatal(err)
	}
	// Move the test.git dir into the expected path at .git
	if err := os.Rename(dest+"/test.git", dest+"/.git"); err != nil {
		log.Fatal(err)
	}
	// Move the test.glab-cli dir into the expected path at .glab-cli
	if err := os.Rename(dest+"/test.glab-cli", dest+"/.glab-cli"); err != nil {
		log.Fatal(err)
	}
	return dest
}

func FirstLine(output []byte) string {
	if i := bytes.IndexAny(output, "\n"); i >= 0 {
		return strings.ReplaceAll(string(output)[0:i], "PASS", "")
	}
	return string(output)
}

func Eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func StubFactory() *cmdutils.Factory {
	if CachedTestFactory != nil {
		return CachedTestFactory
	}
	conf := config.NewBlankConfig()
	CachedTestFactory = cmdutils.New(conf, nil)
	CachedTestFactory, _ = CachedTestFactory.NewClient("https://gitlab.com/glab-cli/test")

	return CachedTestFactory
}