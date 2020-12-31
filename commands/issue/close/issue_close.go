package close

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdClose(f *cmdutils.Factory) *cobra.Command {
	var issueCloseCmd = &cobra.Command{
		Use:     "close <id>",
		Short:   `Close an issue`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			issueID := strings.TrimSpace(args[0])

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintln(f.IO.StdOut, "- Closing Issue...")
				issue, err := api.UpdateIssue(apiClient, repo.FullName(), utils.StringToInt(i2), l)
				if err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "%s Closed Issue #%s\n", utils.RedCheck(), i2)
				fmt.Fprintln(f.IO.StdOut, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}
	return issueCloseCmd
}
