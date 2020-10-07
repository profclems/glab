package issues

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/spf13/cobra"

	"github.com/xanzy/go-gitlab"
)

func NewCmdIssues(f *cmdutils.Factory) *cobra.Command {
	var mrIssuesCmd = &cobra.Command{
		Use:     "issues <id>",
		Short:   `Get issues related to a particular merge request.`,
		Long:    ``,
		Aliases: []string{"issue"},
		Args:    cobra.ExactArgs(1),
		Example: "$ glab mr issues 46",
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

			mergeID := strings.TrimSpace(args[0])
			l := &gitlab.GetIssuesClosedOnMergeOptions{}

			mrIssues, err := api.GetMRLinkedIssues(apiClient, repo.FullName(), utils.StringToInt(mergeID), l)
			if err != nil {
				return err
			}
			fmt.Fprintln(out, issueutils.DisplayAllIssues(mrIssues, repo.FullName()))
			return nil
		},
	}

	return mrIssuesCmd
}
