package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/profclems/glab/commands"
	"github.com/profclems/glab/commands/alias/expand"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/help"
	"github.com/profclems/glab/commands/update"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/internal/run"
	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/spf13/cobra"
)

// version is set dynamically at build
var version string

// build is set dynamically at build
var build string

// debug is set dynamically at build and can be overridden by
// the configuration file or environment variable
// sets to "true" or "false" as string
var debugMode string
var debug bool // parsed boolean of debugMode

func main() {
	debug = debugMode == "true"

	cmdFactory := cmdutils.NewFactory()

	if glHostFromEnv := config.GetFromEnv("host"); glHostFromEnv != "" {
		glinstance.OverrideDefault(glHostFromEnv)
	}

	rootCmd := commands.NewCmdRoot(cmdFactory, version, build)

	cfg, err := cmdFactory.Config()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read configuration:  %s\n", err)
		os.Exit(2)
	}

	// Set Debug mode
	debugMode, _ = cfg.Get("", "debug")
	if debugSet, _ := strconv.ParseBool(debugMode); debugSet {
		debug = debugSet
	}

	if pager, _ := cfg.Get("", "pager"); pager != "" {
		cmdFactory.IO.SetPager(pager)
	}

	var expandedArgs []string
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	cmd, _, err := rootCmd.Traverse(expandedArgs)
	if err != nil || cmd == rootCmd {
		originalArgs := expandedArgs
		isShell := false
		expandedArgs, isShell, err = expand.ExpandAlias(cfg, os.Args, nil)
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

	// Override the default column separator of tableprinter
	tableprinter.DefaultSeparator = "  "

	rootCmd.SetArgs(expandedArgs)

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		printError(os.Stderr, err, cmd, debug)
		cmd.Print("\n")
		os.Exit(1)
	}

	if help.HasFailed() {
		os.Exit(1)
	}

	checkUpdate, _ := cfg.Get("", "check_update")
	if checkUpdate, err := strconv.ParseBool(checkUpdate); err == nil && checkUpdate {
		err = update.CheckUpdate(rootCmd, version, build, true)
		if err != nil && debug {
			printError(os.Stderr, err, rootCmd, debug)
		}
	}
	cmd.Print("\n")
}

func printError(out io.Writer, err error, cmd *cobra.Command, debug bool) {
	if err == cmdutils.SilentError {
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

	var flagError *cmdutils.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			_, _ = fmt.Fprintln(out)
		}
		_, _ = fmt.Fprintln(out, cmd.UsageString())
	}
}
