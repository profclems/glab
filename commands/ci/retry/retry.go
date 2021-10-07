package retry

import (
	"fmt"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdRetry(f *cmdutils.Factory) *cobra.Command {
	var pipelineRetryCmd = &cobra.Command{
		Use:     "retry <job-id>",
		Short:   `Retry a CI job`,
		Aliases: []string{},
		Example: heredoc.Doc(`
	$ glab ci retry 871528
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

			jobID := utils.StringToInt(args[0])

			if jobID < 1 {
				fmt.Fprintln(f.IO.StdErr, "invalid job id:", args[0])
				return cmdutils.SilentError
			}

			job, err := api.RetryPipelineJob(apiClient, jobID, repo.FullName())
			if err != nil {
				return cmdutils.WrapError(err, fmt.Sprintf("Could not retry job with ID: %d", jobID))
			}
			fmt.Fprintln(f.IO.StdOut, "Retried job (id:", job.ID, "), status:", job.Status, ", ref:", job.Ref, ", weburl: ", job.WebURL, ")")

			return nil

		},
	}

	return pipelineRetryCmd
}
