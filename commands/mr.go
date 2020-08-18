package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/utils"
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
	if len(m) > 0 {
		table := uitable.New()
		table.MaxColWidth = 70
		fmt.Println()
		fmt.Printf("Showing mergeRequests %d of %d on %s\n\n", len(m), len(m), git.GetRepo())
		for _, mr := range m {
			var mrID string
			if mr.State == "opened" {
				mrID = color.Sprintf("<green>#%d</>", mr.IID)
			} else {
				mrID = color.Sprintf("<red>#%d</>", mr.IID)
			}
			table.AddRow(mrID, mr.Title, color.Sprintf("<cyan>(%s) ← (%s)</>", mr.TargetBranch, mr.SourceBranch))
		}
		fmt.Println(table)
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
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	mrCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	RootCmd.AddCommand(mrCmd)
}
