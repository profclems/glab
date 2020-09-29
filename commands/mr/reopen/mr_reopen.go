package reopen

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdReopen(f *cmdutils.Factory) *cobra.Command {
	var mrReopenCmd = &cobra.Command{
		Use:     "reopen <id>",
		Short:   `Reopen merge requests`,
		Long:    ``,
		Aliases: []string{"open"},
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

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("reopen")
			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")

			for _, i2 := range arrIds {
				fmt.Fprintf(out, "- Reopening Merge request !%s...\n", i2)
				mr, err := api.UpdateMR(apiClient, repo.FullName(), utils.StringToInt(i2), l)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%s Merge request !%s reopened\n", utils.GreenCheck(), i2)
				fmt.Fprintln(out, mrutils.DisplayMR(mr))
			}

			return nil
		},
	}

	return mrReopenCmd
}
