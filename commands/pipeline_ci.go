package commands

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var pipelineCICmd = &cobra.Command{
	Use:   "ci [command] [flags]",
	Short: `Work with GitLab CI pipelines and jobs`,
	Example: heredoc.Doc(`
	$ glab pipeline ci trace
	`),
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 1 {
			cmdErr(cmd, args)
		}
	},
}

func init() {
	pipelineCmd.AddCommand(pipelineCICmd)
}
