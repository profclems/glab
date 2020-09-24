package mr

import (
	"fmt"
	"log"

	"github.com/profclems/glab/internal/config"
	"github.com/profclems/glab/internal/git"
	"github.com/profclems/glab/internal/manip"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
	"github.com/xanzy/go-gitlab"
)

type mrCheckoutConfig struct {
	branch string
	track  bool
}

var (
	mrCheckoutCfg mrCheckoutConfig
)

var mrCheckoutCmd = &cobra.Command{
	Use:   "checkout <mr-id>",
	Short: "Checkout to an open merge request",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mrID := manip.StringToInt(args[0])
		gitlabClient, repo := git.InitGitlabClient()

		mr, _, err := gitlabClient.MergeRequests.GetMergeRequest(repo, mrID, &gitlab.GetMergeRequestsOptions{})

		if err != nil {
			return err
		}
		if mr == nil {
			return fmt.Errorf("merge Request #%d not found\n", mrID)
		}
		if mrCheckoutCfg.branch == "" {
			mrCheckoutCfg.branch = mr.SourceBranch
		}
		fetchToRef := mrCheckoutCfg.branch

		if mrCheckoutCfg.track {
			if _, err := gitconfig.Local("remote." + mr.Author.Username + ".url"); err != nil {
				mrProject, err := getProject(mr.SourceProjectID)
				if err != nil {
					return err
				}
				if _, err := git.AddRemote(mr.Author.Username, mrProject.SSHURLToRepo); err != nil {
					log.Fatal(err)
				}
			}
			fetchToRef = fmt.Sprintf("refs/remotes/%s/%s", mr.Author.Username, mr.SourceBranch)
		}

		mrRef := fmt.Sprintf("refs/merge-requests/%d/head", mrID)
		fetchRefSpec := fmt.Sprintf("%s:%s", mrRef, fetchToRef)
		if err := git.RunCmd([]string{"fetch", config.GetEnv("GIT_REMOTE_URL_VAR"), fetchRefSpec}); err != nil {
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

func init() {
	mrCheckoutCmd.Flags().StringVarP(&mrCheckoutCfg.branch, "branch", "b", "", "checkout merge request with <branch> name")
	mrCheckoutCmd.Flags().BoolVarP(&mrCheckoutCfg.track, "track", "t", false, "set checked out branch to track mr author remote branch, adds remote if needed")
	mrCmd.AddCommand(mrCheckoutCmd)
}
