// This package contains the old `glab pipeline ci` command which has been deprecated
// in favour of the `glab ci` command.
// This package is kept for backward compatibility but prints a deprecation warning
package legacyci

import (
	ciLintCmd "github.com/profclems/glab/commands/ci/lint"
	ciTraceCmd "github.com/profclems/glab/commands/ci/trace"
	ciViewCmd "github.com/profclems/glab/commands/ci/view"
	"github.com/profclems/glab/commands/cmdutils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdCI(f *cmdutils.Factory) *cobra.Command {
	var pipelineCICmd = &cobra.Command{
		Use:   "ci <command> [flags]",
		Short: `Work with GitLab CI pipelines and jobs`,
		Example: heredoc.Doc(`
	$ glab pipeline ci trace
	`),
	}

	pipelineCICmd.AddCommand(ciTraceCmd.NewCmdTrace(f, nil))
	pipelineCICmd.AddCommand(ciViewCmd.NewCmdView(f))
	pipelineCICmd.AddCommand(ciLintCmd.NewCmdLint(f))
	pipelineCICmd.Deprecated = "This command is deprecated. All the commands under it has been moved to `ci` or `pipeline` command. See https://github.com/profclems/glab/issues/372 for more info."
	pipelineCICmd.Hidden = true
	return pipelineCICmd
}
