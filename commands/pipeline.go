package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"glab/internal/git"
	"glab/internal/manip"
	"os"
	"text/tabwriter"

	"github.com/xanzy/go-gitlab"
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

// pipelineCmd is merge request command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline <command> [flags]",
	Short: `Manage pipelines`,
	Long:  ``,
	Aliases: []string{"pipe"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(pipelineCmd)
}
