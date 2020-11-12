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
		Use:   "close [<id> | <branch>]",
		Short: `Close merge requests`,
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			if err = mrutils.MRCheckErrors(mr, mrutils.MRCheckErrOptions{
				Closed: true,
				Merged: true,
			}); err != nil {
				return err
			}

			mergeID := args[0]

			l := &gitlab.UpdateMergeRequestOptions{}
			l.StateEvent = gitlab.String("close")
			arrIds := strings.Split(strings.Trim(mergeID, "[] "), ",")
			for _, i2 := range arrIds {
				fmt.Fprintf(f.IO.StdOut, "- Closing Merge request...")
				mr, err := api.UpdateMR(apiClient, repo.FullName(), utils.StringToInt(i2), l)
				if err != nil {
					return err
				}
				fmt.Fprintf(f.IO.StdOut, "%s Merge request !%s\n", utils.RedCheck(), i2)
				fmt.Fprintln(f.IO.StdOut, mrutils.DisplayMR(mr))
			}

			return nil
		},
	}

	return mrCloseCmd
}
