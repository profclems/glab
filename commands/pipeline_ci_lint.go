package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// ciLintCmd represents the lint command
var pipelineCILintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Checks if your .gitlab-ci.yml file is valid.",
	Long:  ``,
	Example: heredoc.Doc(`
		$ glab pipeline ci lint  # Uses .gitlab-ci.yml in the current directory
		$ glab pipeline ci lint .gitlab-ci.yml
		$ glab pipeline ci lint path/to/.gitlab-ci.yml
	`),
	Run: func(cmd *cobra.Command, args []string) {
		path := ".gitlab-ci.yml"
		if len(args) == 1 {
			path = args[0]
		}
		fmt.Println("Getting contents in", path)

		content, err := ioutil.ReadFile(path)
		if !os.IsNotExist(err) && err != nil {
			log.Fatal(err)
		}
		fmt.Println("Validating...")
		lint, err := pipelineCILint(string(content))
		if err != nil {
			er(err)
			return
		}
		if lint.Status == "invalid" {
			color.Red.Println(path, "is invalid")
			for i, err := range lint.Errors {
				i++
				fmt.Println(i, err)
				return
			}
		}
		color.Green.Println("CI yml is Valid!")
		return
	},
}

func init() {
	pipelineCICmd.AddCommand(pipelineCILintCmd)
}
