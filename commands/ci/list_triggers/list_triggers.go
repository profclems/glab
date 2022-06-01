package trigger

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/ci/ciutils"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdRun(f *cmdutils.Factory) *cobra.Command {
	var pipelineRunCmd = &cobra.Command{
		Use:     "list-triggers [flags]",
		Short:   `List triggers of a pipeline in the CI`,
		Aliases: []string{"lst"},
		Example: heredoc.Doc(`
	$ glab ci list-triggers
	$ glab ci lst
	$ glab ci list-triggers --page 0
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			c := &gitlab.ListPipelineTriggersOptions{
				PerPage: 30,
				Page:    0,
			}

			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				c.PerPage = p
			}

			if p, err := cmd.Flags().GetInt("page"); err != nil {
				c.Page = p
			}

			triggers, err := api.ListPipelineTriggers(apiClient, repo.FullName(), c)
			if err != nil {
				return err
			}

			title := utils.NewListTitle(fmt.Sprintf("[%s] Pipeline trigger", repo.FullName()))
			title.RepoName = repo.FullName()
			title.Page = c.Page
			title.CurrentPageTotal = len(triggers)

			fmt.Fprintf(f.IO.StdOut, "%s\n%s\n", title.Describe(), ciutils.DisplayMultipleTriggers(f.IO, triggers))

			return nil
		},
	}
	pipelineRunCmd.Flags().IntP("page", "p", 0, "Which page of the pipeline triggers to return (default 0)")
	pipelineRunCmd.Flags().IntP("per-page", "P", 30, "Number of items to list per page. (default 30)")

	return pipelineRunCmd
}
