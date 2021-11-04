package _for

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/git"

	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdFor(f *cmdutils.Factory) *cobra.Command {
	var mrForCmd = &cobra.Command{
		Use:     "for",
		Short:   `Create new merge request for an issue`,
		Long:    ``,
		Aliases: []string{"new-for", "create-for", "for-issue"},
		Example: heredoc.Doc(`
	$ glab mr for 34   # Create mr for issue 34
	$ glab mr for 34 --wip   # Create mr and mark as work in progress
	$ glab mr new-for 34
	$ glab mr create-for 34
	`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			issueID := utils.StringToInt(args[0])
			issue, err := api.GetIssue(apiClient, repo.FullName(), issueID)
			if err != nil {
				return err
			}

			remotes, err := f.Remotes()
			if err != nil {
				return err
			}
			repoRemote, err := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())
			if err != nil {
				return err
			}

			var targetBranch string
			if t, _ := cmd.Flags().GetString("target-branch"); t != "" {
				targetBranch = strings.TrimSpace(t)
			} else {
				targetBranch, _ = git.GetDefaultBranch(repoRemote.Name)
			}

			sourceBranch := fmt.Sprintf("%d-%s", issue.IID, utils.ReplaceNonAlphaNumericChars(strings.ToLower(issue.Title), "-"))

			lb := &gitlab.CreateBranchOptions{
				Branch: gitlab.String(sourceBranch),
				Ref:    gitlab.String(targetBranch),
			}

			_, err = api.CreateBranch(apiClient, repo.FullName(), lb)
			if err != nil {
				for branchErr, branchCount := err, 1; branchErr != nil; branchCount++ {

					numberedBranch := fmt.Sprintf("%d-%s-%d", issue.IID, strings.ReplaceAll(strings.ToLower(issue.Title), " ", "-"), branchCount)
					lb = &gitlab.CreateBranchOptions{
						Branch: gitlab.String(numberedBranch),
						Ref:    gitlab.String(targetBranch),
					}
					sourceBranch = numberedBranch
					_, branchErr = api.CreateBranch(apiClient, repo.FullName(), lb)
					fmt.Println(branchErr)
				}

			}

			var mergeTitle string
			mergeTitle = fmt.Sprintf("Resolve \"%s\"", issue.Title)

			isDraft, _ := cmd.Flags().GetBool("draft")
			isWIP, _ := cmd.Flags().GetBool("wip")
			if isDraft || isWIP {
				if isWIP {
					mergeTitle = "WIP: " + mergeTitle
				} else {
					mergeTitle = "Draft: " + mergeTitle
				}
			}

			mergeLabel, _ := cmd.Flags().GetString("label")

			l := &gitlab.CreateMergeRequestOptions{}
			l.Title = gitlab.String(mergeTitle)
			l.Description = gitlab.String(fmt.Sprintf("Closes #%d", issue.IID))
			l.Labels = gitlab.Labels{mergeLabel}
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
			if withLables, _ := cmd.Flags().GetBool("with-labels"); withLables {
				l.Labels = issue.Labels
			}

			if a, _ := cmd.Flags().GetString("assignee"); a != "" {
				arrIds := strings.Split(strings.Trim(a, "[] "), ",")
				var t2 []int

				for _, i := range arrIds {
					j := utils.StringToInt(i)
					t2 = append(t2, j)
				}
				l.AssigneeIDs = t2
			}

			mr, err := api.CreateMR(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(f.IO.Color(), mr))

			return nil
		},
	}

	mrForCmd.Flags().BoolP("draft", "", true, "Mark merge request as a draft. Default is true")
	mrForCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Overrides --draft")
	mrForCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrForCmd.Flags().StringP("assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrForCmd.Flags().BoolP("allow-collaboration", "", false, "Allow commits from other members")
	mrForCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")
	mrForCmd.Flags().IntP("milestone", "m", -1, "add milestone by <id> for merge request")
	mrForCmd.Flags().StringP("target-branch", "b", "", "The target or base branch into which you want your code merged")
	mrForCmd.Flags().BoolP("with-labels", "", false, "Copy labels from issue to the merge request")

	mrForCmd.Deprecated = "use `glab mr create --related-issue <issueID>`"

	return mrForCmd
}
