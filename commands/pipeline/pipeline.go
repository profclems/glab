package pipeline

import (
	"fmt"
	"log"
	"math"
	"os"
	"text/tabwriter"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func displayMultiplePipelines(m []*gitlab.PipelineInfo, repo ...string) {
	var (
		projectID string
		err       error
	)
	if len(repo) > 0 {
		projectID = repo[0]
	} else {
		projectID, err = git.GetRepo()
		if err != nil {
			log.Fatal(err)
		}
	}
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing pipelines %d of %d on %s\n\n", len(m), len(m), projectID)
		for _, pipeline := range m {
			duration := utils.TimeToPrettyTimeAgo(*pipeline.CreatedAt)
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
		fmt.Println("No Pipelines available on " + projectID)
	}
}

func fmtDuration(duration float64) string {
	s := math.Mod(duration, 60)
	m := (duration - s) / 60
	s = math.Round(s)
	return fmt.Sprintf("%02vm %02vs", m, s)
}

// pipelineCmd is merge request command
var pipelineCmd = &cobra.Command{
	Use:     "pipeline <command> [flags]",
	Short:   `Manage pipelines`,
	Long:    ``,
	Aliases: []string{"pipe"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 2 {
			return cmd.Help()
		}
		return nil
	},
}

func init() {
	pipelineCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(pipelineCmd)
}
