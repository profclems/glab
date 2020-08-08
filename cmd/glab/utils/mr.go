package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintHelpMr : display merge request help
func PrintHelpMr() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("USAGE")
	fmt.Println("  glab mr <subcommand> [flags]")
	fmt.Println()
	fmt.Println("SUBCOMMANDS")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "create", "Create a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "merge, accept", "Merge a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "list, ls", "List merge requests")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "close", "Close a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "reopen", "Reopen a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "delete", "Delete a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "subscribe", "Subscribe to a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "unsubscribe", "Unsubscribe from a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "issues", "List the issues that will close when the merge request is accepted")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "approve", "Approve a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "revoke", "Unapprove a merge request")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "Update", "Update an existing merge request")
}

// PrintHelpMrCreate : display merge request create help
func PrintHelpMrCreate() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Create a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr create [flags]")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--title", `Supply a title. Otherwise, you will be prompted for one. (--title="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--description", `Supply a description. Otherwise, you will be prompted for one. (--description="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--source", `Supply the source branch. Otherwise, you will be prompted for one. (--source="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--target", `Supply the target branch. Otherwise, you will be prompted for one. (--target="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--label", `Add label by name. Multiple labels should be comma separated. Otherwise, you will be prompted for one, though optional (--label="string,string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--assigns", `Assign merge request to people by their ID. Multiple values should be comma separated (--assigns=value,value)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--milestone", `Add the merge request to a milestone by ID. (--milestone=value)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--allow-collaboration", `Allow commits from members who can merge to the target branch. Optional boolean value (--allow-collaboration) or (--allow-collaboration=true)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--remove-source-branch", `removes the source branch when merged. Optional boolean value (--remove-source-branch) or (--remove-source-branch=true)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--target-project", `The target project ID`)
}

// PrintHelpMrList : display merge request list help
func PrintHelpMrList() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("List merge requests")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr list [flags]")
	fmt.Println("  glab mr ls [flags]")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--opened", `Get all opened merge requests (default)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--all", `Show all opened and closed merge requests`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--closed", `Get the list of closed merge requests`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--label, --labels", `Search for merge requests by label. Multiple labels should be comma separated. Otherwise, you will be prompted for one, though optional (--labels="string,string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--milestone", `Search for merge requests by milestone ID. (--milestone=value)`)
}

// PrintHelpMrDelete : display merge request delete help
func PrintHelpMrDelete() {
	fmt.Println("Delete merge requests")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr delete <ID>")
	fmt.Println("  glab mr delete <comma,separated,IDs>")
}

// PrintHelpMrSubscribe : display merge request subscribe help
func PrintHelpMrSubscribe() {
	fmt.Println("Subscribe to merge requests")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr subscribe <ID>")
	fmt.Println("  glab mr subscribe <comma,separated,IDs>")
}

// PrintHelpMrUnsubscribe : display merge request unsubscribe help
func PrintHelpMrUnsubscribe() {
	fmt.Println("Unsubscribe from merge requests")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr unsubscribe <ID>")
	fmt.Println("  glab mr unsubscribe <comma,separated,IDs>")
}

// PrintHelpMrAccept : display merge request accept help
func PrintHelpMrAccept() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Accept a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr accept <ID>")
	fmt.Println("  glab mr merge <ID>")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--message", `Custom merge commit message`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--squash-message", `Custom squash commit message`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--squash", `Squashes the commits into a single commit on merge`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--remove-source-branch", `Removes the source branch`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--when-pipeline-succeed", `Merges when the pipeline succeeds`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--sha", `If present, then this SHA must match the HEAD of the source branch, otherwise the merge will fail`)
}

// PrintHelpMrClose : display merge request close help
func PrintHelpMrClose() {
	fmt.Println("Close a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr close <ID>")
	fmt.Println("  glab mr close <comma,separated,IDs>")
}

// PrintHelpMrReopen : display merge request reopen help
func PrintHelpMrReopen() {
	fmt.Println("Reopen a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr reopen <ID>")
	fmt.Println("  glab mr reopen <comma,separated,IDs>")
}

// PrintHelpMrIssues : display merge request issues help
func PrintHelpMrIssues() {
	fmt.Println("List the issues that will close when the merge request is accepted")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr issues <ID>")
}

// PrintHelpMrApprove : display merge request approve help
func PrintHelpMrApprove() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Approve a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr approve <ID>")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--sha", `The HEAD of the merge request`)
}

// PrintHelpMrRevoke : display merge request revoke help
func PrintHelpMrRevoke() {
	fmt.Println("Unapprove a merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr revoke <ID>")
}

// PrintHelpMrUpdate : display merge request update help
func PrintHelpMrUpdate() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Update an existing merge request")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab mr update <ID>")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--title", `Update the title of the merge request`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--description", `Update the description of the merge request`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--lock-discussion", `Boolean to set if the discussion should be locked.`)
}
