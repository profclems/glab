package root

import (
	"bytes"
	"github.com/profclems/glab/commands/cmdutils"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/profclems/glab/test"
	"github.com/stretchr/testify/assert"
)

var (
	glabBinaryPath = "../../bin/glab"
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	// Build a glab binary with test symbols. If the parent test binary was run
	// with coverage enabled, enable coverage on the child binary, too.
	var err error
	glabBinaryPath, err = filepath.Abs(os.ExpandEnv("$GOPATH/src/github.com/profclems/glab/test/testdata/glab"))
	if err != nil {
		log.Fatal(err)
	}
	testCmd := []string{"test", "-c", "-o", glabBinaryPath, "github.com/profclems/glab/cmd/glab"}
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
	repo := test.CopyTestRepo(log.New(os.Stderr, "", log.LstdFlags))
	if err := os.Chdir(repo); err != nil {
		log.Fatalf("Error chdir to test/testdata: %s", err)
	}
	code := m.Run()

	if err := os.Chdir(originalWd); err != nil {
		log.Fatalf("Error chdir to original working dir: %s", err)
	}
	os.Remove(glabBinaryPath)
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

func TestRootVersion(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	rootCmd := NewCmdRoot(&cmdutils.Factory{}, "v1.0.0", "2020-01-01")
	assert.Nil(t, rootCmd.Flag("version").Value.Set("true"))
	assert.Nil(t, rootCmd.Execute())

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Contains(t, out,"glab v1.0.0 (2020-01-01)")
}

func TestRootNoArg(t *testing.T) {
	cmd := exec.Command(glabBinaryPath)
	b, _ := cmd.CombinedOutput()
	assert.Contains(t, string(b), `GLab is an open source Gitlab Cli tool bringing GitLab to your command line`)
	assert.Contains(t, string(b), `USAGE
  glab <command> <subcommand> [flags]

CORE COMMANDS`)
}
