package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

var repoCloneCmd = &cobra.Command{
	Use:   "clone <command> [flags]",
	Short: `Clone a Gitlab repository/project`,
	Example: heredoc.Doc(`
	$ glab repo clone profclems/glab
	$ glab repo clone https://gitlab.com/profclems/glab
	$ glab repo clone profclems/glab mydirectory  # Clones repo into mydirectory
	$ glab repo clone glab   # clones repo glab for current user 
	$ glab repo clone 4356677   # finds the project by the ID provided and clones it
	`),
	Long: heredoc.Doc(`
	Clone supports these shorthands
	- repo
	- namespace/repo
	- namespace/group/repo
	- project ID
	`),
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmdErr(cmd, args)
			return nil
		}

		var (
			project *gitlab.Project = nil
			err error
		)

		repo := args[0]
		fmt.Println(repo)
		u, _ := currentUser()
		if !git.IsValidURL(repo) {
			// Assuming that repo is a project ID if it is an integer
			if _, err := strconv.ParseInt(repo,10,64); err != nil {
				// Assuming that "/" in the project name means its owned by an organisation
				if !strings.Contains(repo, "/") {
					repo = fmt.Sprintf("%s/%s", u, repo)
				}
			}
			project, err = getProject(repo)
			if err != nil {
				return err
			}
			repo = project.SSHURLToRepo
		} else if !strings.HasSuffix(repo, ".git") {
			repo += ".git"
		}
		_, err = git.RunClone(repo, args[1:])
		if err != nil {
			return err
		}
		// Cloned project was a fork belonging to the user; user is
		// treating fork's ssh url as origin. Add upstream as remote pointing
		// to forked repo's ssh url
		if project != nil {
			if project.ForkedFromProject != nil &&
				strings.Contains(project.PathWithNamespace, u) {
				var dir string
				if len(args) > 1 {
					dir = args[1]
				} else {
					dir = "./"+project.Path
				}
				fProject, err := getProject(project.ForkedFromProject.PathWithNamespace)
				if err != nil {
					return err
				}
				err = git.AddUpstreamRemote(fProject.SSHURLToRepo, dir)
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	repoCmd.AddCommand(repoCloneCmd)
}
