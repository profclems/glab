package lint

import (
	"testing"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/alecthomas/assert"
	"github.com/profclems/glab/commands/cmdtest"
)

func TestMain(m *testing.M) {
	cmdtest.InitTest(m, "ci_lint_test")
}

func Test_pipelineCILint(t *testing.T) {
	io, _, stdout, stderr := iostreams.IOTest()
	fac := cmdtest.StubFactory("")
	fac.IO = io
	fac.IO.StdErr = stderr
	fac.IO.StdOut = stdout

	tests := []struct {
		Name    string
		Args    string
		StdOut  string
		StdErr  string
		WantErr error
	}{
		{
			Name:   "with no path specified",
			Args:   "",
			StdOut: "✓ CI yml is Valid!\n",
			StdErr: "Getting contents in .gitlab-ci.yml\nValidating...\n",
		},
		{
			Name:   "with path specified as url",
			Args:   "https://gitlab.com/profclems/glab/-/raw/trunk/.gitlab-ci.yml",
			StdOut: "✓ CI yml is Valid!\n",
			StdErr: "Getting contents in https://gitlab.com/profclems/glab/-/raw/trunk/.gitlab-ci.yml\nValidating...\n",
		},
	}

	cmd := NewCmdLint(fac)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := cmdtest.RunCommand(cmd, test.Args)
			if err != nil {
				if test.WantErr == nil {
					t.Fatal(err)
				}
				assert.Equal(t, err, test.WantErr)
			}
			assert.Equal(t, test.StdErr, stderr.String())
			assert.Equal(t, test.StdOut, stdout.String())
			stdout.Reset()
			stderr.Reset()
		})
	}
}

func Test_lintRun(t *testing.T) {
	io, _, stdout, stderr := iostreams.IOTest()
	fac := cmdtest.StubFactory("")
	fac.IO = io
	fac.IO.StdErr = stderr
	fac.IO.StdOut = stdout

	tests := []struct {
		name    string
		path    string
		StdOut  string
		StdErr  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "with invalid path specified",
			path:    "WRONG_PATH",
			StdOut:  "",
			StdErr:  "Getting contents in WRONG_PATH\n",
			wantErr: true,
			errMsg:  "WRONG_PATH: no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lintRun(fac, tt.path)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("lintRun() error = %v, wantErr %v", err, tt.wantErr)
				}
				assert.Equal(t, tt.errMsg, err.Error())
			}

			assert.Equal(t, tt.StdErr, stderr.String())
			assert.Equal(t, tt.StdOut, stdout.String())
		})
	}
}
