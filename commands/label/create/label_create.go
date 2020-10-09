package create

import (
	"fmt"

	"github.com/profclems/glab/pkg/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdList(f *cmdutils.Factory) *cobra.Command {
	var labelCreateCmd = &cobra.Command{
		Use:     "create [flags]",
		Short:   `Create labels for repository/project`,
		Long:    ``,
		Aliases: []string{"new"},
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			l := &gitlab.CreateLabelOptions{}

			if s, _ := cmd.Flags().GetString("name"); s != "" {
				l.Name = gitlab.String(s)
			}

			if s, _ := cmd.Flags().GetString("color"); s != "" {
				l.Color = gitlab.String(s)
			}
			if s, _ := cmd.Flags().GetString("description"); s != "" {
				l.Description = gitlab.String(s)
			}
			label, err := api.CreateLabel(apiClient, repo.FullName(), l)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "Created label: %s\nWith color: %s\n", label.Name, label.Color)

			return nil
		},
	}
	labelCreateCmd.Flags().StringP("name", "n", "", "Name of label")
	_ = labelCreateCmd.MarkFlagRequired("name")
	labelCreateCmd.Flags().StringP("color", "c", "#428BCA", "Color of label in plain or HEX code. (Default: #428BCA)")
	labelCreateCmd.Flags().StringP("description", "d", "", "Label description")

	return labelCreateCmd
}
