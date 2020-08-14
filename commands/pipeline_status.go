package commands

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmdErr(cmd, args)
			return
		}
		branch, _ := cmd.Flags().GetString("branch")
		var repo string
		rep, _ := cmd.Flags().GetString("repo")
		if rep != "" {
			repo = rep
		}
		live, _ := cmd.Flags().GetBool("live")
		var err error
		if branch == "" {
			branch, err = git.CurrentBranch()
			if err != nil {
				er(err)
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
		pipes, err := getPipelines(l, repo)
		if err != nil {
			er(err)
		}
		if len(pipes) == 1 {
			runningPipeline := pipes[0]
			isRunning := true
			retry := "Exit"
			writer := uilive.New()

			// start listening for updates and render
			writer.Start()
			for isRunning {
				jobs, err := getPipelineJobs(runningPipeline.ID, repo)
				if err != nil {
					er(err)
					return
				}
				for _, job := range jobs {
					duration := fmtDuration(job.Duration)
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
					pipes, err = getPipelines(l, repo)
					if err != nil {
						er(err)
						return
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
							_, err = retryPipeline(runningPipeline.ID, repo)
							if err != nil {
								er(err)
							}
							pipes, err = getPipelines(l, repo)
							if err != nil {
								er(err)
								return
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
					//args = []string{strconv.FormatInt(int64(runningPipeline.ID), 10)}
					pipelineCITrace(cmd, args)
					fmt.Println("logs")
				}
			}
		} else {
			er("No pipelines running or available on " + branch + " branch")
		}
	},
}

func init() {
	pipelineStatusCmd.Flags().BoolP("live", "l", false, "Show status in realtime till pipeline ends")
	pipelineStatusCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is current branch)")
	pipelineCmd.AddCommand(pipelineStatusCmd)
}
