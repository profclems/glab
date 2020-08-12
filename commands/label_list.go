package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
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
	gitlabClient, repo := git.InitGitlabClient()
	if r, _ := cmd.Flags().GetString("repo"); r != "" {
		repo = r
	}
	// List all labels
	labels, _, err := gitlabClient.Labels.ListLabels(repo, nil)
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
	labelCmd.AddCommand(labelListCmd)
}
