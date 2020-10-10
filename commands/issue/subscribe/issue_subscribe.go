package subscribe

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"

	"github.com/spf13/cobra"
)

func NewCmdSubscribe(f *cmdutils.Factory) *cobra.Command {
	var issueSubscribeCmd = &cobra.Command{
		Use:     "subscribe <id>",
		Short:   `Subscribe to an issue`,
		Long:    ``,
		Aliases: []string{"sub"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := utils.ColorableOut(cmd)
			var err error

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
				fmt.Fprintln(out, "- Subscribing to Issue #"+i2)
				issue, err := api.SubscribeToIssue(apiClient, repo.FullName(), utils.StringToInt(i2), nil)
				if err != nil {
					return err
				}
				fmt.Fprintln(out, utils.GreenCheck(), "Subscribed to issue #"+i2)
				fmt.Fprintln(out, issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}

	return issueSubscribeCmd
}
