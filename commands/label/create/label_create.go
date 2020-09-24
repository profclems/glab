package create

import (
	"github.com/profclems/glab/internal/git"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var labelCreateCmd = &cobra.Command{
	Use:     "create [flags]",
	Short:   `Create labels for repository/project`,
	Long:    ``,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(0),
	Run:     createLabel,
}

func createLabel(cmd *cobra.Command, args []string) {

	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo, _ = fixRepoNamespace(r)
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
	label, m, err := gitlabClient.Labels.CreateLabel(repo, l)
	if err != nil {
		if m != nil {
			if m.StatusCode == 409 {
				er("Label already exists")
			}
			er(m.Body)
		}
		er(err)
	}
	color.Printf("Created label: %s\nWith color: %s\n", label.Name, label.Color)
}

func init() {
	labelCreateCmd.Flags().StringP("name", "n", "", "Name of label")
	labelCreateCmd.MarkFlagRequired("name")
	labelCreateCmd.Flags().StringP("color", "c", "#428BCA", "Color of label in plain or HEX code. (Default: #428BCA)")
	labelCreateCmd.Flags().StringP("description", "d", "", "Label description")
	labelCmd.AddCommand(labelCreateCmd)
}
