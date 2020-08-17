package commands

import (
	"fmt"
	"strings"

	"glab/internal/git"
	"glab/internal/manip"

	"github.com/gookit/color"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var mrApproversCmd = &cobra.Command{
	Use:     "approvers <id> [flags]",
	Short:   `List merge requests eligible approvers`,
	Long:    ``,
	Aliases: []string{},
	Args:    cobra.ExactArgs(1),
	RunE:    listMergeRequestEligibleApprovers,
}

func listMergeRequestEligibleApprovers(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		mergeID := strings.Trim(args[0], " ")

		fmt.Printf("\nListing Merge Request #%v eligible approvers\n", mergeID)
		gitlabClient, repo := git.InitGitlabClient()
		if r, _ := cmd.Flags().GetString("repo"); r != "" {
			repo = r
		}
		mrApprovals, _, err := gitlabClient.MergeRequestApprovals.GetApprovalState(repo, manip.StringToInt(mergeID))
		if err != nil {
			return err
		}
		if mrApprovals.ApprovalRulesOverwritten {
			color.Yellow.Println("Approval rules overwritten")
		}
		for _, rule := range mrApprovals.Rules {
			table := uitable.New()
			table.MaxColWidth = 70
			if rule.Approved {
				color.Green.Println(fmt.Sprintf("Rule %q sufficient approvals (%d/%d required):", rule.Name, len(rule.ApprovedBy), rule.ApprovalsRequired))
			} else {
				color.Yellow.Println(fmt.Sprintf("Rule %q insufficient approvals (%d/%d required):", rule.Name, len(rule.ApprovedBy), rule.ApprovalsRequired))
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
			fmt.Println(table)
		}

	} else {
		cmdErr(cmd, args)
	}
	return nil
}

func init() {
	mrCmd.AddCommand(mrApproversCmd)
}
