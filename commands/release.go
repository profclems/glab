package commands

import (
	"fmt"
	"github.com/profclems/glab/internal/utils"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/xanzy/go-gitlab"
)

func displayRelease(r *gitlab.Release) {
	duration := utils.TimeToPrettyTimeAgo(*r.CreatedAt)
	color.Printf("%s <green>%s</> %s <magenta>(%s)</>\n", r.Name, r.TagName, duration, r.Description)
}

func displayAllReleases(rs []*gitlab.Release) {
	DisplayList(ListInfo{
		Name:    "releases",
		Columns: []string{"Name", "TagName", "Author", "Description", "CreatedAt"},
		Total:   len(rs),
		GetCellValue: func(ri int, ci int) interface{} {
			row := rs[ri]
			switch ci {
			case 0:
				return row.Name
			case 1:
				return row.TagName
			case 2:
				return fmt.Sprintf("%s (%s)", row.Author.Name, row.Author.Username)
			case 3:
				return row.Description
			case 4:
				return utils.TimeToPrettyTimeAgo(*row.CreatedAt)
			default:
				return ""
			}
		},
	})
}

// releaseCmd is release command
var releaseCmd = &cobra.Command{
	Use:   "release <command> [flags]",
	Short: `Create, view and manage releases`,
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || len(args) > 2 {
			_ = cmd.Help()
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(releaseCmd)
}
