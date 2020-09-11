package main

import (
	"errors"
	"fmt"
	"github.com/profclems/glab/internal/run"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/profclems/glab/commands"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
)

// Version is set at build
var version string

// build is set at build
var build string

// usage mode is set at build to either "dev" or "prod" depending how binary is created
var usageMode string
var debug bool

func main() {
	commands.Version = version
	commands.Build = build

	initConfig()
	if usageMode == "dev" {
		debug = true
	}

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	cmd, _, err := commands.RootCmd.Traverse(expandedArgs)
	if err != nil || cmd == commands.RootCmd {
		originalArgs := expandedArgs
		isShell := false
		expandedArgs, isShell, err = commands.ExpandAlias(os.Args, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Failed to process alias: %s\n", err)
			os.Exit(2)
		}

		if debug {
			fmt.Printf("%v -> %v\n", originalArgs, expandedArgs)
		}

		if isShell {
			externalCmd := exec.Command(expandedArgs[0], expandedArgs[1:]...)
			externalCmd.Stderr = os.Stderr
			externalCmd.Stdout = os.Stdout
			externalCmd.Stdin = os.Stdin
			preparedCmd := run.PrepareCmd(externalCmd)

			err = preparedCmd.Run()
			if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					os.Exit(ee.ExitCode())
				}

				fmt.Fprintf(os.Stdout, "failed to run external command: %s", err)
				os.Exit(3)
			}

			os.Exit(0)
		}
	}

	commands.RootCmd.SetArgs(expandedArgs)

	if cmd, err := commands.RootCmd.ExecuteC(); err != nil {
		printError(os.Stderr, err, cmd, debug)
		os.Exit(1)
	}
}

func initConfig() {
	config.SetGlobalPathDir()
	config.UseGlobalConfig = true

	if config.GetEnv("GITLAB_URI") == "" {
		config.SetEnv("GITLAB_URI", "https://gitlab.com")
	}
	if config.GetEnv("GIT_REMOTE_URL_VAR") == "" {
		config.SetEnv("GIT_REMOTE_URL_VAR", "origin")
	}

	config.UseGlobalConfig = false
}

func printError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	if err == utils.SilentError {
		return
	}

	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		_, _ = fmt.Fprintf(out, "error connecting to %s\n", dnsError.Name)
		if debug {
			_, _ = fmt.Fprintln(out, dnsError)
		}
		_, _ = fmt.Fprintln(out, "check your internet connection or status.gitlab.com or 'Run sudo gitlab-ctl status' on your server if self-hosted")
		return
	}
	_, _ = fmt.Fprintln(out, err)

	var flagError *utils.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			_, _ = fmt.Fprintln(out)
		}
		_, _ = fmt.Fprintln(out, cmd.UsageString())
	}
}
