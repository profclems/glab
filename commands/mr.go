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

func displayMergeRequest(hm *gitlab.MergeRequest) {
	duration := manip.TimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		color.Printf("<green>#%d</> %s <magenta>(%s)</> %s\n", hm.IID, hm.Title, hm.SourceBranch, duration)
	} else {
		color.Printf("<red>#%d</> %s <magenta>(%s)</> %s\n", hm.IID, hm.Title, hm.SourceBranch, duration)
	}
	fmt.Println(hm.WebURL)
}

func displayAllMergeRequests(m []*gitlab.MergeRequest) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Println()
		fmt.Printf("Showing mergeRequests %d of %d on %s\n", len(m), len(m), git.GetRepo())
		for _, mr := range m {
			if mr.State == "opened" {
				_, _ = fmt.Fprintln(w, color.Sprintf("<green>#%d</>\t%s\t\t<cyan>(%s) ← (%s)</>", mr.IID, mr.Title, mr.TargetBranch, mr.SourceBranch))
			} else {
				_, _ = fmt.Fprintln(w, color.Sprintf("<green>#%d</>\t%s\t\t<cyan>(%s) ← (%s)</>", mr.IID, mr.Title, mr.TargetBranch, mr.SourceBranch))
			}
		}
		fmt.Println()
	} else {
		fmt.Println("No Merge Requests available on " + git.GetRepo())
	}
}

// mrCmd is merge request command
var mrCmd = &cobra.Command{
	Use:   "mr <command> [flags]",
	Short: `Create, view and manage merge requests`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(mrCmd)
}