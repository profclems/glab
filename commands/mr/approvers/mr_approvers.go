package approvers

import (
	"fmt"
	"strings"

	"github.com/profclems/glab/pkg/api"

	"github.com/profclems/glab/commands/cmdutils"
	"github.com/profclems/glab/internal/utils"

	"github.com/gosuri/uitable"
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
			var err error
			out := utils.ColorableOut(cmd)
			if r, _ := cmd.Flags().GetString("repo"); r != "" {
				f, err = f.NewClient(r)
				if err != nil {
					return err
				}
			}

			apiClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			repo, err := f.BaseRepo()
			if err != nil {
				return err
			}

			mergeID := strings.Trim(args[0], " ")

			fmt.Fprintf(out, "\nListing Merge Request !%v eligible approvers\n", mergeID)

			mrApprovals, err := api.GetMRApprovalState(apiClient, repo.FullName(), utils.StringToInt(mergeID))
			if err != nil {
				return err
			}
			if mrApprovals.ApprovalRulesOverwritten {
				fmt.Fprintln(out, utils.Yellow("Approval rules overwritten"))
			}
			for _, rule := range mrApprovals.Rules {
				table := uitable.New()
				table.MaxColWidth = 70
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
