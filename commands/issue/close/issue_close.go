package close

import (
	"fmt"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"strings"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/issue/issueutils"
	"github.com/profclems/glab/internal/manip"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdClose(f *cmdutils.Factory) *cobra.Command {
	var issueCloseCmd = &cobra.Command{
		Use:     "close <id>",
		Short:   `Close an issue`,
		Long:    ``,
		Aliases: []string{"unsub"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			issueID := strings.TrimSpace(args[0])
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

			l := &gitlab.UpdateIssueOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintln(utils.ColorableOut(cmd), "- Closing Issue...")
				issue, err := api.UpdateIssue(apiClient, repo.FullName(), manip.StringToInt(i2), l)
				if err != nil {
					return err
				}
				fmt.Fprintln(utils.ColorableOut(cmd), utils.GreenCheck()+" Issue #"+i2+" closed\n")
				fmt.Fprintln(utils.ColorableOut(cmd), issueutils.DisplayIssue(issue))
			}
			return nil
		},
	}
	return issueCloseCmd
}
