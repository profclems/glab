package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/pkg/errors"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/spf13/cobra"
)

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
					Name: "https://gitlab.com/api/v4",
				}),
				cmd:   nil,
				debug: false,
			},
			wantOut: `error connecting to https://gitlab.com/api/v4
check your internet connection or status.gitlab.com or 'Run sudo gitlab-ctl status' on your server if self-hosted
`,
		},
		{
			name: "DNS error with debug",
			args: args{
				err: fmt.Errorf("DNS oopsie: %w", &net.DNSError{
					Name: "https://gitlab.com/api/v4",
				}),
				cmd:   nil,
				debug: true,
			},

			wantOut: `error connecting to https://gitlab.com/api/v4
lookup https://gitlab.com/api/v4: 
check your internet connection or status.gitlab.com or 'Run sudo gitlab-ctl status' on your server if self-hosted
`,
		},
		{
			name: "Cobra flag error",
			args: args{
				err:   &cmdutils.FlagError{Err: errors.New("unknown flag --foo")},
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

// Test started when the test binary is started
// and calls the main function
func TestGlab(t *testing.T) {
	main()
}
