package commands

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"glab/internal/config"
	"glab/internal/git"
	"strings"
)

var repoCloneCmd = &cobra.Command{
	Use:   "clone <command> [flags]",
	Short: `Clone or download a repository/project`,
	Example: heredoc.Doc(`
	$ glab repo clone profclems/glab
	$ glab repo clone https://gitlab.com/profclems/glab
	$ glab repo clone profclems/glab mydirectory  // Clones repo into mydirectory
	$ glab repo clone glab --format=zip   // Finds repo for current user and download in zip format 
	`),
	Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	- url/namespace/group/repo
	`),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmdErr(cmd, args)
			return
		}

		repo := args[0]
		fmt.Println(repo)
		if git.IsValidURL(repo) == false {
			repo = config.GetEnv("GITLAB_URI") + "/" + repo
		}
		if !strings.HasSuffix(repo, ".git") {
			repo += ".git"
		}
		git.RunClone(repo, args[1:])
	},
}

func init() {
	repoCmd.AddCommand(repoCloneCmd)
}