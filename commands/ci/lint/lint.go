package lint

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/git"
	"github.com/rsteube/carapace"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdLint(f *cmdutils.Factory) *cobra.Command {
	var pipelineCILintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Checks if your .gitlab-ci.yml file is valid.",
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
		$ glab ci lint  
		#=> Uses .gitlab-ci.yml in the current directory

		$ glab ci lint .gitlab-ci.yml

		$ glab ci lint path/to/.gitlab-ci.yml
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ".gitlab-ci.yml"
			if len(args) == 1 {
				path = args[0]
			}
			return lintRun(f, path)
		},
	}

	carapace.Gen(pipelineCILintCmd).PositionalCompletion(
		carapace.ActionFiles(".gitlab-ci.yml"),
	)

	return pipelineCILintCmd
}

func lintRun(f *cmdutils.Factory, path string) error {
	var err error
	out := f.IO.StdOut
	c := f.IO.Color()

	apiClient, err := f.HttpClient()
	if err != nil {
		return err
	}

	fmt.Fprintln(f.IO.StdErr, "Getting contents in", path)

	var content []byte
	var stdout bytes.Buffer

	if git.IsValidURL(path) {
		resp, err := http.Get(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(&stdout, resp.Body)
		if err != nil {
			return err
		}
		content = stdout.Bytes()
	} else {
		content, err = ioutil.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%s: no such file or directory", path)
			}
			return err
		}
	}

	fmt.Fprintln(f.IO.StdErr, "Validating...")

	lint, err := api.PipelineCILint(apiClient, string(content))
	if err != nil {
		return err
	}

	if lint.Status == "invalid" {
		fmt.Fprintln(out, c.Red(path+" is invalid"))
		for i, err := range lint.Errors {
			i++
			fmt.Fprintln(out, i, err)
		}
		return cmdutils.SilentError
	}
	fmt.Fprintln(out, c.GreenCheck(), "CI yml is Valid!")
	return nil
}
