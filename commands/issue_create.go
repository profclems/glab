package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
)

var issueCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Short:   `Create an issue`,
	Long:    ``,
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmdErr(cmd, args)
			return
		}

		l := &gitlab.CreateIssueOptions{}
		var issueTitle string
		var issueLabel string
		var issueDescription string
		if title, _ := cmd.Flags().GetString("title"); title != "" {
			issueTitle = strings.Trim(title, " ")
		} else {
			issueTitle = manip.AskQuestionWithInput("Title", "", true)
		}
		if label, _ := cmd.Flags().GetString("label"); label != "" {
			issueLabel = strings.Trim(label, "[] ")
		} else {
			issueLabel = manip.AskQuestionWithInput("Label(s) [Comma Separated]", "", false)
		}
		if description, _ := cmd.Flags().GetString("description"); description != "" {
			issueDescription = strings.Trim(description, " ")
		} else {
			issueDescription = manip.AskQuestionMultiline("Description", "")
		}
		//issueDate := manip.AskQuestionWithInput("Due Date (Format: YYYY-MM-DD):", "", false)
		l.Title = gitlab.String(issueTitle)
		l.Labels = &gitlab.Labels{issueLabel}
		l.Description = &issueDescription
		//l.DueDate = &gitlab.ISOTime{issueDate}
		if confidential, _ := cmd.Flags().GetBool("confidential"); confidential {
			l.Confidential = gitlab.Bool(confidential)
		}
		if weight, _ := cmd.Flags().GetInt("weight"); weight != 0 {
			l.Weight = gitlab.Int(weight)
		}
		if a, _ := cmd.Flags().GetInt("linked-merge-request"); a != 0 {
			l.MergeRequestToResolveDiscussionsOf = gitlab.Int(a)
		}
		if a, _ := cmd.Flags().GetInt("milestone"); a != 0 {
			l.MilestoneID = gitlab.Int(a)
		}
		if a, _ := cmd.Flags().GetString("assignee"); a != "" {
			assignID := a
			arrIds := strings.Split(strings.Trim(assignID, "[] "), ",")
			var t2 []int

			for _, i := range arrIds {
				j := manip.StringToInt(i)
				t2 = append(t2, j)
			}
			l.AssigneeIDs = t2
		}
		gitlabClient, repo := git.InitGitlabClient()
		issue, _, err := gitlabClient.Issues.CreateIssue(repo, l)
		if err != nil {
			log.Fatal(err)
		}
		displayIssue(issue)
	},
}

func init() {
	issueCreateCmd.Flags().StringP("title", "t", "", "Supply a title for issue")
	issueCreateCmd.Flags().StringP("description", "d", "", "Supply a description for issue")
	issueCreateCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	issueCreateCmd.Flags().StringP("assignee", "a", "", "Assign issue to people by their ID. Multiple values should be comma separated ")
	issueCreateCmd.Flags().StringP("milestone", "m", "", "add milestone by <id> for issue")
	issueCreateCmd.Flags().BoolP("allow-collaboration", "", false, "Allow collaboration")
	issueCreateCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch after merge")
	issueCmd.AddCommand(issueCreateCmd)
}
