package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"strings"
)

var labelListCmd = &cobra.Command{
	Use:     "list [flags]",
	Short:   `List labels in repository`,
	Long:    ``,
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(0),
	Run:     listLabels,
}

func listLabels(cmd *cobra.Command, args []string) {
	l := &gitlab.ListLabelsOptions{}
	if p, _ := cmd.Flags().GetInt("page"); p != 0 {
		l.Page = p
	}
	if p, _ := cmd.Flags().GetInt("per-page"); p != 0 {
		l.PerPage = p
	}

	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo = r
	}
	// List all labels
	labels, _, err := gitlabClient.Labels.ListLabels(repo, l)
	if err != nil {
		er(err)
	}
	fmt.Printf("Showing label %d of %d on %s", len(labels), len(labels), repo)
	fmt.Println()
	for _, label := range labels {
		color.HEX(strings.Trim(label.Color, "#")).Printf("#%d %s\n", label.ID, label.Name)
	}
}

func init() {
	labelListCmd.Flags().IntP("page", "p", 1, "Page number")
	labelListCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	labelCmd.AddCommand(labelListCmd)
}
