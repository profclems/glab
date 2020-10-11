package search

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdSearch(f *cmdutils.Factory) *cobra.Command {
	var projectSearchCmd = &cobra.Command{
		Use:     "search [flags]",
		Short:   `Search for GitLab repositories and projects by name`,
		Long:    ``,
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"find", "lookup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			out := utils.ColorableOut(cmd)
			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			search, _ := cmd.Flags().GetString("search")
			page, _ := cmd.Flags().GetInt("page")
			perPage, _ := cmd.Flags().GetInt("per-page")

			projects, _, err := apiClient.Search.Projects(search, &gitlab.SearchOptions{
				Page:    page,
				PerPage: perPage,
			})

			if err != nil {
				return err
			}

			title := fmt.Sprintf("Showing results for \"%s\"", search)
			if len(projects) == 0 {
				title = fmt.Sprintf("No results found for \"%s\"", search)
			}

			table := tableprinter.NewTablePrinter()
			for _, p := range projects {
				table.AddCell(utils.Green(string(rune(p.ID))))

				var description string
				if p.Description != "" {
					description = fmt.Sprintf("\n%s", utils.Cyan(p.Description))
				}

				table.AddCellf("%s%s\n%s", strings.ReplaceAll(p.PathWithNamespace, "/", " / "), description, utils.Gray(p.WebURL))
				table.AddCellf("%d stars %d forks %d issues", p.StarCount, p.ForksCount, p.OpenIssuesCount)
				table.AddCellf("updated %s", utils.TimeToPrettyTimeAgo(*p.LastActivityAt))
				table.EndRow()
			}

			fmt.Fprintf(out, "%s\n%s\n", title, table.Render())
			return nil
		},
	}

	projectSearchCmd.Flags().IntP("page", "p", 1, "Page number")
	projectSearchCmd.Flags().IntP("per-page", "P", 20, "Number of items to list per page")
	projectSearchCmd.Flags().StringP("search", "s", "", "A string contained in the project name")
	_ = projectSearchCmd.MarkFlagRequired("search")

	return projectSearchCmd
}
