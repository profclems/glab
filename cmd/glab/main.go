package main

import (
	"errors"
	"fmt"
	"glab/internal/utils"
	"io"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"glab/commands"
	"glab/internal/config"
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
	if cmd, err := commands.Execute(); err != nil {
		printError(os.Stderr, err, cmd, debug)
		os.Exit(1)
	}
}

func initConfig() {
	config.SetGlobalPathDir()
	config.UseGlobalConfig = true

	if config.GetEnv("GITLAB_URI") == "NOTFOUND" || config.GetEnv("GITLAB_URI") == "OK" {
		config.SetEnv("GITLAB_URI", "https://gitlab.com")
	}
	if config.GetEnv("GIT_REMOTE_URL_VAR") == "NOTFOUND" || config.GetEnv("GIT_REMOTE_URL_VAR") == "OK" {
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

	re := regexp.MustCompile(`(?s){(.*)}`)
	m := re.FindAllStringSubmatch(err.Error(), -1)
	if len(m) != 0 {
		if len(m[0]) >= 1 {
			_, _ = fmt.Fprintln(out, m[0][1])
		}
	} else {
		_, _ = fmt.Fprintln(out, err)
	}

	var flagError *utils.FlagError
	if errors.As(err, &flagError) || strings.HasPrefix(err.Error(), "unknown command ") {
		if !strings.HasSuffix(err.Error(), "\n") {
			_, _ = fmt.Fprintln(out)
		}
		_, _ = fmt.Fprintln(out, cmd.UsageString())
	}
}
