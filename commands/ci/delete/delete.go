package delete

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdDelete(f *cmdutils.Factory) *cobra.Command {
	var pipelineDeleteCmd = &cobra.Command{
		Use:   "delete <id> [flags]",
		Short: `Delete a CI pipeline`,
		Example: heredoc.Doc(`
	$ glab ci delete 34
	$ glab ci delete 12,34,2
	`),
		Long: ``,
		Args: cobra.ExactArgs(1),
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

			if m, _ := cmd.Flags().GetString("status"); m != "" {
				l := &gitlab.ListProjectPipelinesOptions{}
				l.Status = gitlab.BuildState(gitlab.BuildStateValue(m))
				pipes, err := api.ListProjectPipelines(apiClient, repo.FullName(), l)
				if err != nil {
					return err
				}
				for _, item := range pipes {
					err := api.DeletePipeline(apiClient, repo.FullName(), item.ID)
					if err != nil {
						return err
					}

					fmt.Fprintln(f.IO.StdOut, utils.RedCheck(), "Pipeline #"+strconv.Itoa(item.ID)+" Deleted Successfully")
				}

			} else {
				pipelineID := args[0]

				arrIds := strings.Split(strings.Trim(pipelineID, "[] "), ",")
				for _, i2 := range arrIds {
					fmt.Fprintln(f.IO.StdOut, "Deleting Pipeline #"+i2)
					err := api.DeletePipeline(apiClient, repo.FullName(), utils.StringToInt(i2))
					if err != nil {
						return err
					}

					fmt.Fprintln(f.IO.StdOut, utils.RedCheck(), "Pipeline #"+i2+" Deleted Successfully")
				}
				fmt.Println()
			}

			return nil

		},
	}

	pipelineDeleteCmd.Flags().StringP("status", "s", "", "delete pipelines by status: {running|pending|success|failed|canceled|skipped|created|manual}")

	return pipelineDeleteCmd
}
