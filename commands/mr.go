package commands

import (
	"fmt"

	"github.com/profclems/glab/internal/utils"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func displayMergeRequest(hm *gitlab.MergeRequest) {
	duration := utils.TimeToPrettyTimeAgo(*hm.CreatedAt)
	if hm.State == "opened" {
		color.Printf("<green>#%d</> %s <magenta>(%s)</> %s\n", hm.IID, hm.Title, hm.SourceBranch, duration)
	} else {
		color.Printf("<red>#%d</> %s <magenta>(%s)</> %s\n", hm.IID, hm.Title, hm.SourceBranch, duration)
	}
	fmt.Println(hm.WebURL)
}

func displayAllMergeRequests(m []*gitlab.MergeRequest) {
	DisplayList(ListInfo{
		Name:    "Merge Requests",
		Columns: []string{"ID", "Title", "Branch"},
		Total:   len(m),
		GetCellValue: func(ri int, ci int) interface{} {
			mr := m[ri]
			switch ci {
			case 0:
				if mr.State == "opened" {
					return color.Sprintf("<green>#%d</>", mr.IID)
				} else {
					return color.Sprintf("<red>#%d</>", mr.IID)
				}
			case 1:
				return mr.Title
			case 2:
				return color.Sprintf("<cyan>(%s) ← (%s)</>", mr.TargetBranch, mr.SourceBranch)
			default:
				return ""
			}
		},
	})
}

// mrCmd is merge request command
var mrCmd = &cobra.Command{
	Use:   "mr <command> [flags]",
	Short: `Create, view and manage merge requests`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	mrCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(mrCmd)
}
