package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"io"
	"math"
	"os"
	"text/tabwriter"

	"github.com/xanzy/go-gitlab"
)

func displayMultiplePipelines(m []*gitlab.PipelineInfo) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing pipelines %d of %d on %s\n\n", len(m), len(m), git.GetRepo())
		for _, pipeline := range m {
			duration := manip.TimeAgo(*pipeline.CreatedAt)
			var pipeState string
			if pipeline.Status == "success" {
				pipeState = color.Sprintf("<green>(%s) • #%d</>", pipeline.Status, pipeline.ID)
			} else if pipeline.Status == "failed" {
				pipeState = color.Sprintf("<red>(%s) • #%d</>", pipeline.Status, pipeline.ID)
			} else {
				pipeState = color.Sprintf("<gray>(%s) • #%d</>", pipeline.Status, pipeline.ID)
			}

			color.Printf("%s\t%s\t<magenta>(%s)</>\n", pipeState, pipeline.Ref, duration)
		}
	} else {
		fmt.Println("No Pipelines available on " + git.GetRepo())
	}
}

func retryPipelineJob(pid int) *gitlab.Pipeline {
	gitlabClient, repo := git.InitGitlabClient()
	pipe, _, err := gitlabClient.Pipelines.RetryPipelineBuild(repo, pid)
	if err != nil {
		er(err)
	}
	return pipe
}

func getPipelineJob(jid int) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	job, _, err := gitlabClient.Jobs.GetJob(repo, jid)
	return job, err
}

func fmtDuration(duration float64) string {
	s := math.Mod(duration, 60)
	m := (duration - s) / 60
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

func getPipelineJobs(pid int) []*gitlab.Job {
	gitlabClient, repo := git.InitGitlabClient()
	l := &gitlab.ListJobsOptions{}
	pipeJobs, _, err := gitlabClient.Jobs.ListPipelineJobs(repo, pid, l)
	if err != nil {
		er(err)
	}
	return pipeJobs
}

func getPipelineJobLog(jobID int) io.Reader {
	gitlabClient, repo := git.InitGitlabClient()
	pipeJobs, _, err := gitlabClient.Jobs.GetTraceFile(repo, jobID)
	if err != nil {
		er(err)
	}
	return pipeJobs
}

func getSinglePipeline(pid int) (*gitlab.Pipeline, error) {
	gitlabClient, repo := git.InitGitlabClient()
	pipes, _, err := gitlabClient.Pipelines.GetPipeline(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

// pipelineCmd is merge request command
var pipelineCmd = &cobra.Command{
	Use:     "pipeline <command> [flags]",
	Short:   `Manage pipelines`,
	Long:    ``,
	Aliases: []string{"pipe"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	pipelineCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(pipelineCmd)
}
