package lint

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/profclems/glab/internal/git"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdLint(f *cmdutils.Factory) *cobra.Command {
	var pipelineCILintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Checks if your .gitlab-ci.yml file is valid.",
		Long:  ``,
		Example: heredoc.Doc(`
		$ glab pipeline ci lint  # Uses .gitlab-ci.yml in the current directory
		$ glab pipeline ci lint .gitlab-ci.yml
		$ glab pipeline ci lint path/to/.gitlab-ci.yml
	`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := f.IO.StdOut

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			path := ".gitlab-ci.yml"
			if len(args) == 1 {
				path = args[0]
			}

			fmt.Fprintln(out, "Getting contents in", path)

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
			}

			if !os.IsNotExist(err) && err != nil {
				return err
			}

			fmt.Fprintln(out, "Validating...")

			lint, err := api.PipelineCILint(apiClient, string(content))
			if err != nil {
				return err
			}

			if lint.Status == "invalid" {
				fmt.Fprintln(out, utils.Red(path+" is invalid"))
				for i, err := range lint.Errors {
					i++
					fmt.Fprintln(out, i, err)
				}
				return nil
			}
			fmt.Fprintln(out, utils.GreenCheck(), "CI yml is Valid!")
			return nil
		},
	}

	return pipelineCILintCmd
}
