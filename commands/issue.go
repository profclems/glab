package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

func displayAllIssues(m []*gitlab.Issue) {
	if len(m) > 0 {
		fmt.Printf("\nShowing issues %d of %d on %s\n\n", len(m), len(m), git.GetRepo())
		table := uitable.New()
		table.MaxColWidth = 70
		for _, issue := range m {
			var labels string
			for _, l := range issue.Labels {
				labels += " " + l + ","
			}
			labels = strings.Trim(labels, ", ")
			if labels != ""{
				labels = "(" + labels + ")"
			}
			var issueID string
			duration := manip.TimeAgo(*issue.CreatedAt)
			if issue.State == "opened" {
				issueID = color.Sprintf("<green>#%d</>", issue.IID)
			} else {
				issueID = color.Sprintf("<red>#%d</>", issue.IID)
			}
			table.AddRow(issueID, issue.Title, color.Cyan.Sprintf(labels), color.Gray.Sprintf(duration))
		}
		fmt.Println(table)
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
