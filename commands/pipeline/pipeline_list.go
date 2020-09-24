package pipeline

import (
	"github.com/profclems/glab/internal/git"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var pipelineListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: `Get the list of pipelines`,
	Example: heredoc.Doc(`
	$ glab pipeline list
	$ glab pipeline list --state=failed
	`),
	Long: ``,
	Args: cobra.ExactArgs(0),
	Run:  listPipelines,
}

func listPipelines(cmd *cobra.Command, args []string) {
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo, _ = fixRepoNamespace(r)
	}
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
	if p, _ := cmd.Flags().GetInt("page"); p != 0 {
		l.Page = p
	}
	if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
		l.PerPage = p
	}

	pipes, _, err := gitlabClient.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		er(err)
	}

	displayMultiplePipelines(pipes, repo)
}

func init() {
	pipelineListCmd.Flags().StringP("status", "s", "", "Get pipeline with status: {running|pending|success|failed|canceled|skipped|created|manual}")
	pipelineListCmd.Flags().StringP("orderBy", "o", "", "Order pipeline by <string>")
	pipelineListCmd.Flags().StringP("sort", "", "desc", "Sort pipeline by {asc|desc}. (Defaults to desc)")
	pipelineListCmd.Flags().IntP("page", "p", 1, "Page number")
	pipelineListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	pipelineCmd.AddCommand(pipelineListCmd)
}
