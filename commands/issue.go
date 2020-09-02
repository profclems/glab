package commands

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/utils"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func displayAllIssues(m []*gitlab.Issue) {
	DisplayList(ListInfo{
		Name:    "issues",
		Columns: []string{"IssueID", "Title", "Labels", "CreatedAt"},
		Total:   len(m),
		GetCellValue: func(ri int, ci int) interface{} {
			issue := m[ri]
			switch ci {
			case 0:
				if issue.State == "opened" {
					return color.Sprintf("<green>#%d</>", issue.IID)
				} else {
					return color.Sprintf("<red>#%d</>", issue.IID)
				}
			case 1:
				return issue.Title
			case 2:
				if len(issue.Labels) > 0 {
					return color.Cyan.Sprintf("(%s)", strings.Trim(strings.Join(issue.Labels, ", "), ","))
				}
				return ""
			case 3:
				return color.Gray.Sprintf(utils.TimeToPrettyTimeAgo(*issue.CreatedAt))
			default:
				return ""
			}
		},
	})
}

func displayIssue(hm *gitlab.Issue) {
	duration := utils.TimeToPrettyTimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		color.Printf("<green>#%d</> %s <magenta>(%s)</>\n", hm.IID, hm.Title, duration)
	} else {
		color.Printf("<red>#%d</> %s <magenta>(%s)</>\n", hm.IID, hm.Title, duration)
	}
	fmt.Println(hm.WebURL)
}

// mrCmd is merge request command
var issueCmd = &cobra.Command{
	Use:   "issue [command] [flags]",
	Short: `Create, view and manage remote issues`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	issueCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(issueCmd)
}
