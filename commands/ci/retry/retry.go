package retry

import (
	"fmt"
	"time"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/git"
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
	$ glab ci retry 871528    # retries a specific job, 871528
	$ glab ci retry           # retries most recent pipeline, if retry is necessary
	$ glab ci retry --follow  # continues to retry most recent pipeline, until interrupted
`),
		Long: ``,
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

			for i := range args {
				jobID := utils.StringToInt(args[i])

				if jobID < 1 {
					fmt.Fprintln(f.IO.StdErr, "invalid job id:", args[0])
					return cmdutils.SilentError
				}

				job, err := api.RetryPipelineJob(apiClient, jobID, repo.FullName())
				if err != nil {
					return cmdutils.WrapError(err, fmt.Sprintf("Could not retry job with ID: %d", jobID))
				}
				fmt.Fprintln(f.IO.StdOut, "Retried job (id:", job.ID, "), status:", job.Status, ", ref:", job.Ref, ", weburl: ", job.WebURL, ")")
			}
			if len(args) > 0 {
				// jobs specified on command line are retried, nothing more to do
				return nil
			}

			// retry all failed jobs in pipeline

			follow, _ := cmd.Flags().GetBool("follow")
			branch, _ := cmd.Flags().GetString("branch")
			if branch == "" {
				branch, err = git.CurrentBranch()
				if err != nil {
					return err
				}
			}

			attempts := map[int]int{} // key is pipeline id, value is how may retries

			for i := 0; i == 0 || follow; i++ {
				if i > 0 {
					// pause for retries triggered by prior iteration
					time.Sleep(30 * time.Minute)
				}

				lastPipeline, err := api.GetLastPipeline(apiClient, repo.FullName(), branch)
				if err != nil {
					if follow {
						continue
					}
					fmt.Fprintf(f.IO.StdOut, "No pipelines running or available on %s branch\n", branch)
					return err
				}

				switch lastPipeline.Status {
				case "canceled", "pending", "success", "skipped":
					// nothing to retry
					continue

				default: // "running", "failed", "created"
					// look for any failed jobs
					failed := false
					jobs, err := api.GetPipelineJobs(apiClient, lastPipeline.ID, repo.FullName())
					if err != nil {
						return err
					}
					for j := range jobs {
						if jobs[j].Status == "failed" {
							if jobs[j].AllowFailure {
								fmt.Fprintf(f.IO.StdErr, "failed job (%s) allows failure, ignoring", jobs[i].WebURL)
								continue
							}

							failed = true
							break
						}
					}
					if !failed {
						continue // continue main loop, nothing to retry
					}
				}

				count := attempts[lastPipeline.ID]
				if count >= 3 {
					fmt.Fprintf(f.IO.StdErr, "giving up on pipeline (%d), too many retries (%d)", lastPipeline.ID, count)
					continue
				}
				attempts[lastPipeline.ID] = count + 1

				fmt.Fprintf(f.IO.StdOut, "retrying pipeline (%s)", lastPipeline.WebURL)
				_, err = api.RetryPipeline(apiClient, lastPipeline.ID, repo.FullName())
				if err != nil {
					fmt.Fprintf(f.IO.StdErr, "failed to retry pipeline (%s): %+v", lastPipeline.WebURL, err)
				}
			}

			return nil
		},
	}

	pipelineRetryCmd.Flags().StringP("branch", "b", "", "Retry latest pipeline associated with branch. (Default is current branch)")
	pipelineRetryCmd.Flags().BoolP("follow", "f", false, "Retry when needed, until interrupted.")

	return pipelineRetryCmd
}
