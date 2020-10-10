package pipeline

import (
	"github.com/profclems/glab/commands/cmdutils"
	ciCmd "github.com/profclems/glab/commands/pipeline/ci"
	pipeDeleteCmd "github.com/profclems/glab/commands/pipeline/delete"
	pipeListCmd "github.com/profclems/glab/commands/pipeline/list"
	pipeRunCmd "github.com/profclems/glab/commands/pipeline/run"
	pipeStatusCmd "github.com/profclems/glab/commands/pipeline/status"

	"github.com/spf13/cobra"
)

func NewCmdPipeline(f *cmdutils.Factory) *cobra.Command {
	var pipelineCmd = &cobra.Command{
		Use:     "pipeline <command> [flags]",
		Short:   `Manage pipelines`,
		Long:    ``,
		Aliases: []string{"pipe"},
	}

	cmdutils.EnableRepoOverride(pipelineCmd, f)

	pipelineCmd.AddCommand(ciCmd.NewCmdCI(f))
	pipelineCmd.AddCommand(pipeDeleteCmd.NewCmdDelete(f))
	pipelineCmd.AddCommand(pipeListCmd.NewCmdList(f))
	pipelineCmd.AddCommand(pipeStatusCmd.NewCmdStatus(f))
	pipelineCmd.AddCommand(pipeRunCmd.NewCmdRun(f))
	return pipelineCmd
}
