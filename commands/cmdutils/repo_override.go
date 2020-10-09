package cmdutils

import (
	"os"

	"github.com/spf13/cobra"
)

func EnableRepoOverride(cmd *cobra.Command, f *Factory) {
	cmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` or `GROUP/NAMESPACE/REPO` format or the project ID or full URL")

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		repoOverride, _ := cmd.Flags().GetString("repo")
		if repoFromEnv := os.Getenv("GITLAB_REPO"); repoOverride == "" && repoFromEnv != "" {
			repoOverride = repoFromEnv
		}
		if repoOverride != "" {
			_ = f.RepoOverride(repoOverride)
		}
	}
}