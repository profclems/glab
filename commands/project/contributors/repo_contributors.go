package contributors

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdContributors(f *cmdutils.Factory) *cobra.Command {
	var repoContributorsCmd = &cobra.Command{
		Use:   "contributors [flags]",
		Short: `Get contributors of the repository.`,
		Example: heredoc.Doc(`
	$ glab repo contributors
	$ glab repo archive  // Downloads zip file of current repository
	`),
		Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	`),
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"users"},
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			out := utils.ColorableOut(cmd)
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

			l := &gitlab.ListContributorsOptions{}
			users, _, err := apiClient.Repositories.Contributors(repo.FullName(), l)
			if err != nil {
				return err
			}
			usersPrintDetails := fmt.Sprintf("Showing users %d of %d on %s\n\n", len(users), len(users), repo.FullName())
			for _, user := range users {
				usersPrintDetails += fmt.Sprintf("%s <%s> - %d commits(-%s +%s\n",
					user.Name, utils.Gray(user.Email), user.Commits, utils.Red(string(rune(user.Deletions))),
					utils.Green(string(rune(user.Additions))))
			}

			fmt.Fprintf(out, usersPrintDetails)
			return err
		},
	}

	repoContributorsCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	repoContributorsCmd.Flags().StringP("order", "f", "zip", "Return contributors ordered by name, email, or commits (orders by commit date) fields. Default is commits")
	repoContributorsCmd.Flags().StringP("sort", "s", "", "Return contributors sorted in asc or desc order. Default is asc")

	return repoContributorsCmd
}
