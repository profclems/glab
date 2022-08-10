package ciutils

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/profclems/glab/pkg/iostreams"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/pkg/tableprinter"
	"github.com/profclems/glab/pkg/utils"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

var (
	once   sync.Once
	offset int64
)

func makeHyperlink(s *iostreams.IOStreams, pipeline *gitlab.PipelineInfo) string {
	return s.Hyperlink(fmt.Sprintf("%d", pipeline.ID), pipeline.WebURL)
}

func DisplayMultiplePipelines(s *iostreams.IOStreams, p []*gitlab.PipelineInfo, projectID string) string {
	c := s.Color()

	table := tableprinter.NewTablePrinter()

	if len(p) > 0 {

		for _, pipeline := range p {
			duration := ""

			if pipeline.CreatedAt != nil {
				duration = c.Magenta("(" + utils.TimeToPrettyTimeAgo(*pipeline.CreatedAt) + ")")
			}

			var pipeState string
			if pipeline.Status == "success" {
				pipeState = c.Green(fmt.Sprintf("(%s) • #%s", pipeline.Status, makeHyperlink(s, pipeline)))
			} else if pipeline.Status == "failed" {
				pipeState = c.Red(fmt.Sprintf("(%s) • #%s", pipeline.Status, makeHyperlink(s, pipeline)))
			} else {
				pipeState = c.Gray(fmt.Sprintf("(%s) • #%s", pipeline.Status, makeHyperlink(s, pipeline)))
			}

			table.AddRow(pipeState, pipeline.Ref, duration)
		}

		return table.Render()
	}

	return "No Pipelines available on " + projectID
}

func DisplayMultipleTriggers(s *iostreams.IOStreams, triggers []*gitlab.PipelineTrigger) string {
	c := s.Color()

	table := tableprinter.NewTablePrinter()

	if len(triggers) > 0 {

		for _, trigger := range triggers {
			duration := ""
			name := "unknown owner"
			if trigger.CreatedAt != nil {
				duration = c.Magenta("(" + utils.TimeToPrettyTimeAgo(*trigger.CreatedAt) + ")")
			}

			if trigger.Owner != nil {
				name = trigger.Owner.Name
			}

			line := c.Green(fmt.Sprintf("[%s] • #%s", name, trigger.ID))

			table.AddRow(line, trigger.ID, duration)
		}

		return table.Render()
	}

	return "" // empty message is provided by the title
}

func RunTraceSha(ctx context.Context, apiClient *gitlab.Client, w io.Writer, pid interface{}, sha, name string) error {
	job, err := api.PipelineJobWithSha(apiClient, pid, sha, name)
	if err != nil || job == nil {
		return errors.Wrap(err, "failed to find job")
	}
	return RunTrace(ctx, apiClient, w, pid, job, name)
}

func RunTrace(ctx context.Context, apiClient *gitlab.Client, w io.Writer, pid interface{}, job *gitlab.Job, name string) error {
	fmt.Fprintln(w, "Getting job trace...")
	for range time.NewTicker(time.Second * 3).C {
		if ctx.Err() == context.Canceled {
			break
		}
		trace, _, err := apiClient.Jobs.GetTraceFile(pid, job.ID)
		if err != nil || trace == nil {
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
