package tally

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

// tally tracks how many times a job or pipeline has been run, and
// some statistics of how long the job took to complete.
type tally struct {
	Count int
	Max   time.Duration
	Min   time.Duration
	Total time.Duration
}

func (t *tally) Add(duration time.Duration) {
	t.Count++
	t.Total += duration
	if duration > t.Max {
		t.Max = duration
	}
	if t.Min == 0 || duration < t.Min {
		t.Min = duration
	}
}

// statusTally tracks a tally for each job/pipeline status.
type statusTally map[string]*tally // key is status, value is tallied statistics

// Count formats the tally count as a string.
func (st statusTally) Count(status string) string {
	t := st[status]
	if t == nil {
		return "0"
	}
	return strconv.Itoa(t.Count)
}

// Avg calculates an average duration, and formats as a string (number of seconds).
func (st statusTally) Avg(status string) string {
	t := st[status]
	if t == nil {
		return ""
	}
	secs := int(t.Total/time.Second) / t.Count
	return strconv.Itoa(secs)
}

// Max formats the longest duration as a string (number of seconds).
func (st statusTally) Max(status string) string {
	t := st[status]
	if t == nil {
		return ""
	}
	secs := int(t.Max / time.Second)
	return strconv.Itoa(secs)
}

// Min formats the shortest duration as a string (number of seconds).
func (st statusTally) Min(status string) string {
	t := st[status]
	if t == nil {
		return ""
	}
	secs := int(t.Min / time.Second)
	return strconv.Itoa(secs)
}

var (
	statusList = []string{"success", "failed", "running", "pending", "canceled", "skipped", "created", "manual"}
)

func NewCmdTally(f *cmdutils.Factory) *cobra.Command {
	var tallyCmd = &cobra.Command{
		Use:   "tally [flags]",
		Short: `Calculate statistics of CI pipelines and jobs`,
		Example: heredoc.Doc(`
	$ glab ci tally
	$ glab ci tally --branch=master
	`),
		Long: ``,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			opt := &gitlab.ListProjectPipelinesOptions{}

			if m, _ := cmd.Flags().GetString("status"); m != "" {
				opt.Status = gitlab.BuildState(gitlab.BuildStateValue(m))
			}
			if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
				opt.PerPage = p
			}

			branch, _ := cmd.Flags().GetString("branch")
			if branch != "" {
				opt.Ref = &branch
			}

			doPipelines, err := cmd.Flags().GetBool("pipeline")
			if err != nil {
				return err
			}
			doJobs, err := cmd.Flags().GetBool("job")
			if err != nil {
				return err
			}

			pipes, err := api.ListProjectPipelines(apiClient, repo.FullName(), opt)
			if err != nil {
				return err
			}

			type tallyKey struct {
				Ref string // branch
				Job string // if "", tally is for a pipeline
			}

			branchStat := map[tallyKey]statusTally{}

			addStats := func(key tallyKey, status string, start, end *time.Time) {
				// initialize tally
				if branchStat[key] == nil {
					branchStat[key] = statusTally{}
				}
				if branchStat[key][status] == nil {
					branchStat[key][status] = &tally{}
				}
				if end != nil && start != nil {
					branchStat[key][status].Add(end.Sub(*start))
				} else {
					branchStat[key][status].Add(0) // increase count but not duration
				}
			}

			for i := range pipes {
				if doPipelines {
					key := tallyKey{Ref: pipes[i].Ref}
					status := pipes[i].Status
					addStats(key, status, pipes[i].CreatedAt, pipes[i].UpdatedAt)
				}

				if doJobs {
					job, err := api.GetPipelineJobs(apiClient, pipes[i].ID, repo.FullName())
					if err != nil {
						return err
					}

					for j := range job {
						key := tallyKey{
							Ref: pipes[i].Ref,
							Job: job[j].Name,
						}
						status := job[j].Status
						if job[j].FinishedAt != nil {
							addStats(key, status, job[j].StartedAt, job[j].FinishedAt)
						}
					} // end each job in pipeline
				} // end tally job

			} // end loop each pipeline

			// prepare to write comma-separated values
			w := csv.NewWriter(os.Stdout)
			head := []string{
				"branch",
				"job",
			}
			for _, status := range statusList {
				head = append(head,
					status, // how many with this status
					fmt.Sprintf("%s avg (sec)", status),
					fmt.Sprintf("%s min (sec)", status),
					fmt.Sprintf("%s max (sec)", status),
				)
			}
			w.Write(head)
			for key, stat := range branchStat {
				col := []string{
					key.Ref,
					key.Job,
				}
				for _, status := range statusList {
					col = append(col,
						stat.Count(status),
						stat.Avg(status),
						stat.Min(status),
						stat.Max(status),
					)
				}
				w.Write(col)
			}

			w.Flush()

			if w.Error() != nil {
				return w.Error()
			}

			return nil
		},
	}

	tallyCmd.Flags().StringP("status", "s", "", fmt.Sprintf("Tally pipelines with status: {%s}", strings.Join(statusList, "|")))
	tallyCmd.Flags().IntP("per-page", "P", 30, "Number of recent pipelines to tally.")

	tallyCmd.Flags().StringP("branch", "b", "", "Limit tally to pipelines on a particular branch.")
	tallyCmd.Flags().BoolP("job", "j", true, "Tally statistics for jobs.")
	tallyCmd.Flags().BoolP("pipeline", "", true, "Tally statistics for pipelines.")

	return tallyCmd
}
