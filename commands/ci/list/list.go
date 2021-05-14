package list

import (
	"fmt"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/ci/ciutils"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"
	"github.com/rsteube/carapace"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var pipelineListCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: `Get the list of CI pipelines`,
		Example: heredoc.Doc(`
	$ glab ci list
	$ glab ci list --state=failed
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var titleQualifier string

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			l := &gitlab.ListProjectPipelinesOptions{}
			l.Page = 1
			l.PerPage = 30

			if m, _ := cmd.Flags().GetString("status"); m != "" {
				l.Status = gitlab.BuildState(gitlab.BuildStateValue(m))
				titleQualifier = m
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

			pipes, err := api.ListProjectPipelines(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}

			title := utils.NewListTitle(fmt.Sprintf("%s pipeline", titleQualifier))
			title.RepoName = repo.FullName()
			title.Page = l.Page
			title.CurrentPageTotal = len(pipes)

			fmt.Fprintf(f.IO.StdOut, "%s\n%s\n", title.Describe(), ciutils.DisplayMultiplePipelines(f.IO.Color(), pipes, repo.FullName()))
			return nil
		},
	}
	pipelineListCmd.Flags().StringP("status", "s", "", "Get pipeline with status: {running|pending|success|failed|canceled|skipped|created|manual}")
	pipelineListCmd.Flags().StringP("orderBy", "o", "", "Order pipeline by <string>")
	pipelineListCmd.Flags().StringP("sort", "", "desc", "Sort pipeline by {asc|desc}. (Defaults to desc)")
	pipelineListCmd.Flags().IntP("page", "p", 1, "Page number")
	pipelineListCmd.Flags().IntP("per-page", "P", 30, "Number of items to list per page. (default 30)")

	carapace.Gen(pipelineListCmd).FlagCompletion(carapace.ActionMap{
		"orderBy": carapace.ActionValues("id", "status", "ref", "updated_at", "user_id"),
		"sort":    carapace.ActionValues("asc", "desc"),
		"status":  carapace.ActionValues("running", "pending", "success", "failed", "canceled", "skipped", "created", "manual"),
	})

	return pipelineListCmd
}
