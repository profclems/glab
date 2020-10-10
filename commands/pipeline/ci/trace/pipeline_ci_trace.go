package trace

import (
	"context"
	"fmt"
	"regexp"

	"github.com/profclems/glab/commands/cmdutils"
	ciViewCmd "github.com/profclems/glab/commands/pipeline/ci/view"
	"github.com/profclems/glab/commands/pipeline/pipelineutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdTrace(f *cmdutils.Factory) *cobra.Command {
	var pipelineCITraceCmd = &cobra.Command{
		Use:   "trace <job-id> [flags]",
		Short: `Work with GitLab CI pipelines and jobs`,
		Example: heredoc.Doc(`
	$ glab pipeline ci trace
	`),
		Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return TraceCmdFunc(cmd, args, f)
		},
	}

	pipelineCITraceCmd.Flags().StringP("branch", "b", "", "Check pipeline status for a branch. (Default is the current branch)")
	return pipelineCITraceCmd
}

func TraceCmdFunc(cmd *cobra.Command, args []string, f *cmdutils.Factory) error {
	var jobID int
	var err error

	out := utils.ColorableOut(cmd)

	apiClient, err := f.HttpClient()
	if err != nil {
		return err
	}

	repo, err := f.BaseRepo()
	if err != nil {
		return err
	}

	branch, _ := cmd.Flags().GetString("branch")
	if branch == "" {
		branch, err = git.CurrentBranch()
		if err != nil {
			return err
		}
	}

	if len(args) != 0 {
		jobID = utils.StringToInt(args[0])
	} else {
		l := &gitlab.ListProjectPipelinesOptions{
			Ref:     gitlab.String(branch),
			OrderBy: gitlab.String("updated_at"),
			Sort:    gitlab.String("desc"),
		}

		l.Page = 1
		l.PerPage = 1

		fmt.Fprintf(out, "Searching for latest pipeline on %s...\n", branch)

		pipes, err := api.GetPipelines(apiClient, l, repo.FullName())
		if err != nil {
			return err
		}

		if len(pipes) == 0 {
			fmt.Fprintln(out, "No pipeline running or available on "+branch+"branch")
			return nil
		}

		pipeline := pipes[0]
		fmt.Fprintf(out, "Getting jobs for pipeline %d...\n", pipeline.ID)

		jobs, err := api.GetPipelineJobs(apiClient, pipeline.ID, repo.FullName())
		if err != nil {
			return err
		}

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
			jobID = utils.StringToInt(m[0][1])
		} else {
			jobID = jobs[0].ID
		}
	}

	commit, err := api.GetCommit(apiClient, repo.FullName(), branch)
	if err != nil {
		return err
	}

	ciViewCmd.CommitSHA = commit.ID

	job, err := api.GetPipelineJob(apiClient, jobID, repo.FullName())
	if err != nil {
		return err
	}

	err = pipelineutils.RunTrace(apiClient, context.Background(), out, repo.FullName(), job.Pipeline.Sha, job.Name)
	if err != nil {
		return err
	}

	return nil
}
