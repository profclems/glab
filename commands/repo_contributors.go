package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var repoContributorsCmd = &cobra.Command{
	Use:   "contributors <command> [flags]",
	Short: `Get an archive of the repository.`,
	Example: heredoc.Doc(`
	$ glab repo archive profclems/glab
	$ glab repo archive  // Downloads zip file of current repository
	$ glab repo clone profclems/glab mydirectory  // Clones repo into mydirectory
	$ glab repo clone profclems/glab --format=zip   // Finds repo for current user and download in zip format
	`),
	Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	`),
	Aliases: []string{"users"},
	Run: func (cmd *cobra.Command, args []string) {
		gitlabClient, repo := git.InitGitlabClient()
		l := &gitlab.ListContributorsOptions{}
		users, _, err := gitlabClient.Repositories.Contributors(repo, l)
		if err != nil {
			er(err)
		}
		fmt.Printf("Showing users %d of %d on %s\n\n", len(users), len(users), git.GetRepo())
		for _, user := range users {
			color.Printf("%s <gray><%s></> - %d commits - <red>%d deletions</> - <green>%d additions</>\n", user.Name, user.Email, user.Commits, user.Deletions, user.Additions)
		}
	},
}

func init() {
	repoContributorsCmd.Flags().StringP("order", "f", "zip", "Return contributors ordered by name, email, or commits (orders by commit date) fields. Default is commits")
	repoContributorsCmd.Flags().StringP("sort", "s", "", "Return contributors sorted in asc or desc order. Default is asc")
	repoCmd.AddCommand(repoContributorsCmd)
}