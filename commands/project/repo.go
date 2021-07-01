package project

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	repoCmdArchive "github.com/profclems/glab/commands/project/archive"
	repoCmdClone "github.com/profclems/glab/commands/project/clone"
	repoCmdContributors "github.com/profclems/glab/commands/project/contributors"
	repoCmdCreate "github.com/profclems/glab/commands/project/create"
	repoCmdDelete "github.com/profclems/glab/commands/project/delete"
	repoCmdFork "github.com/profclems/glab/commands/project/fork"
	repoCmdSearch "github.com/profclems/glab/commands/project/search"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
)

type RepoOpts struct {
	OpenInBrowser bool
}

func NewCmdRepo(f *cmdutils.Factory) *cobra.Command {
	opts := RepoOpts{}
	var repoCmd = &cobra.Command{
		Use:     "repo <command> [flags]",
		Short:   `Work with GitLab repositories and projects`,
		Long:    ``,
		Aliases: []string{"project"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.OpenInBrowser {
				cfg, err := f.Config()
				if err != nil {
					return err
				}
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				browser, err := cfg.Get(baseRepo.RepoHost(), "browser")
				if err != nil {
					return err
				}
				return utils.OpenInBrowser(fmt.Sprintf("https://%s/%s", baseRepo.RepoHost(), baseRepo.FullName()), browser)
			}
			return nil
		},
	}

	repoCmd.AddCommand(repoCmdArchive.NewCmdArchive(f))
	repoCmd.AddCommand(repoCmdClone.NewCmdClone(f, nil))
	repoCmd.AddCommand(repoCmdContributors.NewCmdContributors(f))
	repoCmd.AddCommand(repoCmdCreate.NewCmdCreate(f))
	repoCmd.AddCommand(repoCmdDelete.NewCmdDelete(f))
	repoCmd.AddCommand(repoCmdFork.NewCmdFork(f, nil))
	repoCmd.AddCommand(repoCmdSearch.NewCmdSearch(f))

	repoCmd.Flags().BoolVarP(&opts.OpenInBrowser, "web", "w", false, "Open repo in a browser. Uses default browser or browser specified in BROWSER variable")

	return repoCmd
}
