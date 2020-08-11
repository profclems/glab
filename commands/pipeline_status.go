package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"math"
	"time"
)

var pipelineStatusCmd = &cobra.Command{
	Use:   "status <command> [flags]",
	Short: `Check the status of a single pipeline`,
	Aliases: []string{"stats"},
	Example: heredoc.Doc(`
	$ glab pipeline status 177883
	$ glab pipeline status --live
	$ glab pipeline status --branch=master   // Get pipeline for master branch
	$ glab pipe status   // Get pipeline for current branch
	`),
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmdErr(cmd, args)
			return
		}
		branch, _ := cmd.Flags().GetString("branch")
		live, _ := cmd.Flags().GetBool("live")
		var err error
		if branch == "" {
			branch, err = git.CurrentBranch()
			if err != nil {
				er(err)
			}
		}
		l := &gitlab.ListProjectPipelinesOptions{
			Ref: gitlab.String(branch),
			OrderBy: gitlab.String("updated_at"),
			Sort: gitlab.String("desc"),
		}
		l.Page = 1
		l.PerPage = 1
		//pid := manip.StringToInt(args[0])
		pipes, err := getPipelines(l)
		if err != nil {
			er(err)
		}
		runningPipeline := pipes[0]
		if len(pipes) == 1 {
			isRunning := true
			retry := false
			writer := uilive.New()

			// start listening for updates and render
			writer.Start()
			for isRunning {
				jobs := getPipelineJobs(runningPipeline.ID)
				for _, job := range jobs {
					duration := fmtDuration(job.Duration)
					var status string
					switch s:=job.Status; s {
					case "failed":
						status = color.Red.Sprint(s)
					case "success":
						status = color.Green.Sprint(s)
					default:
						status = color.Gray.Sprint(s)
					}
					fmt.Fprintf(writer, "(%s) • %s\t\t%s\t\t%s\n", status, duration, job.Stage, job.Name)
				}

				fmt.Fprintf(writer.Newline(), "\n%s\n", runningPipeline.WebURL)
				fmt.Fprintf(writer.Newline(), "SHA: %s\n", runningPipeline.SHA)
				fmt.Fprintf(writer.Newline(), "Pipeline State: %s\n", runningPipeline.Status)
				if runningPipeline.Status == "running" && live {
					pipes, err = getPipelines(l)
					if err != nil {
						er(err)
					}
					runningPipeline = pipes[0]
				} else {
					if runningPipeline.Status == "failed" || runningPipeline.Status == "canceled" {
						prompt := &survey.Confirm{
							Message: "Do you want to retry?",
						}
						survey.AskOne(prompt, &retry)
					}
					if retry {
						retryPipelineJob(runningPipeline.ID)
						pipes, err = getPipelines(l)
						runningPipeline = pipes[0]
					} else {
						isRunning = false
					}
				}
				time.Sleep(time.Millisecond * 0)
			}
		} else {
			er("No pipelines running on " + branch + " branch")
		}
	},
}

func retryPipelineJob(pid int) *gitlab.Pipeline {
	gitlabClient, repo := git.InitGitlabClient()
	pipe, _, err := gitlabClient.Pipelines.RetryPipelineBuild(repo, pid)
	if err != nil {
		er(err)
	}
	return pipe
}

func fmtDuration(duration float64) string {
	s := math.Mod(duration, 60)
	m := (duration-s)/60
	s = math.Round(s)
	return fmt.Sprintf("%02vm %02vs", m, s)
}

func getPipelines(l *gitlab.ListProjectPipelinesOptions) ([]*gitlab.PipelineInfo, error) {
	gitlabClient, repo := git.InitGitlabClient()
	pipes, _, err := gitlabClient.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func getSinglePipeline(pid int) (*gitlab.Pipeline, error) {
	gitlabClient, repo := git.InitGitlabClient()
	pipes, _, err := gitlabClient.Pipelines.GetPipeline(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func init() {
	pipelineStatusCmd.Flags().BoolP("live", "l", false, "Show status in realtime till pipeline ends")
	pipelineStatusCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is current branch)")
	pipelineCmd.AddCommand(pipelineStatusCmd)
}