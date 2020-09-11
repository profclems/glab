package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

// Test started when the test binary is started
// and calls the main function
func TestGlab(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	flag.Set("test", "../coverage-"+strconv.Itoa(int(rand.Uint64()))+".out")
	main()
}

func Test_printError(t *testing.T) {
	cmd := &cobra.Command{}

	type args struct {
		err   error
		cmd   *cobra.Command
		debug bool
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "generic error",
			args: args{
				err:   errors.New("the app exploded"),
				cmd:   nil,
				debug: false,
			},
			wantOut: "the app exploded\n",
		},
		{
			name: "DNS error",
			args: args{
				err: fmt.Errorf("DNS oopsie: %w", &net.DNSError{
					Name: config.GetEnv("GITLAB_URI") + "/api/v4",
				}),
				cmd:   nil,
				debug: false,
			},
			wantOut: `error connecting to ` + config.GetEnv("GITLAB_URI") + `/api/v4
check your internet connection or status.gitlab.com or 'Run sudo gitlab-ctl status' on your server if self-hosted
`,
		},
		{
			name: "Cobra flag error",
			args: args{
				err:   &utils.FlagError{Err: errors.New("unknown flag --foo")},
				cmd:   cmd,
				debug: false,
			},
			wantOut: "unknown flag --foo\n\nUsage:\n\n",
		},
		{
			name: "unknown Cobra command error",
			args: args{
				err:   errors.New("unknown command foo"),
				cmd:   cmd,
				debug: false,
			},
			wantOut: "unknown command foo\n\nUsage:\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printError(out, tt.args.err, tt.args.cmd, tt.args.debug)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printError() = %q, want %q", gotOut, tt.wantOut)
			}
		})
	}
}

func TestExpandAlias(t *testing.T) {
	err := config.SetAlias("test-co", "mr checkout")
	if err != nil {
		t.Error(err)
	}
	err = config.SetAlias("test-il", "issue list --author=\"$1\" --label=\"$2\"")
	if err != nil {
		t.Error(err)
	}
	err = config.SetAlias("test-ia", "issue list --author=\"$1\" --assignee=\"$1\"")
	if err != nil {
		t.Error(err)
	}
	for _, c := range []struct {
		Args         string
		ExpectedArgs []string
		Err          string
	}{
		{"glab test-co", []string{"mr", "checkout"}, ""},
		{"glab test-il", nil, `not enough arguments for alias: issue list --author="$1" --label="$2"`},
		{"glab test-il vilmibm", nil, `not enough arguments for alias: issue list --author="vilmibm" --label="$2"`},
		{"glab test-co 123", []string{"mr", "checkout", "123"}, ""},
		{"glab test-il vilmibm epic", []string{"issue", "list", `--author=vilmibm`, `--label=epic`}, ""},
		{"glab test-ia vilmibm", []string{"issue", "list", `--author=vilmibm`, `--assignee=vilmibm`}, ""},
		{"glab test-ia $coolmoney$", []string{"issue", "list", `--author=$coolmoney$`, `--assignee=$coolmoney$`}, ""},
		{"glab mr status", []string{"mr", "status"}, ""},
		{"glab test-il vilmibm epic -R vilmibm/testing", []string{"issue", "list", "--author=vilmibm", "--label=epic", "-R", "vilmibm/testing"}, ""},
		{"glab test-dne", []string{"test-dne"}, ""},
		{"glab", []string{}, ""},
		{"", []string{}, ""},
	} {
		var args []string
		if c.Args != "" {
			args = strings.Split(c.Args, " ")
		}

		out, err := expandAlias(args)

		if err == nil && c.Err != "" {
			t.Logf("expected error %s for %s", c.Err, c.Args)
			continue
		}

		if err != nil {
			eq(t, err.Error(), c.Err)
			continue
		}

		if len(out) == 0 && len(c.ExpectedArgs) == 0 {
			continue
		}
		eq(t, out, c.ExpectedArgs)
	}

	err = config.DeleteAlias("test-co")
	if err != nil {
		t.Log(err)
	}
	err = config.DeleteAlias("test-il")
	if err != nil {
		t.Log(err)
	}
	err = config.DeleteAlias("test-ia")
	if err != nil {
		t.Log(err)
	}
}

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func Test_initConfig(t *testing.T) {
	initConfig()
	config.UseGlobalConfig = true
	eq(t, config.GetEnv("GITLAB_URI"), "https://gitlab.com")
	eq(t, config.GetEnv("GIT_REMOTE_URL_VAR"), "origin")
	config.UseGlobalConfig = false
}
