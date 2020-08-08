package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintHelpIssue : display issue command help
func PrintHelpIssue() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("USAGE")
	fmt.Println("  glab issue <subcommand> [flags]")
	fmt.Println()
	fmt.Println("SUBCOMMANDS")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "create", "Create an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "list, ls", "List issues")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "delete", "Delete an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "subscribe", "Subscribe to an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "unsubscribe", "Unsubscribe from an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "open", "Open an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "reopen", "Reopen an issue")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "close", "Close an issue")
}

// PrintHelpIssueCreate : display issue create command help
func PrintHelpIssueCreate() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Create an issue")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue create [flags]")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--title", `Supply a title. Otherwise, you will be prompted for one. (--title="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--description", `Supply a description. Otherwise, you will be prompted for one. (--description="string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--label", `Add label by name. Multiple labels should be comma-separated. Otherwise, you will be prompted for one, though optional (--label="string,string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--assigns", `Assign issue to people by their ID. Multiple values should be comma separated (--assigns=value,value)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--milestone", `Add the issue to a milestone by ID. (--milestone=value)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--confidential", `Set issue as confidential. Optional boolean value (--confidential) or (--confidential=true)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--mr, --resolved-by-merge-request", `Link issue to a merge request by ID. (--mr=id)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--weight", `Set weight of issue`)
}

// PrintHelpIssueList : display issue list command help
func PrintHelpIssueList() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println(`List issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue list [flags]")
	fmt.Println("  glab issue ls [flags]")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--opened", `Get all opened issues (default)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--all", `Show all opened and closed issues`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--closed", `Get the list of closed issues`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--label, --labels", `Search by label name. Multiple labels should be comma-separated. Otherwise, you will be prompted for one, though optional (--label="string,string")`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--milestone", `Search for issues from a milestone by milestone ID. (--milestone=value)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--confidential", `Search for confidential issues. Optional boolean value (--confidential) or (--confidential=true)`)
}

// PrintHelpIssueDelete : display issue delete command help
func PrintHelpIssueDelete() {
	fmt.Println(`Delete issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue delete <id>")
	fmt.Println("  glab issue delete <comma,separated,ids>")
}

// PrintHelpIssueSubscribe : display issue subscribe command help
func PrintHelpIssueSubscribe() {
	fmt.Println(`Subscribe to issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue subscribe <id>")
	fmt.Println("  glab issue subscribe <comma,separated,ids>")
}

// PrintHelpIssueUnsubscribe : display issue unsubscribe command help
func PrintHelpIssueUnsubscribe() {
	fmt.Println(`Unsubscribe from issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue unsubscribe <id>")
	fmt.Println("  glab issue unsubscribe <comma,separated,ids>")
}

// PrintHelpIssueOpen : display issue open command help
func PrintHelpIssueOpen() {
	fmt.Println(`Reopen issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue open <id>")
	fmt.Println("  glab issue open <comma,separated,ids>")
}

// PrintHelpIssueReopen : display issue reopen command help
func PrintHelpIssueReopen() {
	fmt.Println(`Reopen issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue reopen <id>")
	fmt.Println("  glab issue reopen <comma,separated,ids>")
}

// PrintHelpIssueClose : display issue close command help
func PrintHelpIssueClose() {
	fmt.Println(`Close issues`)
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab issue close <id>")
	fmt.Println("  glab issue close <comma,separated,ids>")
}
