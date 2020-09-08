package commands

import (
	"bytes"
	"fmt"
	"io"
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

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

var (
	glabBinaryPath = "../../bin/glab"
)

func Test_isLatestVersion(t *testing.T) {
	type args struct {
		latestVersion  string
		currentVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "latest is newer",
			args: args{"v1.10.0", "v1.9.1"},
			want: true,
		},
		{
			name: "latest is current",
			args: args{"v1.9.2", "v1.9.2"},
			want: false,
		},
		{
			name: "latest is older",
			args: args{"v1.9.0", "v1.9.2-pre.1"},
			want: false,
		},
		{
			name: "current is prerelease",
			args: args{"v1.9.0", "v1.9.0-pre.1"},
			want: true,
		},
		{
			name: "latest is older (against prerelease)",
			args: args{"v1.9.0", "v1.10.0-pre.1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLatestVersion(tt.args.latestVersion, tt.args.currentVersion); got != tt.want {
				t.Errorf("isLatestVersion(%s, %s) = %v, want %v",
					tt.args.latestVersion, tt.args.currentVersion, got, tt.want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	// Build a lab binary with test symbols. If the parent test binary was run
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
	repo := copyTestRepo(log.New(os.Stderr, "", log.LstdFlags))
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

	RootCmd.Flag("version").Value.Set("true")
	RootCmd.Run(RootCmd, nil)

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

	assert.Contains(t, out, fmt.Sprintf("glab version %s (%s)", Version, Build))
	assert.Contains(t, out, "git version")
	assert.Contains(t, out, "Made with ❤ by Clement Sam <clementsam75@gmail.com> and contributors")
}

type fatalLogger interface {
	Fatal(...interface{})
}

func copyTestRepo(log fatalLogger) string {
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

func TestRootNoArg(t *testing.T) {
	cmd := exec.Command(glabBinaryPath)
	b, _ := cmd.CombinedOutput()
	assert.Contains(t, string(b), `GLab is an open source Gitlab Cli tool bringing GitLab to your command line`)
	assert.Contains(t, string(b), `Usage:
  glab <command> <subcommand> [flags]
  glab [command]`)
}

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func firstLine(output []byte) string {
	if i := bytes.IndexAny(output, "\n"); i >= 0 {
		return strings.ReplaceAll(string(output)[0:i], "PASS", "")
	}
	return string(output)
}