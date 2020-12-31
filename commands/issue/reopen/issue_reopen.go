package reopen

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

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var issueReopenCmd = &cobra.Command{
		Use:     "reopen <id>",
		Short:   `Reopen a closed issue`,
		Long:    ``,
		Aliases: []string{"open"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := f.IO.StdOut

			gLabClient, err := f.HttpClient()

			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			issueID := strings.TrimSpace(args[0])

			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("reopen")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				if f.IO.IsaTTY && f.IO.IsErrTTY {
					fmt.Fprintln(out, "- Reopening Issue...")
				}

				issue, err := api.UpdateIssue(gLabClient, repo.FullName(), utils.StringToInt(i2), l)
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s Reopened Issue #%s\n", utils.GreenCheck(), i2)
				fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}

	return issueReopenCmd
}
