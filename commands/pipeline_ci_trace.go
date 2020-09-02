package commands

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
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
	Run:  pipelineCITrace,
}

func pipelineCITrace(cmd *cobra.Command, args []string) {
	if len(args) > 1 {
		cmdErr(cmd, args)
	}
	var jobID int
	var repo string
	repo, _ = cmd.Flags().GetString("repo")
	if repo == "" {
		repo = git.GetRepo()
	}
	branch, _ := cmd.Flags().GetString("branch")
	var err error
	if branch == "" {
		branch, err = git.CurrentBranch()
		if err != nil {
			er(err)
		}
	}

	if len(args) != 0 {
		jobID = manip.StringToInt(args[0])
	} else {
		l := &gitlab.ListProjectPipelinesOptions{
			Ref:     gitlab.String(branch),
			OrderBy: gitlab.String("updated_at"),
			Sort:    gitlab.String("desc"),
		}
		l.Page = 1
		l.PerPage = 1
		fmt.Printf("Searching for latest pipeline on %s...\n", branch)
		pipes, err := getPipelines(l, repo)
		if err != nil {
			er(err)
		}
		if len(pipes) == 0 {
			er("No pipelines running or available on " + branch + "branch")
		}
		pipeline := pipes[0]
		fmt.Printf("Getting jobs for pipeline %d...\n", pipeline.ID)
		jobs, err := getPipelineJobs(pipeline.ID, repo)
		if err != nil {
			er(err)
			return
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
			jobID = manip.StringToInt(m[0][1])
		} else {
			jobID = jobs[0].ID
		}
	}

	commit, err := getCommit(repo, branch)
	if err != nil {
		log.Fatal(err)
	}
	commitSHA = commit.ID
	job, err := getPipelineJob(jobID, repo)
	if err != nil {
		er(err)
		return
	}
	err = runTrace(context.Background(), os.Stdout, repo, job.Pipeline.Sha, job.Name)
	if err != nil {
		log.Fatal(err)
	}
}

func runTrace(ctx context.Context, w io.Writer, pid interface{}, sha, name string) error {
	var (
		once   sync.Once
		offset int64
	)
	fmt.Fprintln(w, "Getting job trace...")
	for range time.NewTicker(time.Second * 3).C {
		if ctx.Err() == context.Canceled {
			break
		}
		trace, job, err := pipelineJobTraceWithSha(pid, sha, name)
		if err != nil || job == nil || trace == nil {
			return errors.Wrap(err, "failed to find job")
		}
		switch job.Status {
		case "pending":
			fmt.Fprintf(w, "%s is pending... waiting for job to start\n", job.Name)
			continue
		case "manual":
			fmt.Fprintf(w, "Manual job %s not started, waiting for job to start\n", job.Name)
			continue
		case "skipped":
			fmt.Fprintf(w, "%s has been skipped\n", job.Name)
			break
		}
		once.Do(func() {
			if name == "" {
				name = job.Name
			}
			fmt.Fprintf(w, "Showing logs for %s job #%d\n", job.Name, job.ID)
		})
		_, _ = io.CopyN(ioutil.Discard, trace, offset)
		lenT, err := io.Copy(w, trace)
		if err != nil {
			return err
		}
		offset += int64(lenT)

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
