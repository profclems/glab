package close

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

func NewCmdClose(f *cmdutils.Factory) *cobra.Command {
	var mrCloseCmd = &cobra.Command{
		Use:   "close <id>",
		Short: `Close merge requests`,
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

			mergeID := args[0]

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintf(out, "- Closing Merge request...")
				mr, err := api.UpdateMR(apiClient, repo.FullName(), utils.StringToInt(i2), l)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, "%s Merge request #%s\n", utils.RedCheck(), i2)
				fmt.Fprintln(out, mrutils.DisplayMR(mr))
			}

			return nil
		},
	}

	return mrCloseCmd
}
