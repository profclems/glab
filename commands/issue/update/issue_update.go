package update

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdUpdate(f *cmdutils.Factory) *cobra.Command {
	var issueUpdateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: `Update issue`,
		Long:  ``,
		Example: heredoc.Doc(`
	$ glab issue update 42 --label ui,ux
	$ glab issue update 42 --unlabel working
	`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)
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

			issueID := utils.StringToInt(args[0])
			l := &gitlab.UpdateIssueOptions{}

			if m, _ := cmd.Flags().GetString("title"); m != "" {
				l.Title = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetBool("lock-discussion"); m {
				l.DiscussionLocked = gitlab.Bool(m)
			}
			if m, _ := cmd.Flags().GetString("description"); m != "" {
				l.Description = gitlab.String(m)
			}
			if m, _ := cmd.Flags().GetStringArray("label"); len(m) != 0 {
				l.AddLabels = gitlab.Labels(m)
			}
			if m, _ := cmd.Flags().GetStringArray("unlabel"); len(m) != 0 {
				l.RemoveLabels = gitlab.Labels(m)
			}
			fmt.Fprintf(out,"- Updating issue #%d", issueID)
			issue, err := api.UpdateIssue(apiClient, repo.FullName(), issueID, l)
			if err != nil {
				return err
			}

			fmt.Fprintln(out, utils.GreenCheck(), "Updated")

			fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			return nil
		},
	}

	issueUpdateCmd.Flags().StringP("title", "t", "", "Title of issue")
	issueUpdateCmd.Flags().BoolP("lock-discussion", "", false, "Lock discussion on issue")
	issueUpdateCmd.Flags().StringP("description", "d", "", "Issue description")
	issueUpdateCmd.Flags().StringArrayP("label", "l", []string{}, "add labels")
	issueUpdateCmd.Flags().StringArrayP("unlabel", "u", []string{}, "remove labels")

	return issueUpdateCmd
}
