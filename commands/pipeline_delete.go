package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var pipelineDeleteCmd = &cobra.Command{
	Use:   "delete <id> [flags]",
	Short: `Delete a pipeline`,
	Example: heredoc.Doc(`
	$ glab pipeline delete 34
	$ glab pipeline delete 12,34,2
	`),
	Long:  ``,
	Run: deletePipeline,
}

func deletePipeline(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		pipelineID := strings.Trim(args[0], " ")
		gitlabClient, repo := git.InitGitlabClient()
		arrIds := strings.Split(strings.Trim(pipelineID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Pipeline #" + i2)
			pipeline, _ := gitlabClient.Pipelines.DeletePipeline(repo, manip.StringToInt(i2))
			if pipeline.StatusCode == 204 {
				color.Green.Println("Pipeline Deleted Successfully")
			} else if pipeline.StatusCode == 404 {
				er("Pipeline does not exist")
			} else {
				er("Could not complete request." +pipeline.Status)
			}
		}
		fmt.Println()
	} else {
		cmdErr(cmd, args)
	}
}

func init()  {
	pipelineCmd.AddCommand(pipelineDeleteCmd)
}