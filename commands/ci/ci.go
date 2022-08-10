package ci

import (
	jobArtifactCmd "github.com/profclems/glab/commands/ci/artifact"
	pipeDeleteCmd "github.com/profclems/glab/commands/ci/delete"
	legacyCICmd "github.com/profclems/glab/commands/ci/legacyci"
	ciLintCmd "github.com/profclems/glab/commands/ci/lint"
	pipeListCmd "github.com/profclems/glab/commands/ci/list"
	pipeRetryCmd "github.com/profclems/glab/commands/ci/retry"
	pipeRunCmd "github.com/profclems/glab/commands/ci/run"
	pipeStatusCmd "github.com/profclems/glab/commands/ci/status"
	"github.com/profclems/glab/commands/ci/tally"
	ciTraceCmd "github.com/profclems/glab/commands/ci/trace"
	ciViewCmd "github.com/profclems/glab/commands/ci/view"
	"github.com/profclems/glab/commands/cmdutils"

	"github.com/spf13/cobra"
)

func NewCmdCI(f *cmdutils.Factory) *cobra.Command {
	var ciCmd = &cobra.Command{
		Use:     "ci <command> [flags]",
		Short:   `Work with GitLab CI pipelines and jobs`,
		Long:    ``,
		Aliases: []string{"pipe", "pipeline"},
	}

	cmdutils.EnableRepoOverride(ciCmd, f)

	ciCmd.AddCommand(legacyCICmd.NewCmdCI(f))
	ciCmd.AddCommand(ciTraceCmd.NewCmdTrace(f, nil))
	ciCmd.AddCommand(ciViewCmd.NewCmdView(f))
	ciCmd.AddCommand(ciLintCmd.NewCmdLint(f))
	ciCmd.AddCommand(pipeDeleteCmd.NewCmdDelete(f))
	ciCmd.AddCommand(pipeListCmd.NewCmdList(f))
	ciCmd.AddCommand(pipeStatusCmd.NewCmdStatus(f))
	ciCmd.AddCommand(pipeRetryCmd.NewCmdRetry(f))
	ciCmd.AddCommand(pipeRunCmd.NewCmdRun(f))
	ciCmd.AddCommand(jobArtifactCmd.NewCmdRun(f))
	ciCmd.AddCommand(tally.NewCmdTally(f))
	return ciCmd
}
