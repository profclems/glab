package cmdutils

import (
	"os"

	"github.com/spf13/cobra"
)

func EnableRepoOverride(cmd *cobra.Command, f *Factory) {
	cmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` or `GROUP/NAMESPACE/REPO` format or full URL or git URL")

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		repoOverride, err := cmd.Flags().GetString("repo")
		if err != nil {
			return err
		}
		if repoFromEnv := os.Getenv("GITLAB_REPO"); repoOverride == "" && repoFromEnv != "" {
			repoOverride = repoFromEnv
		}
		if repoOverride != "" {
			return f.RepoOverride(repoOverride)
		}
		return nil
	}
}
