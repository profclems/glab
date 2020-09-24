package issue

import (
	"github.com/profclems/glab/commands/cmdutils"
	issueCloseCmd "github.com/profclems/glab/commands/issue/close"

	"github.com/spf13/cobra"
)

func NewCmdIssue(f *cmdutils.Factory) *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "issue [command] [flags]",
		Short: `Work with GitLab issues`,
		Long:  ``,
	}
	issueCmd.AddCommand(issueCloseCmd.NewCmdClose(f))
	issueCmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces")
	return issueCmd
}
