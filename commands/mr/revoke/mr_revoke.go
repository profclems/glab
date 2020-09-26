package revoke

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdRevoke(f *cmdutils.Factory) *cobra.Command {
	var mrRevokeCmd = &cobra.Command{
		Use:     "revoke <id>",
		Short:   `Revoke approval on a merge request <id>`,
		Long:    ``,
		Aliases: []string{"unapprove"},
		Args:    cobra.ExactArgs(1),
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

			fmt.Fprintln(out,"- Revoking approval for Merge Request #" + mergeID + "...")

			err = api.UnapproveMR(apiClient, repo.FullName(), utils.StringToInt(mergeID))
			if err != nil {
				return err
			}

			fmt.Fprintln(out, utils.GreenCheck(), "Merge Request approval revoked")

			return nil
		},
	}

	return mrRevokeCmd
}