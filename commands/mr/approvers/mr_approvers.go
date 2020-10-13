package approvers

import (
	"fmt"

	"github.com/profclems/glab/commands/mr/mrutils"
	"github.com/profclems/glab/pkg/api"
	"github.com/profclems/glab/pkg/tableprinter"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func NewCmdApprovers(f *cmdutils.Factory) *cobra.Command {
	var mrApproversCmd = &cobra.Command{
		Use:     "approvers <id> [flags]",
		Short:   `List merge request eligible approvers`,
		Long:    ``,
		Aliases: []string{},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := utils.ColorableOut(cmd)

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			mr, repo, err := mrutils.MRFromArgs(f, args)
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "\nListing Merge Request !%d eligible approvers\n", mr.IID)

			mrApprovals, err := api.GetMRApprovalState(apiClient, repo.FullName(), mr.IID)
			if err != nil {
				return err
			}
			if mrApprovals.ApprovalRulesOverwritten {
				fmt.Fprintln(out, utils.Yellow("Approval rules overwritten"))
			}
			for _, rule := range mrApprovals.Rules {
				table := tableprinter.NewTablePrinter()
				if rule.Approved {
					fmt.Fprintln(out, utils.Green(fmt.Sprintf("Rule %q sufficient approvals (%d/%d required):", rule.Name, len(rule.ApprovedBy), rule.ApprovalsRequired)))
				} else {
					fmt.Fprintln(out, utils.Yellow(fmt.Sprintf("Rule %q insufficient approvals (%d/%d required):", rule.Name, len(rule.ApprovedBy), rule.ApprovalsRequired)))
				}

				eligibleApprovers := rule.EligibleApprovers

				approvedBy := map[string]*gitlab.BasicUser{}
				for _, by := range rule.ApprovedBy {
					approvedBy[by.Username] = by
				}

				for _, eligibleApprover := range eligibleApprovers {
					approved := "-"
					source := ""
					if _, exists := approvedBy[eligibleApprover.Username]; exists {
						approved = "👍"
					}
					if rule.SourceRule != nil {
						source = rule.SourceRule.RuleType
					}
					table.AddRow(eligibleApprover.Name, eligibleApprover.Username, approved, source)
					delete(approvedBy, eligibleApprover.Username)
				}

				for _, approver := range approvedBy {
					approved := "👍"
					table.AddRow(approver.Name, approver.Username, approved, "")
				}
				fmt.Fprintln(out, table)
			}
			return nil
		},
	}

	return mrApproversCmd
}
