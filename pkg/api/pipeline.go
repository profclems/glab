package api

import (
	"errors"
	"io"

	"github.com/profclems/glab/internal/git"
	"github.com/xanzy/go-gitlab"
)

var RetryPipeline = func(client *gitlab.Client, pid int, repo string) (*gitlab.Pipeline, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Pipelines.RetryPipelineBuild(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

var PlayPipelineJob = func(client *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Jobs.PlayJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

var RetryPipelineJob = func(client *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Jobs.RetryJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

var CancelPipelineJob = func(client *gitlab.Client, repo string, jobID int) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Jobs.CancelJob(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

var PlayOrRetryJobs = func(client *gitlab.Client, repo string, jobID int, status string) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	switch status {
	case "pending", "running":
		return nil, nil
	case "manual":
		j, err := PlayPipelineJob(client, jobID, repo)
		if err != nil {
			return nil, err
		}
		return j, nil
	default:

		j, err := RetryPipelineJob(client, jobID, repo)
		if err != nil {
			return nil, err
		}

		return j, nil
	}
}

var ErasePipelineJob = func(client *gitlab.Client, pid int, repo string) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Jobs.EraseJob(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

var GetPipelineJob = func(client *gitlab.Client, jid int, repo string) (*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	job, _, err := client.Jobs.GetJob(repo, jid)
	return job, err
}

var GetJobs = func(client *gitlab.Client, repo string, opts *gitlab.ListJobsOptions) ([]*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}

	if opts == nil {
		opts = &gitlab.ListJobsOptions{}
	}
	jobs, _, err := client.Jobs.ListProjectJobs(repo, opts)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var GetPipelines = func(client *gitlab.Client, l *gitlab.ListProjectPipelinesOptions, repo interface{}) ([]*gitlab.PipelineInfo, error) {
	if client == nil {
		client = apiClient
	}
	if l.PerPage == 0 {
		l.PerPage = DefaultListLimit
	}

	pipes, _, err := client.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

var GetPipelineJobs = func(client *gitlab.Client, pid int, repo string) ([]*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	l := &gitlab.ListJobsOptions{}
	pipeJobs, _, err := client.Jobs.ListPipelineJobs(repo, pid, l)
	if err != nil {
		return nil, err
	}
	return pipeJobs, nil
}

var GetPipelineJobLog = func(client *gitlab.Client, jobID int, repo string) (io.Reader, error) {
	if client == nil {
		client = apiClient
	}
	pipeJoblog, _, err := client.Jobs.GetTraceFile(repo, jobID)
	if err != nil {
		return nil, err
	}
	return pipeJoblog, nil
}

var GetSinglePipeline = func(client *gitlab.Client, pid int, repo string) (*gitlab.Pipeline, error) {
	if client == nil {
		client = apiClient
	}
	pipes, _, err := client.Pipelines.GetPipeline(repo, pid)
	if err != nil {
		return nil, err
	}
	return pipes, nil
}

var GetCommit = func(client *gitlab.Client, repo string, ref string) (*gitlab.Commit, error) {
	if client == nil {
		client = apiClient
	}
	c, _, err := client.Commits.GetCommit(repo, ref)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var GetPipelineFromBranch = func(client *gitlab.Client, ref, repo string) ([]*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	var err error
	if ref == "" {
		ref, err = git.CurrentBranch()
		if err != nil {
			return nil, err
		}
	}
	l := &gitlab.ListProjectPipelinesOptions{
		Ref:  gitlab.String(ref),
		Sort: gitlab.String("desc"),
	}
	l.Page = 1
	l.PerPage = 1
	pipes, err := GetPipelines(client, l, repo)
	if err != nil {
		return nil, err
	}
	if len(pipes) == 0 {
		err = errors.New("No pipelines running or available on " + ref + "branch")
		return nil, err
	}
	pipeline := pipes[0]
	jobs, err := GetPipelineJobs(client, pipeline.ID, repo)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var PipelineJobTraceWithSha = func(client *gitlab.Client, pid interface{}, sha, name string) (io.Reader, *gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	jobs, err := PipelineJobsWithSha(client, pid, sha)
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
	r, _, err := client.Jobs.GetTraceFile(pid, job.ID)
	if err != nil {
		return nil, job, err
	}

	return r, job, err
}

// PipelineJobsWithSha returns a list of jobs in a pipeline for a given sha. The jobs are
// returned sorted by their CreatedAt time
var PipelineJobsWithSha = func(client *gitlab.Client, pid interface{}, sha string) ([]*gitlab.Job, error) {
	if client == nil {
		client = apiClient
	}
	pipelines, err := GetPipelines(client, &gitlab.ListProjectPipelinesOptions{
		SHA: gitlab.String(sha),
	}, pid)
	if len(pipelines) == 0 || err != nil {
		return nil, err
	}
	target := pipelines[0].ID
	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 500,
		},
	}
	list, resp, err := client.Jobs.ListPipelineJobs(pid, target, opts)
	if err != nil {
		return nil, err
	}
	if resp.CurrentPage == resp.TotalPages {
		return list, nil
	}
	opts.Page = resp.NextPage
	for {
		jobs, resp, err := client.Jobs.ListPipelineJobs(pid, target, opts)
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

var PipelineCILint = func(client *gitlab.Client, content string) (*gitlab.LintResult, error) {
	if client == nil {
		client = apiClient
	}
	c, _, err := client.Validate.Lint(content)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var DeletePipeline = func(client *gitlab.Client, projectID interface{}, pipeID int) error {
	if client == nil {
		client = apiClient
	}
	_, err := client.Pipelines.DeletePipeline(projectID, pipeID)
	if err != nil {
		return err
	}
	return nil
}

var ListProjectPipelines = func(client *gitlab.Client, projectID interface{}, opts *gitlab.ListProjectPipelinesOptions) ([]*gitlab.PipelineInfo, error) {
	if client == nil {
		client = apiClient
	}
	pipes, _, err := client.Pipelines.ListProjectPipelines(projectID, opts)
	if err != nil {
		return pipes, err
	}
	return pipes, nil
}

var CreatePipeline = func(client *gitlab.Client, projectID interface{}, opts *gitlab.CreatePipelineOptions) (*gitlab.Pipeline, error) {
	if client == nil {
		client = apiClient
	}
	pipe, _, err := client.Pipelines.CreatePipeline(projectID, opts)
	return pipe, err
}
