package project

import (
	"github.com/profclems/glab/commands/cmdutils"
	repoCmdArchive "github.com/profclems/glab/commands/project/archive"
	repoCmdClone "github.com/profclems/glab/commands/project/clone"
	repoCmdContributors "github.com/profclems/glab/commands/project/contributors"
	repoCmdCreate "github.com/profclems/glab/commands/project/create"
	repoCmdDelete "github.com/profclems/glab/commands/project/delete"
	repoCmdFork "github.com/profclems/glab/commands/project/fork"
	repoCmdMirror "github.com/profclems/glab/commands/project/mirror"
	repoCmdList "github.com/profclems/glab/commands/project/list"
	repoCmdSearch "github.com/profclems/glab/commands/project/search"
	repoCmdView "github.com/profclems/glab/commands/project/view"

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
	repoCmd.AddCommand(repoCmdClone.NewCmdClone(f, nil))
	repoCmd.AddCommand(repoCmdContributors.NewCmdContributors(f))
	repoCmd.AddCommand(repoCmdList.NewCmdList(f))
	repoCmd.AddCommand(repoCmdCreate.NewCmdCreate(f))
	repoCmd.AddCommand(repoCmdDelete.NewCmdDelete(f))
	repoCmd.AddCommand(repoCmdFork.NewCmdFork(f, nil))
	repoCmd.AddCommand(repoCmdSearch.NewCmdSearch(f))
	repoCmd.AddCommand(repoCmdView.NewCmdView(f))
	repoCmd.AddCommand(repoCmdMirror.NewCmdMirror(f))

	return repoCmd
}
