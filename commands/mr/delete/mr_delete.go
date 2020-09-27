package delete

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutils.Factory) *cobra.Command {
	var mrDeleteCmd = &cobra.Command{
		Use:     "delete <id>",
		Short:   `Delete merge requests`,
		Long:    ``,
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		Example: "$ glab delete 123",
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

			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintln(out, "- Deleting Merge Request !"+i2)
				err := api.DeleteMR(apiClient, repo.FullName(), utils.StringToInt(i2))
				if err != nil {
					return err
				}

				fmt.Fprintf(out, "%s Merge request !%s deleted\n", utils.RedCheck(), i2)
			}
			return nil
		},
	}

	return mrDeleteCmd
}
