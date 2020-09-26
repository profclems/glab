package status

import (
	"fmt"
	"time"

	"github.com/profclems/glab/commands/cmdutils"
	ciTraceCmd "github.com/profclems/glab/commands/pipeline/ci/trace"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdStatus(f *cmdutils.Factory) *cobra.Command {
	var pipelineStatusCmd = &cobra.Command{
		Use:     "status [flags]",
		Short:   `View a running pipeline on current or other branch specified`,
		Aliases: []string{"stats"},
		Example: heredoc.Doc(`
	$ glab pipeline status --live
	$ glab pipeline status --branch=master   // Get pipeline for master branch
	$ glab pipe status   // Get pipeline for current branch
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			out := utils.ColorableOut(cmd)
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}
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

			if branch == "" {
				branch, err = git.CurrentBranch()
				if err != nil {
					return err
				}
			}
			l := &gitlab.ListProjectPipelinesOptions{
				Ref:     gitlab.String(branch),
				OrderBy: gitlab.String("updated_at"),
				Sort:    gitlab.String("desc"),
			}
			l.Page = 1
			l.PerPage = 1
			//pid := manip.StringToInt(args[0])
			pipes, err := api.GetPipelines(apiClient, l, repo.FullName())
			if err != nil {
				return err
			}
			if len(pipes) == 1 {
				runningPipeline := pipes[0]
				isRunning := true
				retry := "Exit"
				writer := uilive.New()

				// start listening for updates and render
				writer.Start()
				for isRunning {
					jobs, err := api.GetPipelineJobs(apiClient, runningPipeline.ID, repo.FullName())
					if err != nil {
						return err
					}
					for _, job := range jobs {
						duration := utils.PrettyTimeAgo(time.Duration(job.Duration))
						var status string
						switch s := job.Status; s {
						case "failed":
							status = color.Red.Sprint(s)
						case "success":
							status = color.Green.Sprint(s)
						default:
							status = color.Gray.Sprint(s)
						}
						//fmt.Println(job.Tag)
						_, _ = fmt.Fprintf(writer, "(%s) • %s\t\t%s\t\t%s\n", status, duration, job.Stage, job.Name)
					}

					_, _ = fmt.Fprintf(writer.Newline(), "\n%s\n", runningPipeline.WebURL)
					_, _ = fmt.Fprintf(writer.Newline(), "SHA: %s\n", runningPipeline.SHA)
					_, _ = fmt.Fprintf(writer.Newline(), "Pipeline State: %s\n", runningPipeline.Status)
					if runningPipeline.Status == "running" && live {
						pipes, err = api.GetPipelines(apiClient, l, repo.FullName())
						if err != nil {
							return err
						}
						runningPipeline = pipes[0]
					} else {
						if runningPipeline.Status == "failed" || runningPipeline.Status == "canceled" {
							prompt := &survey.Select{
								Message: "Choose an action:",
								Options: []string{"View Logs", "Retry", "Exit"},
								Default: "Exit",
							}
							_ = survey.AskOne(prompt, &retry)
						}
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
					time.Sleep(time.Millisecond * 0)
					if retry == "View Logs" {
						// ToDo: bad idea to call another sub-command. should be fixed to avoid cyclo imports
						//    and the a shared function placed in the pipeutils sub-module
						return ciTraceCmd.TraceCmdFunc(cmd, args, f)
					}
				}
			}
			redCheck := utils.Red("✘")
			fmt.Fprintf(out, "%s No pipelines running or available on %s branch\n", redCheck, branch)
			return nil
		},
	}

	pipelineStatusCmd.Flags().BoolP("live", "l", false, "Show status in real-time till pipeline ends")
	pipelineStatusCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is current branch)")

	return pipelineStatusCmd
}
