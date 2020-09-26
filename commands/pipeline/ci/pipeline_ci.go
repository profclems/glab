package ci

import (
	"github.com/profclems/glab/commands/cmdutils"
	ciTraceCmd "github.com/profclems/glab/commands/pipeline/ci/trace"
	ciViewCmd "github.com/profclems/glab/commands/pipeline/ci/view"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdCI(f *cmdutils.Factory) *cobra.Command {
	var pipelineCICmd = &cobra.Command{
		Use:   "ci [command] [flags]",
		Short: `Work with GitLab CI pipelines and jobs`,
		Example: heredoc.Doc(`
	$ glab pipeline ci trace
	`),
		Long: ``,
	}

	pipelineCICmd.AddCommand(ciTraceCmd.NewCmdTrace(f))
	pipelineCICmd.AddCommand(ciViewCmd.NewCmdView(f))
	return pipelineCICmd
}
