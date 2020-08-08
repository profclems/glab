package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintHelpPipeline : display pipeline command help
func PrintHelpPipeline() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("USAGE")
	fmt.Println("  glab pipeline <subcommand> [flags]")
	fmt.Println()
	fmt.Println("SUBCOMMANDS")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "list", "List pipelines")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "delete", "Delete pipelines")
}

// PrintHelpPipelineList : display pipeline list command help
func PrintHelpPipelineList() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("List pipelines")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab pipeline list [flags]")
	fmt.Println()
	fmt.Println("FLAGS")
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--success", `Retrieve successful pipelines (default). Optional boolean value (--success) or (--sucess=true)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--failed", `Retrieve failed pipelines. Optional boolean value (--failed) or (--failed=true)`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--order_by", `Order the output. Valid options: id, status, red, updated_at, user_id`)
	fmt.Fprintf(tabWriter, "  %s\t%s\n", "--sort", `Sort the output. Valid options: asc, desc`)
}

// PrintHelpPipelineDelete : display pipeline delete command help
func PrintHelpPipelineDelete() {
	fmt.Println("Delete pipelines")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab pipeline delete <ID>")
}
