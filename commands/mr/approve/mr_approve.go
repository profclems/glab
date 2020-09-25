package approve

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/manip"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdApprove(f *cmdutils.Factory) *cobra.Command {
	var mrApproveCmd = &cobra.Command{
		Use:   "approve <id> [flags]",
		Short: `Approve merge requests`,
		Long:  ``,
		Args:  cobra.ExactArgs(1),
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

			mergeID := strings.Trim(args[0], " ")
			l := &gitlab.ApproveMergeRequestOptions{}
			if s, _ := cmd.Flags().GetString("sha"); s != "" {
				l.SHA = gitlab.String(s)
			}
			//if s, _ := cmd.Flags().GetString("password"); s  {
			// ToDo:
			//}

			fmt.Fprintf(out, "- Approving Merge Request #%s\n", mergeID)
			_, err = api.ApproveMR(apiClient, repo.FullName(), manip.StringToInt(mergeID), l)
			if err != nil {
				return err
			}
			fmt.Fprintln(out, utils.GreenCheck(), "Approved successfully")

			return nil
		},
	}

	//mrApproveCmd.Flags().StringP("password", "p", "", "Current user’s password. Required if 'Require user password to approve' is enabled in the project settings.")
	mrApproveCmd.Flags().StringP("sha", "s", "", "The HEAD of the merge request")
	return mrApproveCmd
}
