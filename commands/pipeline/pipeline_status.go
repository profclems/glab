package pipeline

import (
	"fmt"
	"github.com/profclems/glab/internal/utils"
	"time"

	"github.com/profclems/glab/internal/git"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if len(args) != 0 {
			cmdErr(cmd, args)
			return nil
		}
		branch, _ := cmd.Flags().GetString("branch")
		var repo string
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, _ = fixRepoNamespace(r)
		} else {
			repo, err = git.GetRepo()
			if err != nil {
				return err
			}
		}
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
		pipes, err := getPipelines(l, repo)
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
				jobs, err := getPipelineJobs(runningPipeline.ID, repo)
				if err != nil {
					return err
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
							_, err = retryPipeline(runningPipeline.ID, repo)
							if err != nil {
								return err
							}
							pipes, err = getPipelines(l, repo)
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
					//args = []string{strconv.FormatInt(int64(runningPipeline.ID), 10)}
					pipelineCITrace(cmd, args)
					fmt.Println("logs")
				}
			}
		}
		redCheck := utils.Red("✘")
		fmt.Fprintf(colorableErr(cmd), "%s No pipelines running or available on %s branch\n", redCheck, branch)
		return nil
	},
}

func init() {
	pipelineStatusCmd.Flags().BoolP("live", "l", false, "Show status in realtime till pipeline ends")
	pipelineStatusCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is current branch)")
	pipelineCmd.AddCommand(pipelineStatusCmd)
}
