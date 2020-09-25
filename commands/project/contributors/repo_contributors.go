package contributors

import (
	"fmt"
	"github.com/profclems/glab/commands/project"

	"github.com/profclems/glab/internal/git"

	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

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
	Run: func(cmd *cobra.Command, args []string) {
		gitlabClient, repo := git.InitGitlabClient()
		l := &gitlab.ListContributorsOptions{}
		users, _, err := gitlabClient.Repositories.Contributors(repo, l)
		if err != nil {
			er(err)
		}
		fmt.Printf("Showing users %d of %d on %s\n\n", len(users), len(users), repo)
		for _, user := range users {
			color.Printf("%s <gray><%s></> - %d commits - <red>%d deletions</> - <green>%d additions</>\n", user.Name, user.Email, user.Commits, user.Deletions, user.Additions)
		}
	},
}

func init() {
	repoContributorsCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	repoContributorsCmd.Flags().StringP("order", "f", "zip", "Return contributors ordered by name, email, or commits (orders by commit date) fields. Default is commits")
	repoContributorsCmd.Flags().StringP("sort", "s", "", "Return contributors sorted in asc or desc order. Default is asc")
	project.repoCmd.AddCommand(repoContributorsCmd)
}
