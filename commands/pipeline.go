package commands

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"text/tabwriter"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
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

func retryPipeline(pid int, rep string) (*gitlab.Pipeline, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipe, _, err := gitlabClient.Pipelines.RetryPipelineBuild(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}
func playPipelineJob(pid int, rep string) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipe, _, err := gitlabClient.Jobs.PlayJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func retryPipelineJob(pid int, rep string) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipe, _, err := gitlabClient.Jobs.RetryJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func cancelPipelineJob(rep string, jobID int) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipe, _, err := gitlabClient.Jobs.CancelJob(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func playOrRetryJobs(repo string, jobID int, status string) (*gitlab.Job, error) {
	switch status {
	case "pending", "running":
		return nil, nil
	case "manual":
		j, err := playPipelineJob(jobID, repo)
		if err != nil {
			return nil, err
		}
		return j, nil
	default:

		j, err := retryPipelineJob(jobID, repo)
		if err != nil {
			return nil, err
		}

		return j, nil
	}
}

func erasePipelineJob(pid int, rep string) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipe, _, err := gitlabClient.Jobs.EraseJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func getPipelineJob(jid int, rep string) (*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	job, _, err := gitlabClient.Jobs.GetJob(repo, jid)
	return job, err
}

func getJobs(rep string, opts *gitlab.ListJobsOptions) ([]gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}

	if opts == nil {
		opts = &gitlab.ListJobsOptions{}
	}
	jobs, _, err := gitlabClient.Jobs.ListProjectJobs(repo, opts)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func fmtDuration(duration float64) string {
	s := math.Mod(duration, 60)
	m := (duration - s) / 60
	s = math.Round(s)
	return fmt.Sprintf("%02vm %02vs", m, s)
}

func getPipelines(l *gitlab.ListProjectPipelinesOptions, rep string) ([]*gitlab.PipelineInfo, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipes, _, err := gitlabClient.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func getPipelineJobs(pid int, rep string) ([]*gitlab.Job, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	l := &gitlab.ListJobsOptions{}
	pipeJobs, _, err := gitlabClient.Jobs.ListPipelineJobs(repo, pid, l)
	if err != nil {
		return nil, err
	}
	return pipeJobs, nil
}

func getPipelineJobLog(jobID int, rep string) (io.Reader, error) {
	gitlabClient, repo := git.InitGitlabClient()
	if rep != "" {
		repo = rep
	}
	pipeJoblog, _, err := gitlabClient.Jobs.GetTraceFile(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipeJoblog, nil
}

func getSinglePipeline(pid int, rep string) (*gitlab.Pipeline, error) {
	gitlabClient, repo := git.InitGitlabClient()

	if rep != "" {
		repo = rep
	}
	pipes, _, err := gitlabClient.Pipelines.GetPipeline(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func getCommit(repository string, ref string) (*gitlab.Commit, error) {
	gitlabClient, repo := git.InitGitlabClient()

	if repository != "" {
		repo = repository
	}
	c, _, err := gitlabClient.Commits.GetCommit(repo, ref)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getPipelineFromBranch(ref, repo string) ([]*gitlab.Job, error) {
	var err error
	if ref == "" {
		ref, err = git.CurrentBranch()
		if err != nil {
			return nil, err
		}
	}
	l := &gitlab.ListProjectPipelinesOptions{
		Ref:     gitlab.String(ref),
		OrderBy: gitlab.String("updated_at"),
		Sort:    gitlab.String("desc"),
	}
	l.Page = 1
	l.PerPage = 1
	if repo == "" {
		repo = git.GetRepo()
	}
	pipes, err := getPipelines(l, repo)
	if err != nil {
		return nil, err
	}
	if len(pipes) == 0 {
		err = errors.New("No pipelines running or available on " + ref + "branch")
		return nil, err
	}
	pipeline := pipes[0]
	jobs, err := getPipelineJobs(pipeline.ID, repo)
	return jobs, nil
}

func pipelineJobTraceWithSha(pid interface{}, sha, name string) (io.Reader, *gitlab.Job, error) {
	gitlabClient, _ := git.InitGitlabClient()
	jobs, err := pipelineJobsWithSha(pid, sha)
	if len(jobs) == 0 || err != nil {
		return nil, nil, err
	}
	var (
		job          *gitlab.Job
		lastRunning  *gitlab.Job
		firstPending *gitlab.Job
	)

	for _, j := range jobs {
		if j.Status == "running" {
			lastRunning = j
		}
		if j.Status == "pending" && firstPending == nil {
			firstPending = j
		}
		if j.Name == name {
			job = j
			// don't break because there may be a newer version of the job
		}
	}
	if job == nil {
		job = lastRunning
	}
	if job == nil {
		job = firstPending
	}
	if job == nil {
		job = jobs[len(jobs)-1]
	}
	r, _, err := gitlabClient.Jobs.GetTraceFile(pid, job.ID)
	if err != nil {
		return nil, job, err
	}

	return r, job, err
}

// CIJobs returns a list of jobs in a pipeline for a given sha. The jobs are
// returned sorted by their CreatedAt time
func pipelineJobsWithSha(pid interface{}, sha string) ([]*gitlab.Job, error) {
	gitlabClient, _ := git.InitGitlabClient()
	pipelines, _, err := gitlabClient.Pipelines.ListProjectPipelines(pid, &gitlab.ListProjectPipelinesOptions{
		SHA: gitlab.String(sha),
	})
	if len(pipelines) == 0 || err != nil {
		return nil, err
	}
	target := pipelines[0].ID
	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 500,
		},
	}
	list, resp, err := gitlabClient.Jobs.ListPipelineJobs(pid, target, opts)
	if err != nil {
		return nil, err
	}
	if resp.CurrentPage == resp.TotalPages {
		return list, nil
	}
	opts.Page = resp.NextPage
	for {
		jobs, resp, err := gitlabClient.Jobs.ListPipelineJobs(pid, target, opts)
		if err != nil {
			return nil, err
		}
		opts.Page = resp.NextPage
		list = append(list, jobs...)
		if resp.CurrentPage == resp.TotalPages {
			break
		}
	}
	return list, nil
}

func pipelineCILint(content string) (*gitlab.LintResult, error) {
	gitlabClient, _ := git.InitGitlabClient()
	c, _, err := gitlabClient.Validate.Lint(content)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// pipelineCmd is merge request command
var pipelineCmd = &cobra.Command{
	Use:     "pipeline <command> [flags]",
	Short:   `Manage pipelines`,
	Long:    ``,
	Aliases: []string{"pipe"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	pipelineCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(pipelineCmd)
}
