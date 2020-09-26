package pipeline

import (
	"github.com/profclems/glab/commands/cmdutils"
	ciCmd "github.com/profclems/glab/commands/pipeline/ci"

	"github.com/spf13/cobra"
)

func NewCmdPipeline(f *cmdutils.Factory) *cobra.Command {
	var pipelineCmd = &cobra.Command{
		Use:     "pipeline <command> [flags]",
		Short:   `Manage pipelines`,
		Long:    ``,
		Aliases: []string{"pipe"},
	}

	pipelineCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	pipelineCmd.AddCommand(ciCmd.NewCmdCI(f))
	return pipelineCmd
}
