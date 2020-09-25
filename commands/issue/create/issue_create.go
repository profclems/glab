package create

import (
	"fmt"
	"github.com/profclems/glab/pkg/api"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/manip"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdCreate(f *cmdutils.Factory) *cobra.Command {
	var issueCreateCmd = &cobra.Command{
		Use:     "create [flags]",
		Short:   `Create an issue`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			l := &gitlab.CreateIssueOptions{}
			var (
				issueTitle       string
				issueLabel       string
				issueDescription string
				err              error
			)

			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}
			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			if title, _ := cmd.Flags().GetString("title"); title != "" {
				issueTitle = strings.Trim(title, " ")
			} else {
				issueTitle = manip.AskQuestionWithInput("Title", "", true)
			}
			if description, _ := cmd.Flags().GetString("description"); description != "" {
				issueDescription = strings.Trim(description, " ")
			} else {
				if editor, _ := cmd.Flags().GetBool("no-editor"); editor {
					issueDescription = manip.AskQuestionMultiline("Description:", "")
				} else {
					issueDescription = manip.Editor(manip.EditorOptions{
						Label:    "Description:",
						Help:     "Enter the issue description. ",
						FileName: "*_ISSUE_EDITMSG.md",
					})
				}
			}
			if label, _ := cmd.Flags().GetString("label"); label != "" {
				issueLabel = strings.Trim(label, "[] ")
			} else {
				labelsEntry := config.GetEnv("PROJECT_LABELS")
				if labelsEntry != "" {
					labels := strings.Split(labelsEntry, ",")
					issueLabel = strings.Join(manip.AskQuestionWithMultiSelect("Label(s)", labels), ",")
				} else {
					issueLabel = manip.AskQuestionWithInput("Label(s) [Comma Separated]", "", false)
				}
			}
			l.Title = gitlab.String(issueTitle)
			l.Labels = gitlab.Labels{issueLabel}
			l.Description = &issueDescription
			if confidential, _ := cmd.Flags().GetBool("confidential"); confidential {
				l.Confidential = gitlab.Bool(confidential)
			}
			if weight, _ := cmd.Flags().GetInt("weight"); weight != -1 {
				l.Weight = gitlab.Int(weight)
			}
			if a, _ := cmd.Flags().GetInt("linked-mr"); a != -1 {
				l.MergeRequestToResolveDiscussionsOf = gitlab.Int(a)
			}
			if a, _ := cmd.Flags().GetInt("milestone"); a != -1 {
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
			issue, err := api.CreateIssue(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}
			fmt.Fprintln(utils.ColorableOut(cmd), issueutils.DisplayIssue(issue))
			return nil
		},
	}
	issueCreateCmd.Flags().StringP("title", "t", "", "Supply a title for issue")
	issueCreateCmd.Flags().StringP("description", "d", "", "Supply a description for issue")
	issueCreateCmd.Flags().StringP("label", "l", "", "Add label by name. Multiple labels should be comma separated")
	issueCreateCmd.Flags().StringP("assignee", "a", "", "Assign issue to people by their ID. Multiple values should be comma separated ")
	issueCreateCmd.Flags().IntP("milestone", "m", -1, "The global ID of a milestone to assign issue")
	issueCreateCmd.Flags().BoolP("confidential", "c", false, "Set an issue to be confidential. Default is false")
	issueCreateCmd.Flags().IntP("linked-mr", "", -1, "The IID of a merge request in which to resolve all issues")
	issueCreateCmd.Flags().IntP("weight", "w", -1, "The weight of the issue. Valid values are greater than or equal to 0.")
	issueCreateCmd.Flags().BoolP("no-editor", "", false, "Don't open editor to enter description. If set to true, uses prompt. Default is false")

	return issueCreateCmd
}
