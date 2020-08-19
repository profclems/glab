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
	Use:     "for <issue_id>",
	Short:   `Create new merge request for existing issue`,
	Long:    ``,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			cmdErr(cmd, args)
			return nil
		}
		pid := manip.StringToInt(args[0])

		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}

		issue, _, err := gitlabClient.Issues.GetIssue(repo, pid)
		if err != nil {
			return err
		}

		var mergeTitle string
		var mergeDescription string
		var targetBranch string

		if t, _ := cmd.Flags().GetString("target-branch"); t != "" {
			targetBranch = strings.Trim(t, "[] ")
		} else {
			targetBranch = "master"
		}
		if fill, _ := cmd.Flags().GetBool("fill"); !fill {
			if title, _ := cmd.Flags().GetString("title"); title != "" {
				mergeTitle = strings.Trim(title, " ")
			} else {
				mergeTitle = fmt.Sprintf("Resolve: %s (#%d)", issue.Description, issue.IID)
			}
			mergeDescription, _ = cmd.Flags().GetString("description")
		} else {
			branch, _ := git.CurrentBranch()
			commit, _ := git.LatestCommit(branch)
			_, err := getCommit(repo, targetBranch)
			if err != nil {
				return fmt.Errorf("target branch %s does not exist on remote. Specify target branch with --target-branch flag", targetBranch)
			}
			mergeDescription, err = git.CommitBody(branch)
			if err != nil {
				return err
			}
			mergeTitle = strings.Trim(commit.Title, "'")
		}
		if closesIssue, _ := cmd.Flags().GetBool("closes-issue"); closesIssue {
			mergeDescription = fmt.Sprintf("Closes #%d\n%s", issue.IID, mergeDescription)
		}

		newBranch := fmt.Sprintf("%d-%s", issue.IID, manip.ReplaceNonAlphaNumericChars(strings.ToLower(issue.Description), "-"))

		cbo := &gitlab.CreateBranchOptions{}
		cbo.Branch = gitlab.String(newBranch)
		cbo.Ref = gitlab.String(targetBranch)

		fmt.Println("Creating related branch...")
		branch, _, err := gitlabClient.Branches.CreateBranch(repo, cbo)
		if err == nil {
			fmt.Println("Branch created: ", branch.WebURL)
		} else {
			for counter, branchErr := 1, err; branchErr == nil; counter++ {
				newBranch = fmt.Sprintf("%d-%s-%d", issue.IID, manip.ReplaceNonAlphaNumericChars(strings.ToLower(issue.Description), "-"), counter)
				cbo.Branch = gitlab.String(newBranch)
				branch, _, branchErr = gitlabClient.Branches.CreateBranch(repo, cbo)
				if branchErr == nil {
					fmt.Println("Branch created: ", branch.WebURL)
				}
			}
		}

		l := &gitlab.CreateMergeRequestOptions{}

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
		l.Title = gitlab.String(mergeTitle)
		l.Description = gitlab.String(mergeDescription)
		l.Labels = &gitlab.Labels{mergeLabel}
		l.SourceBranch = cbo.Branch
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
	mrForCmd.Flags().BoolP("closes-issue", "", true, "This merge request closes the issue when it is merged or closed")
	mrForCmd.Flags().BoolP("fill", "f", false, "Do not prompt for title/description and just use commit info")
	mrForCmd.Flags().BoolP("draft", "", true, "Mark merge request as a draft")
	mrForCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrForCmd.Flags().StringP("title", "t", "", "Supply a title for merge request")
	mrForCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrForCmd.Flags().StringP("assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrForCmd.Flags().StringP("target-branch", "b", "", "The target or base branch into which you want your code merged (default master)")
	mrForCmd.Flags().IntP("target-project", "", -1, "Add target project by id")
	mrForCmd.Flags().IntP("milestone", "m", -1, "add milestone by <id> for merge request")
	mrForCmd.Flags().BoolP("allow-collaboration", "", false, "Allow commits from other members")
	mrForCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")
	mrForCmd.Flags().BoolP("no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")
	mrCmd.AddCommand(mrForCmd)
}
