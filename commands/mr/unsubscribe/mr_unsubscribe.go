package unsubscribe

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdUnsubscribe(f *cmdutils.Factory) *cobra.Command {
	var mrUnsubscribeCmd = &cobra.Command{
		Use:     "unsubscribe <id>",
		Short:   `Unsubscribe from merge requests`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
				Unsubscribed: true,
			}); err != nil {
				return err
			}

			fmt.Fprintf(out, "- Unsubscribing from Merge Request !%d\n", mr.IID)

			mr, err = api.UnsubscribeFromMR(apiClient, repo.FullName(), mr.IID, nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s You have successfully unsubscribed from merge request !%d\n", utils.GreenCheck(), mr.IID)
			fmt.Fprintln(out, mrutils.DisplayMR(mr))

			return nil
		},
	}

	return mrUnsubscribeCmd
}
