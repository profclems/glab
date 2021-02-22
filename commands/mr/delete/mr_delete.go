package delete

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"

	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutils.Factory) *cobra.Command {
	var mrDeleteCmd = &cobra.Command{
		Use:     "delete [<id> | <branch>]",
		Short:   `Delete merge requests`,
		Long:    ``,
		Aliases: []string{"del"},
		Example: heredoc.Doc(`
			$ glab mr delete 123
			$ glab mr delete 123 branch-name 789  # close multiple branches
			$ glab mr delete 1,2,branch-related-to-mr-3,4,5  # close MRs !1,!2,!3,!4,!5
			$ glab mr del 123
			$ glab mr delete branch
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c := f.IO.Color()
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args)
			if err != nil {
				return err
			}

			for _, mr := range mrs {
				fmt.Fprintf(f.IO.StdOut, "- Deleting Merge Request !%d\n", mr.IID)
				if err = api.DeleteMR(apiClient, repo.FullName(), mr.IID); err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "%s Merge request !%d deleted\n", c.RedCheck(), mr.IID)
			}

			return nil
		},
	}

	return mrDeleteCmd
}
