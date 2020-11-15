package checkout

import (
	"fmt"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/utils"
	"github.com/profclems/glab/pkg/api"
	"github.com/tcnksm/go-gitconfig"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

type mrCheckoutConfig struct {
	branch string
	track  bool
}

var (
	mrCheckoutCfg mrCheckoutConfig
)

func NewCmdCheckout(f *cmdutils.Factory) *cobra.Command {
	var mrCheckoutCmd = &cobra.Command{
		Use:   "checkout <id>",
		Short: "Checkout to an open merge request",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mrID := utils.StringToInt(args[0])

			mr, err := api.GetMR(apiClient, repo.FullName(), mrID, &gitlab.GetMergeRequestsOptions{})

			if err != nil {
				return err
			}
			if mr == nil {
				return fmt.Errorf("merge Request !%d not found\n", mrID)
			}
			if mrCheckoutCfg.branch == "" {
				mrCheckoutCfg.branch = mr.SourceBranch
			}
			fetchToRef := mrCheckoutCfg.branch

			if mrCheckoutCfg.track {
				if _, err := gitconfig.Local("remote." + mr.Author.Username + ".url"); err != nil {
					mrProject, err := api.GetProject(apiClient, mr.SourceProjectID)
					if err != nil {
						return err
					}
					if _, err := git.AddRemote(mr.Author.Username, mrProject.SSHURLToRepo); err != nil {
						return err
					}
				}
				fetchToRef = fmt.Sprintf("refs/remotes/%s/%s", mr.Author.Username, mr.SourceBranch)
			}
			remotes, err := f.Remotes()
			if err != nil {
				fmt.Println(err)
			}
			repoRemote, err := remotes.FindByRepo(repo.RepoOwner(), repo.RepoName())

			if err != nil {
				fmt.Println(err)
			}

			mrRef := fmt.Sprintf("refs/merge-requests/%d/head", mrID)
			fetchRefSpec := fmt.Sprintf("%s:%s", mrRef, fetchToRef)
			if err := git.RunCmd([]string{"fetch", repoRemote.Name, fetchRefSpec}); err != nil {
				return err
			}

			if mrCheckoutCfg.track {
				if err := git.RunCmd([]string{"branch", "--track", mrCheckoutCfg.branch, fetchToRef}); err != nil {
					return err
				}
			}

			// Check out branch
			if err := git.CheckoutBranch(mrCheckoutCfg.branch); err != nil {
				return err
			}
			return nil
		},
	}
	mrCheckoutCmd.Flags().StringVarP(&mrCheckoutCfg.branch, "branch", "b", "", "checkout merge request with <branch> name")
	mrCheckoutCmd.Flags().BoolVarP(&mrCheckoutCfg.track, "track", "t", false, "set checked out branch to track remote branch, adds remote if needed")

	return mrCheckoutCmd
}
