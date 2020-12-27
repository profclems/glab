package ciutils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

var (
	once   sync.Once
	offset int64
)

func DisplayMultiplePipelines(p []*gitlab.PipelineInfo, projectID string) string {

	table := tableprinter.NewTablePrinter()

	if len(p) > 0 {

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

			table.AddRow(pipeState, pipeline.Ref, utils.Magenta("("+duration+")"))
		}

		return table.Render()
	}

	return "No Pipelines available on " + projectID
}

func RunTrace(ctx context.Context, apiClient *gitlab.Client, w io.Writer, pid interface{}, sha, name string) error {
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
		offset += lenT

		if job.Status == "success" ||
			job.Status == "failed" ||
			job.Status == "cancelled" {
			return nil
		}
	}
	return nil
}
