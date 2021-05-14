package issue

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/commands/cmdutils/action"
	issueBoardCmd "github.com/profclems/glab/commands/issue/board"
	issueCloseCmd "github.com/profclems/glab/commands/issue/close"
	issueCreateCmd "github.com/profclems/glab/commands/issue/create"
	issueDeleteCmd "github.com/profclems/glab/commands/issue/delete"
	issueListCmd "github.com/profclems/glab/commands/issue/list"
	issueNoteCmd "github.com/profclems/glab/commands/issue/note"
	issueReopenCmd "github.com/profclems/glab/commands/issue/reopen"
	issueSubscribeCmd "github.com/profclems/glab/commands/issue/subscribe"
	issueUnsubscribeCmd "github.com/profclems/glab/commands/issue/unsubscribe"
	issueUpdateCmd "github.com/profclems/glab/commands/issue/update"
	issueViewCmd "github.com/profclems/glab/commands/issue/view"

	"github.com/spf13/cobra"
)

func NewCmdIssue(f *cmdutils.Factory) *cobra.Command {
	var issueCmd = &cobra.Command{
		Use:   "issue [command] [flags]",
		Short: `Work with GitLab issues`,
		Long:  ``,
		Example: heredoc.Doc(`
			$ glab issue list
			$ glab issue create --label --confidential
			$ glab issue view --web
			$ glab issue note -m "closing because !123 was merged" <issue number>
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				An issue can be supplied as argument in any of the following formats:
				- by number, e.g. "123"
				- by URL, e.g. "https://gitlab.com/NAMESPACE/REPO/-/issues/123"
			`),
		},
	}

	cmdutils.EnableRepoOverride(issueCmd, f)
	action.EnableRepoOverride(issueCmd, f)

	issueCmd.AddCommand(issueCloseCmd.NewCmdClose(f))
	issueCmd.AddCommand(issueBoardCmd.NewCmdBoard(f))
	issueCmd.AddCommand(issueCreateCmd.NewCmdCreate(f))
	issueCmd.AddCommand(issueDeleteCmd.NewCmdDelete(f))
	issueCmd.AddCommand(issueListCmd.NewCmdList(f, nil))
	issueCmd.AddCommand(issueNoteCmd.NewCmdNote(f))
	issueCmd.AddCommand(issueReopenCmd.NewCmdReopen(f))
	issueCmd.AddCommand(issueViewCmd.NewCmdView(f))
	issueCmd.AddCommand(issueSubscribeCmd.NewCmdSubscribe(f))
	issueCmd.AddCommand(issueUnsubscribeCmd.NewCmdUnsubscribe(f))
	issueCmd.AddCommand(issueUpdateCmd.NewCmdUpdate(f))
	return issueCmd
}
