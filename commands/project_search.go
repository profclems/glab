package commands

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
	"glab/internal/utils"
	"strings"
)


var projectSearchCmd = &cobra.Command{
	Use:     "search [flags]",
	Short:   `Work with GitLab repositories and projects`,
	Long:    ``,
	Aliases: []string{"find","lookup"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			_ = cmd.Help()
			return nil
		}

		gitlabClient, _ := git.InitGitlabClient()

		search, _ := cmd.Flags().GetString("search")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")

		projects, _, err := gitlabClient.Search.Projects(search, &gitlab.SearchOptions{
			Page: page,
			PerPage: perPage,
		})

		if err != nil {
			return err
		}

		DisplayList(ListInfo{
			Name:    "Projects",
			Columns: []string{"", "", "", ""},
			Total:   len(projects),
			Description: fmt.Sprintf("Showing results for \"%s\"", search),
			EmptyMessage: fmt.Sprintf("No results found for \"%s\"", search),
			TableWrap: true,
			GetCellValue: func(ri int, ci int) interface{} {
				p := projects[ri]
				switch ci {
				case 0:
					return color.Green.Sprint(p.ID)
				case 1:
					var description string
					if p.Description != "" {
						description = color.Sprintf("\n<cyan>%s</>", p.Description)
					}
					return fmt.Sprintf("%s%s\n%s",
						strings.ReplaceAll(p.PathWithNamespace, "/", " / "),
						description, color.Gray.Sprint(p.WebURL))
				case 2:
					return fmt.Sprintf("%d stars %d forks %d issues", p.StarCount, p.ForksCount, p.OpenIssuesCount)
				case 3:
					return "updated " +utils.TimeToPrettyTimeAgo(*p.LastActivityAt)
				default:
					return ""
				}
			},
		})
		return nil
	},
}

func init() {
	projectSearchCmd.Flags().IntP("page", "p", 1, "Page number")
	projectSearchCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	projectSearchCmd.Flags().StringP("search", "s", "", "A string contained in the project name")
	projectCmd.AddCommand(projectSearchCmd)
}
