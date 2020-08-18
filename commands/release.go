package commands

import (
	"fmt"
	"glab/internal/git"
	"glab/internal/utils"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/gosuri/uitable"
	"github.com/xanzy/go-gitlab"
)

func displayRelease(r *gitlab.Release) {
	duration := utils.TimeToPrettyTimeAgo(*r.CreatedAt)
	color.Printf("%s <green>%s</> %s <magenta>(%s)</>\n", r.Name, r.TagName, duration, r.Description)
}

func displayAllReleases(rs []*gitlab.Release) {
	if len(rs) > 0 {

		table := uitable.New()
		table.MaxColWidth = 70
		fmt.Println()
		fmt.Printf("Showing releases %d of %d on %s\n\n", len(rs), len(rs), git.GetRepo())
		for _, r := range rs {
			duration := utils.TimeToPrettyTimeAgo(*r.CreatedAt)
			author := fmt.Sprintf("%s (%s)", r.Author.Name, r.Author.Username)
			table.AddRow(r.Name, r.TagName, author, r.Description, duration)
		}
		fmt.Println(table)
	} else {
		fmt.Println("No releases available on " + git.GetRepo())
	}
}

// releaseCmd is release command
var releaseCmd = &cobra.Command{
	Use:   "release <command> [flags]",
	Short: `Create, view and manage releases`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
