package pipelineutils

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/profclems/glab/pkg/api"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/juju/ansiterm/tabwriter"
	"github.com/profclems/glab/internal/utils"
	"github.com/xanzy/go-gitlab"
)

func DisplayMultiplePipelines(p []*gitlab.PipelineInfo, projectID string) string {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(p) > 0 {
		pipelinePrint := fmt.Sprintf("Showing pipelines %d of %d on %s\n\n", len(p), len(p), projectID)

		for _, pipeline := range p {
			duration := utils.TimeToPrettyTimeAgo(*pipeline.CreatedAt)
			var pipeState string
			if pipeline.Status == "success" {
				pipeState = utils.Green(fmt.Sprintf("(%s) • #%d", pipeline.Status, pipeline.ID))
			} else if pipeline.Status == "failed" {
				pipeState = utils.Red(fmt.Sprintf("(%s) • #%d", pipeline.Status, pipeline.ID))
			} else {
				pipeState = utils.Gray(fmt.Sprintf("(%s) • #%d", pipeline.Status, pipeline.ID))
			}

			pipelinePrint += fmt.Sprintf("%s\t%s\t%s\n", pipeState, pipeline.Ref, utils.Magenta("("+duration+")"))
		}
	}

	return "No Pipelines available on " + projectID
}

func RunTrace(apiClient *gitlab.Client, ctx context.Context, w io.Writer, pid interface{}, sha, name string) error {
	var (
		once   sync.Once
		offset int64
	)
	fmt.Fprintln(w, "Getting job trace...")
	for range time.NewTicker(time.Second * 3).C {
		if ctx.Err() == context.Canceled {
			break
		}
		trace, job, err := api.PipelineJobTraceWithSha(apiClient, pid, sha, name)
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
