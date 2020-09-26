package delete

import (
	"fmt"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdDelete(f *cmdutils.Factory) *cobra.Command {
	var issueDeleteCmd = &cobra.Command{
		Use:     "delete <id>",
		Short:   `Delete an issue`,
		Long:    ``,
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			issueID := strings.TrimSpace(args[0])
			var err error

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

			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Println("- Deleting Issue #" + i2)
				err := api.DeleteIssue(apiClient, repo.FullName(), utils.StringToInt(i2))
				if err != nil {
					return err
				}
				fmt.Fprintln(utils.ColorableOut(cmd), utils.GreenCheck(), "Issue Deleted")
			}
			return nil
		},
	}
	return issueDeleteCmd
}
