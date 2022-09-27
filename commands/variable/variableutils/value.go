package variableutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/iostreams"
)

func GetValue(value string, io *iostreams.IOStreams, args []string) (string, error) {
	if value != "" {
		return value, nil
	} else if len(args) == 2 {
		return args[1], nil
	}

	if io.IsInTTY {
		return "", &cmdutils.FlagError{Err: errors.New("no value specified but nothing on STDIN")}
	}

	// read value from STDIN if not provided
	defer io.In.Close()
	stdinValue, err := ioutil.ReadAll(io.In)
	if err != nil {
		return "", fmt.Errorf("failed to read value from STDIN: %w", err)
	}
	return strings.TrimSpace(string(stdinValue)), nil
}
