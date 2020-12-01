package project

import (
	"github.com/profclems/glab/commands/cmdutils"
	repoCmdArchive "github.com/profclems/glab/commands/project/archive"
	repoCmdClone "github.com/profclems/glab/commands/project/clone"
	repoCmdContributors "github.com/profclems/glab/commands/project/contributors"
	repoCmdCreate "github.com/profclems/glab/commands/project/create"
	repoCmdDelete "github.com/profclems/glab/commands/project/delete"
	repoCmdSearch "github.com/profclems/glab/commands/project/search"

	"github.com/spf13/cobra"
)

func NewCmdRepo(f *cmdutils.Factory) *cobra.Command {
	var repoCmd = &cobra.Command{
		Use:     "repo <command> [flags]",
		Short:   `Work with GitLab repositories and projects`,
		Long:    ``,
		Aliases: []string{"project"},
	}

	repoCmd.AddCommand(repoCmdArchive.NewCmdArchive(f))
	repoCmd.AddCommand(repoCmdClone.NewCmdClone(f))
	repoCmd.AddCommand(repoCmdContributors.NewCmdContributors(f))
	repoCmd.AddCommand(repoCmdCreate.NewCmdCreate(f))
	repoCmd.AddCommand(repoCmdDelete.NewCmdDelete(f))
	repoCmd.AddCommand(repoCmdSearch.NewCmdSearch(f))

	return repoCmd
}
