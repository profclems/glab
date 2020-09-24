package gitlab

import (
	"errors"
	"io"

	"github.com/profclems/glab/internal/git"

	"github.com/xanzy/go-gitlab"
)

func RetryPipeline(gLab *gitlab.Client, pid int, repo string) (*gitlab.Pipeline, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipe, _, err := gLab.Pipelines.RetryPipelineBuild(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func PlayPipelineJob(gLab *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipe, _, err := gLab.Jobs.PlayJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func RetryPipelineJob(gLab *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipe, _, err := gLab.Jobs.RetryJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func CancelPipelineJob(gLab *gitlab.Client, repo string, jobID int) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipe, _, err := gLab.Jobs.CancelJob(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func PlayOrRetryJobs(gLab *gitlab.Client, repo string, jobID int, status string) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	switch status {
	case "pending", "running":
		return nil, nil
	case "manual":
		j, err := PlayPipelineJob(gLab, jobID, repo)
		if err != nil {
			return nil, err
		}
		return j, nil
	default:

		j, err := RetryPipelineJob(gLab, jobID, repo)
		if err != nil {
			return nil, err
		}

		return j, nil
	}
}

func ErasePipelineJob(gLab *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipe, _, err := gLab.Jobs.EraseJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

func GetPipelineJob(gLab *gitlab.Client, jid int, repo string) (*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	job, _, err := gLab.Jobs.GetJob(repo, jid)
	return job, err
}

func GetJobs(gLab *gitlab.Client, repo string, opts *gitlab.ListJobsOptions) ([]gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}

	if opts == nil {
		opts = &gitlab.ListJobsOptions{}
	}
	jobs, _, err := gLab.Jobs.ListProjectJobs(repo, opts)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func GetPipelines(gLab *gitlab.Client, l *gitlab.ListProjectPipelinesOptions, repo string) ([]*gitlab.PipelineInfo, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipes, _, err := gLab.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func GetPipelineJobs(gLab *gitlab.Client, pid int, repo string) ([]*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	l := &gitlab.ListJobsOptions{}
	pipeJobs, _, err := gLab.Jobs.ListPipelineJobs(repo, pid, l)
	if err != nil {
		return nil, err
	}
	return pipeJobs, nil
}

func GetPipelineJobLog(gLab *gitlab.Client, jobID int, repo string) (io.Reader, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipeJoblog, _, err := gLab.Jobs.GetTraceFile(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipeJoblog, nil
}

func GetSinglePipeline(gLab *gitlab.Client, pid int, repo string) (*gitlab.Pipeline, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipes, _, err := gLab.Pipelines.GetPipeline(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

func GetCommit(gLab *gitlab.Client, repo string, ref string) (*gitlab.Commit, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	c, _, err := gLab.Commits.GetCommit(repo, ref)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetPipelineFromBranch(gLab *gitlab.Client, ref, repo string) ([]*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
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
	pipes, err := GetPipelines(gLab, l, repo)
	if err != nil {
		return nil, err
	}
	if len(pipes) == 0 {
		err = errors.New("No pipelines running or available on " + ref + "branch")
		return nil, err
	}
	pipeline := pipes[0]
	jobs, err := GetPipelineJobs(gLab, pipeline.ID, repo)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func PipelineJobTraceWithSha(gLab *gitlab.Client, pid interface{}, sha, name string) (io.Reader, *gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	jobs, err := PipelineJobsWithSha(gLab, pid, sha)
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
	r, _, err := gLab.Jobs.GetTraceFile(pid, job.ID)
	if err != nil {
		return nil, job, err
	}

	return r, job, err
}

// PipelineJobsWithSha returns a list of jobs in a pipeline for a given sha. The jobs are
// returned sorted by their CreatedAt time
func PipelineJobsWithSha(gLab *gitlab.Client, pid interface{}, sha string) ([]*gitlab.Job, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	pipelines, _, err := gLab.Pipelines.ListProjectPipelines(pid, &gitlab.ListProjectPipelinesOptions{
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
	list, resp, err := gLab.Jobs.ListPipelineJobs(pid, target, opts)
	if err != nil {
		return nil, err
	}
	if resp.CurrentPage == resp.TotalPages {
		return list, nil
	}
	opts.Page = resp.NextPage
	for {
		jobs, resp, err := gLab.Jobs.ListPipelineJobs(pid, target, opts)
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

func PipelineCILint(gLab *gitlab.Client, content string) (*gitlab.LintResult, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	c, _, err := gLab.Validate.Lint(content)
	if err != nil {
		return nil, err
	}
	return c, nil
}
