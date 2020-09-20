package commands

import (
	"fmt"

	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
)

var mrCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   `Create new merge request`,
	Long:    ``,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cmdErr(cmd, args)
			return nil
		}
		var sourceBranch string
		var mergeTitle string
		var mergeDescription string
		var err error
		var targetBranch string
		l := &gitlab.CreateMergeRequestOptions{}

		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo, err = fixRepoNamespace(r)
			if err != nil {
				return err
			}
		}
		if t, _ := cmd.Flags().GetString("target-branch"); t != "" {
			targetBranch = t
		} else {
			targetBranch, _ = git.GetDefaultBranch(repo)
		}
		if source, _ := cmd.Flags().GetString("source-branch"); source != "" {
			sourceBranch = strings.Trim(source, "[] ")
		} else {
			if c, _ := cmd.Flags().GetBool("create-source-branch"); c && sourceBranch == "" {
				sourceBranch = manip.ReplaceNonAlphaNumericChars(mergeTitle, "-")
			} else {
				b, err := git.CurrentBranch()
				if err != nil {
					return err
				}
				sourceBranch = b
			}
		}
		if fill, _ := cmd.Flags().GetBool("fill"); !fill {
			if title, _ := cmd.Flags().GetString("title"); title != "" {
				mergeTitle = strings.Trim(title, " ")
			} else {
				mergeTitle = manip.AskQuestionWithInput("Title:", "", true)
			}
			if desc, _ := cmd.Flags().GetString("description"); desc != "" {
				mergeDescription = desc
			} else {
				if editor, _ := cmd.Flags().GetBool("no-editor"); editor {
					mergeDescription = manip.AskQuestionMultiline("Description:", "")
				} else {
					mergeDescription = manip.Editor(manip.EditorOptions{
						Label:    "Description:",
						Help:     "Enter the MR description. ",
						FileName: "*_MR_EDITMSG.md",
					})
				}
			}
		} else {
			branch, _ := git.CurrentBranch()
			commit, _ := git.LatestCommit(branch)
			_, err := getCommit(repo, targetBranch)
			if err != nil {
				return fmt.Errorf("target branch %s does not exist on remote. Specify target branch with --target-branch flag",
					targetBranch)
			}
			mergeDescription, err = git.CommitBody(branch)
			if err != nil {
				return err
			}
			mergeTitle = utils.Humanize(commit.Title)
			if c, err := git.UncommittedChangeCount(); c != 0 {
				if err != nil {
					return err
				}
				fmt.Printf("warning: you have %v uncommitted changes\n", c)
			}
			remoteURL, err := git.GetRemoteURL()
			if err != nil {
				return err
			}
			err = git.Push(remoteURL, sourceBranch)
			if err != nil {
				return err
			}
		}
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
			return err
		}
		displayMergeRequest(mr)
		return nil
	},
}

func init() {
	mrCreateCmd.Flags().BoolP("fill", "f", false, "Do not prompt for title/description and just use commit info")
	mrCreateCmd.Flags().BoolP("draft", "", false, "Mark merge request as a draft")
	mrCreateCmd.Flags().BoolP("wip", "", false, "Mark merge request as a work in progress. Alternative to --draft")
	mrCreateCmd.Flags().BoolP("push", "", false, "Push committed changes after creating merge request. Make sure you have committed changes")
	mrCreateCmd.Flags().StringP("title", "t", "", "Supply a title for merge request")
	mrCreateCmd.Flags().StringP("description", "d", "", "Supply a description for merge request")
	mrCreateCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	mrCreateCmd.Flags().StringP("assignee", "a", "", "Assign merge request to people by their IDs. Multiple values should be comma separated ")
	mrCreateCmd.Flags().StringP("source-branch", "s", "", "The Branch you are creating the merge request. Default is the current branch.")
	mrCreateCmd.Flags().StringP("target-branch", "b", "", "The target or base branch into which you want your code merged")
	mrCreateCmd.Flags().IntP("target-project", "", -1, "Add target project by id")
	mrCreateCmd.Flags().BoolP("create-source-branch", "", false, "Create source branch if it does not exist")
	mrCreateCmd.Flags().IntP("milestone", "m", -1, "add milestone by <id> for merge request")
	mrCreateCmd.Flags().BoolP("allow-collaboration", "", false, "Allow commits from other members")
	mrCreateCmd.Flags().BoolP("remove-source-branch", "", false, "Remove Source Branch on merge")
	mrCreateCmd.Flags().BoolP("no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")
	mrCmd.AddCommand(mrCreateCmd)
}
