package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/manip"
	"strings"
)

var mrForCmd = &cobra.Command{
	Use:     "for",
	Short:   `Create new merge request for an issue`,
	Long:    ``,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || len(args) > 1 {
			cmdErr(cmd, args)
			return nil
		}

		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}

		issueID := manip.StringToInt(args[0])
		issue, _, err := gitlabClient.Issues.GetIssue(repo, issueID)
		if err != nil {
			return err
		}

		sourceBranch := fmt.Sprintf("%d-%s", issue.IID, manip.ReplaceNonAlphaNumericChars(strings.ToLower(issue.Title), "-"))

		lb := &gitlab.CreateBranchOptions{
			Branch: gitlab.String(sourceBranch),
			Ref:    gitlab.String("master"),
		}

		_, _, err = gitlabClient.Branches.CreateBranch(repo, lb)
		if err != nil {
			for branchErr, branchCount := err, 1; branchErr != nil; branchCount++ {

				numberedBranch := fmt.Sprintf("%d-%s-%d", issue.IID, strings.ReplaceAll(strings.ToLower(issue.Title), " ", "-"), branchCount)
				lb = &gitlab.CreateBranchOptions{
					Branch: gitlab.String(numberedBranch),
					Ref:    gitlab.String("master"),
				}
				sourceBranch = numberedBranch
				_, _, branchErr = gitlabClient.Branches.CreateBranch(repo, lb)
				fmt.Println(branchErr)
			}

		}

		var targetBranch string
		if t, _ := cmd.Flags().GetString("target-branch"); t != "" {
			targetBranch = strings.Trim(t, "[] ")
		} else {
			targetBranch = "master"
		}

		var mergeTitle string
		mergeTitle = fmt.Sprintf("%s (#%d)", issue.Title, issue.IID)

		isDraft, _ := cmd.Flags().GetBool("draft")
		isWIP, _ := cmd.Flags().GetBool("wip")
		if isDraft || isWIP {
			if isDraft {
				mergeTitle = "Draft: " + mergeTitle
			} else {
				mergeTitle = "WIP: " + mergeTitle
			}
		}

		mergeLabel, _ := cmd.Flags().GetString("label")

		l := &gitlab.CreateMergeRequestOptions{}
		l.Title = gitlab.String(mergeTitle)
		l.Description = gitlab.String(fmt.Sprintf("Closes #%d", issue.IID))
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

		if a, _ := cmd.Flags().GetString("assignee"); a != "" {
			arrIds := strings.Split(strings.Trim(a, "[] "), ",")
			var t2 []int

			for _, i := range arrIds {
				j := manip.StringToInt(i)
				t2 = append(t2, j)
			}
			l.AssigneeIDs = t2
		}

		mr, _, err := gitlabClient.MergeRequests.CreateMergeRequest(repo, l)
		if err != nil {
			return err
		}
		displayMergeRequest(mr)

		return nil
	},
}

func init() {
	mrForCmd.Flags().BoolP("draft", "", true, "Mark merge request as a draft")
	mrForCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrForCmd.Flags().StringP("title", "t", "", "Supply a title for merge request")
	mrForCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrForCmd.Flags().StringP("assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrForCmd.Flags().BoolP("allow-collaboration", "", false, "Allow commits from other members")
	mrForCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")
	mrForCmd.Flags().BoolP("closes-issue", "", true, "Issue is closed when the PR is merged or closed")
	mrForCmd.Flags().BoolP("create-source-branch", "", true, "Create source branch if it does not exist")
	mrCmd.AddCommand(mrForCmd)
}
