package commands

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	//"strings"
	"sync"
	"time"
	//"github.com/pkg/errors"
)

var (
	once   sync.Once
	offset int64
)

var pipelineCITraceCmd = &cobra.Command{
	Use:   "trace <job-id> [flags]",
	Short: `Work with GitLab CI pipelines and jobs`,
	Example: heredoc.Doc(`
	$ glab pipeline ci trace
	`),
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			cmdErr(cmd, args)
		}
		var jobID int
		if len(args) != 0 {
			jobID = manip.StringToInt(args[0])
		} else {
			branch, _ := cmd.Flags().GetString("branch")
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
			pipes, err := getPipelines(l)
			if err != nil {
				er(err)
			}
			if len(pipes) == 0 {
				er("No pipelines running or available on " + branch + "branch")
			}
			pipeline := pipes[0]
			jobs := getPipelineJobs(pipeline.ID)
			var jobOptions []string
			var selectedJob string
			for _, job := range jobs {
				jobOptions = append(jobOptions, fmt.Sprintf("%s (%d) - %s", job.Name, job.ID, job.Status))
			}
			prompt := &survey.Select{
				Message: "Select pipeline job to trace:",
				Options: jobOptions,
			}
			_ = survey.AskOne(prompt, &selectedJob)
			if selectedJob != "" {
				re := regexp.MustCompile(`(?s)\((.*)\)`)
				m := re.FindAllStringSubmatch(selectedJob, -1)
				jobID = manip.StringToInt(m[0][1])
			}

			if err := runTrace(jobID); err != nil {
				er(err)
			}
		}
	},
}

func runTrace(jobID int) error {
	w := os.Stdout
	ctx := context.Background()
	for range time.NewTicker(time.Second * 3).C {
		if ctx.Err() == context.Canceled {
			break
		}
		job, err := getPipelineJob(jobID)
		trace := getPipelineJobLog(jobID)

		if err != nil || job == nil || trace == nil {
			return errors.Wrap(err, "failed to find job")
		}

		switch job.Status {
		case "pending":
			_, _ = fmt.Fprintf(w, "%s is pending... waiting for job to start...\n", job.Name)
			return nil
		case "manual":
			_, _ = fmt.Fprintf(w, "Manual job %s not started, waiting for job to start\n", job.Name)
			continue
		case "skipped":
			_, _ = fmt.Fprintf(w, "%s has been skipped\n", job.Name)
			continue
		}
		once.Do(func() {
			_, _ = fmt.Fprintf(w, "Showing logs for %s job #%d\n", job.Name, job.ID)
		})
		_, err = io.CopyN(ioutil.Discard, trace, offset)
		lenT, err := io.Copy(w, trace)
		if err != nil {
			return err
		}
		offset += lenT

		if job.Status == "success" ||
			job.Status == "failed" ||
			job.Status == "cancelled" {
			return nil
		}
	}
	return nil
}

func init() {
	pipelineCITraceCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is the current branch)")
	pipelineCICmd.AddCommand(pipelineCITraceCmd)
}
