package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintHelpRepo : display repo command help
func PrintHelpRepo() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("USAGE")
	fmt.Println("  glab repo <subcommand> [flags]")
	fmt.Println()
	fmt.Println("SUBCOMMANDS")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "clone", "Clone a repository")
}

// PrintHelpRepoClone : display repo clone command help
func PrintHelpRepoClone() {
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	defer tabWriter.Flush()
	fmt.Println("Clone a repository")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  glab repo clone <owner/repo> [<dir>] [<format>]")
	fmt.Println()
	fmt.Println("OPTIONS")
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "dir", `The directory to save the repository to`)
	fmt.Fprintf(tabWriter, "  %s:\t%s\n", "format", `Clone the repository as an archive. Valid options: "tar.gz", "tar.bz2", "tbz", "tbz2", "tb2", "bz2", "tar", "zip"`)
}
