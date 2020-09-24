package cmdutils

import (
	"fmt"
	"github.com/profclems/glab/internal/utils"
	"github.com/spf13/cobra"
)

// CmdErr prints "unknown command" error for unknown subcommands.
func CmdErr(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(utils.ColorableErr(cmd), "Error: Unknown command:")
	return cmd.Usage()
}

// IsSuccessful returns true if code falls within the success range
// indicating a request was successful
func IsSuccessful(code int) bool {
	// code 2xx (ie, 200 – 299) are considered successful responses
	// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status.
	if code >= 200 && code <= 299 {
		return true
	}
	return false
}
