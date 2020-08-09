package commands

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var pipelineListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: `Get the list of pipelines`,
	Example: heredoc.Doc(`
	$ glab pipeline list
	$ glab pipeline list --state=failed
	`),
	Long:  ``,
	Run: listPipelines,
}
func listPipelines(cmd *cobra.Command, args []string) {
	gitlabClient, repo := git.InitGitlabClient()
	l := &gitlab.ListProjectPipelinesOptions{}
	if m, _ := cmd.Flags().GetString("status"); m != "" {
		l.Status = gitlab.BuildState(gitlab.BuildStateValue(m))
	}
	if m, _ := cmd.Flags().GetString("orderBy"); m != "" {
		l.OrderBy = gitlab.String(m)
	}
	if m, _ := cmd.Flags().GetString("sort"); m != "" {
		l.Sort = gitlab.String(m)
	}
	pipes, _, err := gitlabClient.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		er(err)
	}

	displayMultiplePipelines(pipes)
}

func init() {
	pipelineListCmd.Flags().StringP("status", "s", "", "Get pipeline with status: {running|pending|success|failed|canceled|skipped|created|manual}")
	pipelineListCmd.Flags().StringP("orderBy", "o", "", "Order pipeline by <string>")
	pipelineListCmd.Flags().StringP("sort", "", "desc", "Sort pipeline by {asc|desc}. (Defaults to desc)")
	pipelineCmd.AddCommand(pipelineListCmd)
}