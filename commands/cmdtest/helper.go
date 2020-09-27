package cmdtest

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/shlex"
	"github.com/otiai10/copy"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/test"
	"github.com/spf13/cobra"
)

var (
	GlabBinaryPath    = "../../bin/glab"
	CachedTestFactory *cmdutils.Factory
	TestRepo          = ""
)

type fatalLogger interface {
	Fatal(...interface{})
}

func InitTest(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	// Build a glab binary with test symbols. If the parent test binary was run
	// with coverage enabled, enable coverage on the child binary, too.
	var err error
	GlabBinaryPath, err = filepath.Abs(os.ExpandEnv("$GOPATH/src/github.com/profclems/glab/test/testdata/glab"))
	if err != nil {
		log.Fatal(err)
	}
	testCmd := []string{"test", "-c", "-o", GlabBinaryPath, "github.com/profclems/glab/cmd/glab"}
	if coverMode := testing.CoverMode(); coverMode != "" {
		testCmd = append(testCmd, "-covermode", coverMode, "-coverpkg", "./...")
	}
	if out, err := exec.Command("go", testCmd...).CombinedOutput(); err != nil {
		log.Fatalf("Error building glab test binary: %s (%s)", string(out), err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Make a copy of the testdata Git test project and chdir to it.
	repo := CopyTestRepo(log.New(os.Stderr, "", log.LstdFlags))
	TestRepo = repo
	if err := os.Chdir(repo); err != nil {
		log.Fatalf("Error chdir to test/testdata: %s", err)
	}
	code := m.Run()

	if err := os.Chdir(originalWd); err != nil {
		log.Fatalf("Error chdir to original working dir: %s", err)
	}
	err = os.Remove(GlabBinaryPath)
	if err != nil {
		log.Fatal(err)
	}
	testdirs, err := filepath.Glob(os.ExpandEnv("$GOPATH/src/github.com/profclems/glab/test/testdata-*"))
	if err != nil {
		log.Printf("Error listing glob test/testdata-*: %s", err)
	}
	for _, dir := range testdirs {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Printf("Error removing dir %s: %s", dir, err)
		}
	}

	os.Exit(code)
}

func RunCommand(cmd *cobra.Command, cli string) (*test.CmdOut, error) {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	argv, err := shlex.Split(cli)
	if err != nil {
		return nil, err
	}
	cmd.SetArgs(argv)

	_, err = cmd.ExecuteC()

	return &test.CmdOut{
		OutBuf: &stdout,
		ErrBuf: &stderr,
	}, err
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
