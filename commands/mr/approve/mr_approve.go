package approve

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/api"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdApprove(f *cmdutils.Factory) *cobra.Command {
	var mrApproveCmd = &cobra.Command{
		Use:   "approve {<id> | <branch>}",
		Short: `Approve merge requests`,
		Long:  ``,
		Example: heredoc.Doc(`
			$ glab mr approve 235
			$ glab mr approve 123 345
			$ glab mr approve branch-1
			$ glab mr approve branch-2 branch-3
			$ glab mr approve    # Finds open merge request from current branch
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			c := f.IO.Color()

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mrs, repo, err := mrutils.MRsFromArgs(f, args, "opened")
			if err != nil {
				return err
			}

			for _, mr := range mrs {
				if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
					WorkInProgress: true,
					Closed:         true,
					Merged:         true,
				}); err != nil {
					return err
				}

				opts := &gitlab.ApproveMergeRequestOptions{}
				if s, _ := cmd.Flags().GetString("sha"); s != "" {
					opts.SHA = gitlab.String(s)
				}

				fmt.Fprintf(f.IO.StdOut, "- Approving Merge Request !%d\n", mr.IID)
				_, err = api.ApproveMR(apiClient, repo.FullName(), mr.IID, opts)
				if err != nil {
					return err
				}
				fmt.Fprintln(f.IO.StdOut, c.GreenCheck(), "Approved")
			}

			return nil
		},
	}

	//mrApproveCmd.Flags().StringP("password", "p", "", "Current user’s password. Required if 'Require user password to approve' is enabled in the project settings.")
	mrApproveCmd.Flags().StringP("sha", "s", "", "SHA which must match the SHA of the HEAD commit of the merge request")

	carapace.Gen(mrApproveCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			branches := action.ActionBranches(mrApproveCmd, f, &gitlab.ListBranchesOptions{}).Invoke(c)
			mergeRequests := action.ActionMergeRequests(mrApproveCmd, f, &gitlab.ListProjectMergeRequestsOptions{}).Invoke(c)
			return branches.Merge(mergeRequests).Filter(c.Args).ToA()
		}),
	)

	return mrApproveCmd
}
