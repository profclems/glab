package delete

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutils.Factory) *cobra.Command {
	var mrDeleteCmd = &cobra.Command{
		Use:     "delete [<id> | <branch>]",
		Short:   `Delete merge requests`,
		Long:    ``,
		Aliases: []string{"del"},
		Args:    cobra.MaximumNArgs(1),
		Example: "$ glab delete 123",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.StdOut, "- Deleting Merge Request !%d\n", mr.IID)

			if err = api.DeleteMR(apiClient, repo.FullName(), mr.IID); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.StdOut, "%s Merge request !%d deleted\n", utils.RedCheck(), mr.IID)

			return nil
		},
	}

	return mrDeleteCmd
}
