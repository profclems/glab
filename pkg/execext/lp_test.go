package execext

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func winExecutable(s string) string {
	if runtime.GOOS == "windows" {
		return s
	}
	return ""
}

func TestLookPath(t *testing.T) {
	root, wderr := os.Getwd()
	if wderr != nil {
		t.Fatal(wderr)
	}
	defaultPath := os.Getenv("PATH")
	paths := []string{
		filepath.Join(root, "testdata", "nonexist"),
		filepath.Join(root, "testdata", "PATH"),
	}
	os.Setenv("PATH", strings.Join(paths, string(filepath.ListSeparator)))

	if err := os.Chdir(filepath.Join(root, "testdata", "cwd")); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc    string
		pathext string
		arg     string
		wants   string
		wantErr bool
	}{
		{
			desc:    "has extension",
			pathext: "",
			arg:     "git.exe",
			wants:   filepath.Join(root, "testdata", "PATH", "git.exe"),
			wantErr: false,
		},
		{
			desc:    "has path",
			pathext: "",
			arg:     filepath.Join("..", "PATH", "git"),
			wants:   filepath.Join("..", "PATH", "git"+winExecutable(".exe")),
			wantErr: false,
		},
		{
			desc:    "no extension",
			pathext: "",
			arg:     "git",
			wants:   filepath.Join(root, "testdata", "PATH", "git"+winExecutable(".exe")),
			wantErr: false,
		},
		{
			desc:    "has path+extension",
			pathext: "",
			arg:     filepath.Join("..", "PATH", "git.bat"),
			wants:   filepath.Join("..", "PATH", "git.bat"),
			wantErr: false,
		},
		{
			desc:    "no extension, PATHEXT",
			pathext: ".com;.bat",
			arg:     "git",
			wants:   filepath.Join(root, "testdata", "PATH", "git"+winExecutable(".bat")),
			wantErr: false,
		},
		{
			desc:    "has extension, PATHEXT",
			pathext: ".com;.bat",
			arg:     "git.exe",
			wants:   filepath.Join(root, "testdata", "PATH", "git.exe"),
			wantErr: false,
		},
		{
			desc:    "no extension, not found",
			pathext: "",
			arg:     "ls",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "has extension, not found",
			pathext: "",
			arg:     "ls.exe",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "no extension, PATHEXT, not found",
			pathext: ".com;.bat",
			arg:     "ls",
			wants:   "",
			wantErr: true,
		},
		{
			desc:    "has extension, PATHEXT, not found",
			pathext: ".com;.bat",
			arg:     "ls.exe",
			wants:   "",
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			os.Setenv("PATHEXT", tC.pathext)
			got, err := LookPath(tC.arg)

			if tC.wantErr != (err != nil) {
				t.Errorf("expects error: %v, got: %v", tC.wantErr, err)
			}
			if err != nil && !errors.Is(err, exec.ErrNotFound) {
				t.Errorf("expected exec.ErrNotFound; got %#v", err)
			}
			if got != tC.wants {
				t.Errorf("expected result %q, got %q", tC.wants, got)
			}
		})
	}
	_ = os.Setenv("PATH", defaultPath)
}
