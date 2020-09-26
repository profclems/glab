package subscribe

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdSubscribe(f *cmdutils.Factory) *cobra.Command {
	var mrSubscribeCmd = &cobra.Command{
		Use:     "subscribe <id>",
		Short:   `Subscribe to merge requests`,
		Long:    ``,
		Aliases: []string{"sub"},
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

			mergeID := args[0]

			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintln(out, "- Subscribing to merge Request #" + i2)
				mr, err := api.SubscribeToMR(apiClient, repo.FullName(), utils.StringToInt(i2), nil)
				if err != nil {
					return err
				}
				fmt.Fprintln(out, utils.GreenCheck(), "You have successfully subscribed to merge request #" + i2)
				fmt.Fprintln(out, mrutils.DisplayMR(mr))
			}

			return nil
		},
	}

	return mrSubscribeCmd
}
