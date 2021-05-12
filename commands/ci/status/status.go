package status

import (
	"fmt"
	"time"

	"github.com/profclems/glab/api"
	ciTraceCmd "github.com/profclems/glab/commands/ci/trace"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/pkg/git"
	"github.com/profclems/glab/pkg/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdStatus(f *cmdutils.Factory) *cobra.Command {
	var pipelineStatusCmd = &cobra.Command{
		Use:     "status [flags]",
		Short:   `View a running CI pipeline on current or other branch specified`,
		Aliases: []string{"stats"},
		Example: heredoc.Doc(`
	$ glab ci status --live
	$ glab ci status --compact // more compact view
	$ glab ci status --branch=master   // Get pipeline for master branch
	$ glab ci status   // Get pipeline for current branch
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			branch, _ := cmd.Flags().GetString("branch")
			live, _ := cmd.Flags().GetBool("live")
			compact, _ := cmd.Flags().GetBool("compact")

			if branch == "" {
				branch, err = git.CurrentBranch()
				if err != nil {
					return err
				}
			}
			l := &gitlab.ListProjectPipelinesOptions{
				Ref:  gitlab.String(branch),
				Sort: gitlab.String("desc"),
			}
			l.Page = 1
			l.PerPage = 1

			pipes, err := api.GetPipelines(apiClient, l, repo.FullName())
			if err != nil {
				return err
			}

			if len(pipes) > 0 {
				runningPipeline := pipes[0]
				isRunning := true
				retry := "Exit"
				writer := uilive.New()

				// start listening for updates and render
				writer.Start()
				defer writer.Stop()
				for isRunning {
					jobs, err := api.GetPipelineJobs(apiClient, runningPipeline.ID, repo.FullName())
					if err != nil {
						return err
					}
					for _, job := range jobs {
						end := time.Now()
						if job.FinishedAt != nil {
							end = *job.FinishedAt
						}
						duration := utils.FmtDuration(end.Sub(*job.StartedAt))
						var status string
						switch s := job.Status; s {
						case "failed":
							if job.AllowFailure {
								status = c.Yellow(s)
							} else {
								status = c.Red(s)
							}
						case "success":
							status = c.Green(s)
						default:
							status = c.Gray(s)
						}
						//fmt.Println(job.Tag)
						if compact {
							fmt.Fprintf(writer, "(%s) • %s [%s]\n", status, job.Name, job.Stage)
						} else {
							fmt.Fprintf(writer, "(%s) • %s\t%s\t\t%s\n", status, c.Gray(duration), job.Stage, job.Name)
						}
					}

					if !compact {
						fmt.Fprintf(writer.Newline(), "\n%s\n", runningPipeline.WebURL)
						fmt.Fprintf(writer.Newline(), "SHA: %s\n", runningPipeline.SHA)
					}
					fmt.Fprintf(writer.Newline(), "Pipeline State: %s\n\n", runningPipeline.Status)

					// break loop if output is a TTY to avoid prompting
					if !f.IO.IsOutputTTY() || !f.IO.PromptEnabled() {
						break
					}
					if runningPipeline.Status == "running" && live {
						pipes, err = api.GetPipelines(apiClient, l, repo.FullName())
						if err != nil {
							return err
						}
						runningPipeline = pipes[0]
					} else {
						prompt := &survey.Select{
							Message: "Choose an action:",
							Options: []string{"View Logs", "Retry", "Exit"},
							Default: "Exit",
						}
						_ = survey.AskOne(prompt, &retry)
						if retry != "" && retry != "Exit" {
							if retry == "View Logs" {
								isRunning = false
							} else {
								_, err = api.RetryPipeline(apiClient, runningPipeline.ID, repo.FullName())
								if err != nil {
									return err
								}
								pipes, err = api.GetPipelines(apiClient, l, repo.FullName())
								if err != nil {
									return err
								}
								runningPipeline = pipes[0]
								isRunning = true
							}
						} else {
							isRunning = false
						}
					}

					if retry == "View Logs" {
						// ToDo: bad idea to call another sub-command. should be fixed to avoid cyclo imports
						//    and the a shared function placed in the ciutils sub-module
						return ciTraceCmd.TraceRun(&ciTraceCmd.TraceOpts{
							Branch:     branch,
							JobID:      0,
							BaseRepo:   f.BaseRepo,
							HTTPClient: f.HttpClient,
							IO:         f.IO,
						})
					}
				}
				return nil
			}
			redCheck := c.Red("✘")
			fmt.Fprintf(f.IO.StdOut, "%s No pipelines running or available on %s branch\n", redCheck, branch)
			return nil
		},
	}

	pipelineStatusCmd.Flags().BoolP("live", "l", false, "Show status in real-time till pipeline ends")
	pipelineStatusCmd.Flags().BoolP("compact", "c", false, "Show status in compact format")
	pipelineStatusCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is current branch)")

	return pipelineStatusCmd
}
