package commands

import (
	"fmt"
	"github.com/gookit/color"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
)

func displayAllIssues(m []*gitlab.Issue) {
	if len(m) > 0 {
		fmt.Printf("\nShowing issues %d of %d on %s\n\n", len(m), len(m), git.GetRepo())

		// initialize tabwriter
		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		defer w.Flush()
		for _, issue := range m {
			var labels string
			for _, l := range issue.Labels {
				labels += " " + l + ","
			}
			labels = strings.Trim(labels, ", ")
			duration := manip.TimeAgo(*issue.CreatedAt)
			if issue.State == "opened" {
				_, _ = fmt.Fprintln(w, color.Sprintf("<green>#%d</>\t%s\t<cyan>(%s)</>\t<gray>%s</>", issue.IID, issue.Title, labels, duration))
			} else {
				_, _ = fmt.Fprintln(w, color.Sprintf("<red>#%d</>\t%s\t<cyan>(%s)</>\t<gray>%s</>", issue.IID, issue.Title, labels, duration))
			}
		}
	} else {
		fmt.Println("No Issues available on " + git.GetRepo())
	}
}

func displayIssue(hm *gitlab.Issue) {
	duration := manip.TimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		color.Printf("<green>#%d</> %s <magenta>(%s)</>\n", hm.IID, hm.Title, duration)
	} else {
		color.Printf("<red>#%d</> %s <magenta>(%s)</>\n", hm.IID, hm.Title, duration)
	}
	fmt.Println(hm.WebURL)
}

// mrCmd is merge request command
var issueCmd = &cobra.Command{
	Use:   "issue <command> [flags]",
	Short: `Create, view and manage remote issues`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(issueCmd)
}
