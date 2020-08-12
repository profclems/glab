package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"log"
	"strings"
)

var mrCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   `Create new merge request`,
	Long:    ``,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmdErr(cmd, args)
			return
		}
		l := &gitlab.CreateMergeRequestOptions{}
		var sourceBranch string
		var targetBranch string
		var mergeTitle string
		var mergeLabel string
		var mergeDescription string
		if title, _ := cmd.Flags().GetString("title"); title != "" {
			mergeTitle = strings.Trim(title, " ")
		} else {
			mergeTitle = manip.AskQuestionWithInput("Title:", "", true)
		}
		if label, _ := cmd.Flags().GetString("label"); label != "" {
			mergeLabel = strings.Trim(label, "[] ")
		} else {
			mergeLabel = manip.AskQuestionWithInput("Label(s) [Comma Separated]:", "", false)
		}
		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			mergeDescription = strings.Trim(desc, " ")
		} else {
			mergeDescription = manip.AskQuestionMultiline("Description:", "")
		}
		if source, _ := cmd.Flags().GetString("source"); source != "" {
			sourceBranch = strings.Trim(source, "[] ")
		} else {
			if c, _ := cmd.Flags().GetBool("create-source-branch"); c {
				sourceBranch = manip.ReplaceNonAlphaNumericChars(mergeTitle, "-")
			} else {
				sourceBranch = manip.AskQuestionWithInput("Source Branch:", "", true)
			}
		}
		if t, _ := cmd.Flags().GetString("target"); t != "" {
			targetBranch = strings.Trim(t, "[] ")
		} else {
			targetBranch = manip.AskQuestionWithInput("Target Branch:", "", true)
		}
		l.Title = gitlab.String(mergeTitle)
		l.Description = gitlab.String(mergeDescription)
		l.Labels = &gitlab.Labels{mergeLabel}
		l.SourceBranch = gitlab.String(sourceBranch)
		l.TargetBranch = gitlab.String(targetBranch)
		if milestone, _ := cmd.Flags().GetInt("milestone"); milestone != -1 {
			l.MilestoneID = gitlab.Int(milestone)
		}
		if allowCol, _ := cmd.Flags().GetBool("allow-collaboration"); allowCol {
			l.AllowCollaboration = gitlab.Bool(true)
		}
		if removeSource, _ := cmd.Flags().GetBool("remove-source-branch"); removeSource {
			l.RemoveSourceBranch = gitlab.Bool(true)
		}
		if targetProject, _ := cmd.Flags().GetInt("target-project"); targetProject != -1 {
			l.TargetProjectID = gitlab.Int(targetProject)
		}
		if a, _ := cmd.Flags().GetString("assigns"); a != "" {
			arrIds := strings.Split(strings.Trim(a, "[] "), ",")
			var t2 []int

			for _, i := range arrIds {
				j := manip.StringToInt(i)
				t2 = append(t2, j)
			}
			l.AssigneeIDs = t2
		}

		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		if c, _ := cmd.Flags().GetBool("create-source-branch"); c {
			lb := &gitlab.CreateBranchOptions{
				Branch: gitlab.String(sourceBranch),
				Ref:    gitlab.String(targetBranch),
			}
			fmt.Println("Creating related branch...")
			branch, resp, _ := gitlabClient.Branches.CreateBranch(repo, lb)
			if resp.StatusCode == 201 {
				fmt.Println("Branch created: ", branch.WebURL)
			} else {
				fmt.Println("Error creating branch: ", resp.Status)
			}
		}

		mr, _, err := gitlabClient.MergeRequests.CreateMergeRequest(repo, l)
		if err != nil {
			log.Fatal(err)
		}
		displayMergeRequest(mr)
	},
}

func init() {
	mrCreateCmd.Flags().StringP("title", "t", "", "Supply a title for merge request")
	mrCreateCmd.Flags().StringP("description", "d", "", "Supply a description for merge request")
	mrCreateCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrCreateCmd.Flags().StringP("assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrCreateCmd.Flags().StringP("source-branch", "s", "", "Source Branch for merge request")
	mrCreateCmd.Flags().StringP("target-branch", "g", "", "Target Branch for merge request")
	mrCreateCmd.Flags().IntP("target-project", "", -1, "Add target project by id")
	mrCreateCmd.Flags().BoolP("create-source-branch", "", false, "Create source branch if it does not exist")
	mrCreateCmd.Flags().IntP("milestone", "m", -1, "add milestone by <id> for merge request")
	mrCreateCmd.Flags().BoolP("allow-collaboration", "", false, "Allow commits from other members")
	mrCreateCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")
	mrCmd.AddCommand(mrCreateCmd)
}
