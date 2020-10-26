package mr

import (
	"github.com/profclems/glab/commands/cmdutils"
	mrApproveCmd "github.com/profclems/glab/commands/mr/approve"
	mrApproversCmd "github.com/profclems/glab/commands/mr/approvers"
	mrCheckoutCmd "github.com/profclems/glab/commands/mr/checkout"
	mrCloseCmd "github.com/profclems/glab/commands/mr/close"
	mrCreateCmd "github.com/profclems/glab/commands/mr/create"
	mrDeleteCmd "github.com/profclems/glab/commands/mr/delete"
	mrDiffCmd "github.com/profclems/glab/commands/mr/diff"
	mrForCmd "github.com/profclems/glab/commands/mr/for"
	mrIssuesCmd "github.com/profclems/glab/commands/mr/issues"
	mrListCmd "github.com/profclems/glab/commands/mr/list"
	mrMergeCmd "github.com/profclems/glab/commands/mr/merge"
	mrNoteCmd "github.com/profclems/glab/commands/mr/note"
	mrRebaseCmd "github.com/profclems/glab/commands/mr/rebase"
	mrReopenCmd "github.com/profclems/glab/commands/mr/reopen"
	mrRevokeCmd "github.com/profclems/glab/commands/mr/revoke"
	mrSubscribeCmd "github.com/profclems/glab/commands/mr/subscribe"
	mrTodoCmd "github.com/profclems/glab/commands/mr/todo"
	mrUnsubscribeCmd "github.com/profclems/glab/commands/mr/unsubscribe"
	mrUpdateCmd "github.com/profclems/glab/commands/mr/update"
	mrViewCmd "github.com/profclems/glab/commands/mr/view"

	"github.com/spf13/cobra"
)

func NewCmdMR(f *cmdutils.Factory) *cobra.Command {
	var mrCmd = &cobra.Command{
		Use:   "mr <command> [flags]",
		Short: `Create, view and manage merge requests`,
		Long:  ``,
	}

	cmdutils.EnableRepoOverride(mrCmd, f)

	mrCmd.AddCommand(mrApproveCmd.NewCmdApprove(f))
	mrCmd.AddCommand(mrApproversCmd.NewCmdApprovers(f))
	mrCmd.AddCommand(mrCheckoutCmd.NewCmdCheckout(f))
	mrCmd.AddCommand(mrCloseCmd.NewCmdClose(f))
	mrCmd.AddCommand(mrCreateCmd.NewCmdCreate(f))
	mrCmd.AddCommand(mrDeleteCmd.NewCmdDelete(f))
	mrCmd.AddCommand(mrDiffCmd.NewCmdDiff(f, nil))
	mrCmd.AddCommand(mrForCmd.NewCmdFor(f))
	mrCmd.AddCommand(mrIssuesCmd.NewCmdIssues(f))
	mrCmd.AddCommand(mrListCmd.NewCmdList(f))
	mrCmd.AddCommand(mrMergeCmd.NewCmdMerge(f))
	mrCmd.AddCommand(mrNoteCmd.NewCmdNote(f))
	mrCmd.AddCommand(mrRebaseCmd.NewCmdRebase(f))
	mrCmd.AddCommand(mrReopenCmd.NewCmdReopen(f))
	mrCmd.AddCommand(mrRevokeCmd.NewCmdRevoke(f))
	mrCmd.AddCommand(mrSubscribeCmd.NewCmdSubscribe(f))
	mrCmd.AddCommand(mrUnsubscribeCmd.NewCmdUnsubscribe(f))
	mrCmd.AddCommand(mrTodoCmd.NewCmdTodo(f))
	mrCmd.AddCommand(mrUpdateCmd.NewCmdUpdate(f))
	mrCmd.AddCommand(mrViewCmd.NewCmdView(f))

	return mrCmd
}
