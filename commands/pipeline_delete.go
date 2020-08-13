package commands

import (
	"fmt"
	"glab/internal/git"
	"glab/internal/manip"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var pipelineDeleteCmd = &cobra.Command{
	Use:   "delete <id> [flags]",
	Short: `Delete a pipeline`,
	Example: heredoc.Doc(`
	$ glab pipeline delete 34
	$ glab pipeline delete 12,34,2
	`),
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run:  deletePipeline,
}

func deletePipeline(cmd *cobra.Command, args []string) {
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo = r
	}
	if m, _ := cmd.Flags().GetString("status"); m != "" {
		l := &gitlab.ListProjectPipelinesOptions{}
		l.Status = gitlab.BuildState(gitlab.BuildStateValue(m))
		pipes, _, err := gitlabClient.Pipelines.ListProjectPipelines(repo, l)
		if err != nil {
			er(err)
		}
		for _, item := range pipes {
			pipeline, _ := gitlabClient.Pipelines.DeletePipeline(repo, item.ID)
			if pipeline.StatusCode == 204 {

				fmt.Println("Pipeline #" + strconv.Itoa(item.ID) + " Deleted Successfully")
			} else {
				er("Could not complete request." + pipeline.Status)
			}
		}

	} else {
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
					er("Could not complete request." + pipeline.Status)
				}
			}
			fmt.Println()
		} else {
			cmdErr(cmd, args)
		}
	}

}

func init() {
	pipelineDeleteCmd.Flags().StringP("status", "s", "", "delete pipelines by status: {running|pending|success|failed|canceled|skipped|created|manual}")

	pipelineCmd.AddCommand(pipelineDeleteCmd)
}
